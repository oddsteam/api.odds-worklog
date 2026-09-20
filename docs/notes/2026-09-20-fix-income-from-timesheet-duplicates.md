# Fix: income_from_timesheet inserted a new row per event instead of updating

## What was wrong

`SyncFromEvent` has had an upsert since day one — look up the row for the period, `Add` if absent,
`Update` if present. The `Update` branch practically never ran, so every timesheet event produced
another row for the same user and period.

The period was never stored on the record. `IncomeFromTimesheet` embeds `Income`, which has no
period field, so `evt.Year`/`evt.Month` were only ever used as lookup arguments and then thrown
away. `GetByUserYearMonth` therefore matched a `submitDate` window instead — and `submitDate` was
stamped `time.Now()` on every write (`repositories/income_from_timesheet.go` on insert,
`models.prepareDataForUpdateIncome` on update). That is the moment the sync ran, not the month being
reported.

A monthly summary arrives after its month has closed, so `submitDate` was never inside the window
being searched and the lookup always missed:

1. Event for June arrives 1 July → no match in `(1 Jun, 1 Jul)` → insert, `submitDate = 1 Jul`
2. Same event redelivered → still no match in `(1 Jun, 1 Jul)` → insert again

The spec calls out at-least-once delivery, so redelivery is normal: rows per user per period equalled
deliveries per user per period.

The same wrong timestamp also put every row in the month it was synced rather than the month it
reports on, which threw off `GetAllByRoleStartDateAndEndDate` — the export's period query.

Origin: `79fd390` (9 Aug 2026), which copied `createQueryByIdAndPeriod` out of
`repositories/income.go`. The pattern is sound there, because a person fills the income form during
the month they are reporting on. It does not transfer to an event-driven sync of a closed month.

It stayed hidden because the local publisher's sample event is for the current month
(`month: 8`, `summary_at: 2026-08-08`), where `submitDate = time.Now()` does land inside the window
and the update branch works; the manual check in the spec only confirmed that a row appears, not
that a second event fails to add one; and every unit test mocks the repository, so the repository's
own query was never exercised.

## What changed

- `IncomeFromTimesheet` carries `Year` and `Month`, and `GetByUserYearMonth` matches on them
  instead of on a `submitDate` window
- `SubmitDate` is anchored to the period being reported (noon UTC on the 1st) rather than to the
  moment of the write, so the export and listing period queries over this collection see the right
  month. The insert path no longer overwrites it
- The same treatment for `mirrorIncomeToTimesheet`, the dual-write from the hand-filled income
  form: those rows share this collection, so they need the same key. Without it they would all
  carry `year: 0, month: 0` and collide on the new index
- A unique index on `userId + year + month`, created at startup, as the actual guard against
  duplicates — the find-then-insert cannot enforce it under at-least-once delivery. Startup only
  logs if it cannot be built, so a not-yet-migrated database still serves

## Deploying this needs a migration

The unique index will not build while the duplicates are still there — Mongo reads a missing field
as null, so every legacy row of one user collides with the others as `(userId, null, null)` — and
rows without a stored period are invisible to the new lookup.

`cmd/resyncincomefromtimesheet` repairs the rows in place:

```bash
go run ./cmd/resyncincomefromtimesheet          # reports what it would do, writes nothing
go run ./cmd/resyncincomefromtimesheet -apply
```

It recovers each legacy row's period, stamps it on, collapses the duplicates of a period to the row
written last, and creates the index. Re-running it is safe.

**It recomputes nothing.** The payroll figures on every row are left exactly as they are, so the
manual OT adjustments the previous release asked people to make through the income form survive.

Two kinds of row live in this collection and the period is recovered differently for each:

- **Written by the timesheet sync.** The sync saves the event log and then the row inside one
  message handler, so a row sits milliseconds after the log entry it came from. Matching a row to
  the nearest event for the same employee within a minute recovers the period it was really for.
- **Mirrored from the income form** by `mirrorIncomeToTimesheet`. These have no event behind them,
  so their own `submitDate` month is used — the period the old lookup already treated them as
  being in, which means they do not move. An earlier draft of this migration replayed the event log
  and deleted every row without a period; that would have destroyed these rows outright, along with
  the rows of anyone whose email no longer matches a user.

## Known behaviour worth watching after deploy

Two events for the same user and period handled concurrently both miss the lookup, so the second
`Add` now fails on the unique index instead of inserting a duplicate. The consumer nacks and
requeues it, the retry finds the row and updates it, and the message is acked — it settles by
itself, but each retry writes another `timesheet_event_log` entry, since the usecase logs the event
before doing anything else.
