# Timesheet Sync → Separate `income_for_timesheet` Table (Phase 1, Revised) — Design

Basecamp card #10127003710 Phase 1 ("ดึงข้อมูลจาก timesheet มา add / update income"), revised scope.

## Why this replaces the earlier Phase 1 design

An earlier design (`docs/superpowers/specs/2026-08-06-timesheet-consumer-design.md`, implemented on
branch `feature/timesheet-sync`) wrote timesheet-sourced data directly onto the real `Income`
collection and added a `User.TimesheetSynced` flag that blocked/altered the existing Add/Update
Income flows. That implementation is complete and merged, but the user has decided to hold off on
touching the real flow until the new pipeline is proven. This design instead captures the same
data into a **new, separate collection** used only for future export, with **zero changes to the
real `Income` collection, the real `User` model, or the existing Add/Update Income usecases**.

The two designs consume the same RabbitMQ event and share most of their infrastructure — only the
persistence target and the User-flag/Add-Income-blocking behavior differ.

## Scope

**In scope:** consume `timesheet.monthly_summary.published` events, compute the same
VAT/WHT/Net figures the real payroll flow would compute, and upsert one row per employee per month
into a new `income_for_timesheet` collection, alongside the raw per-site breakdown.

**Out of scope (deliberately, for this phase):**
- Any change to the `income` collection, `models.Income`, `models.User`, or the `add_income`/
  `update_income` usecases.
- A `TimesheetSynced` flag or any Add/Update Income UI behavior change.
- An export endpoint (CSV/SAP) reading from `income_for_timesheet` — this phase only captures the
  data; exporting it is a later phase.

## Event contract

Same event as the original design: `timesheet.monthly_summary.published` — this is the confirmed
real RabbitMQ routing key (the event catalog at `event-catalog-dev.odt.co.th` currently documents
it without the `.published` suffix, but the user has directly confirmed `.published` is correct
for the actual queue binding; treat the catalog as unreliable on this specific point). One message
per employee per month, full-month snapshot (not a delta), at-least-once delivery. Payload shape
(`business/models/timesheet_event.go`, already schema-verified against the event catalog and a
real sample from the timesheet team):

```json
{
  "event_type": "timesheet.monthly_summary",
  "year": 2026, "month": 6,
  "summary_at": "2026-07-10T15:31:10+07:00",
  "employee": { "email": "employee@odds.team", "english_name": "Jane Doe" },
  "sites": [
    { "client_site": "SITE-A", "customer_name": "Site A Customer",
      "working_days": 12.5, "overtime_days": 2.0 }
  ]
}
```

## Architecture

**Models (`business/models/`)**
- `income_for_timesheet.go` (new file):
  ```go
  type IncomeForTimesheet struct {
      Income `bson:",inline"`
      Sites  []SiteWork `bson:"sites" json:"sites,omitempty"`
  }

  type SiteWork struct {
      ClientSite   string  `bson:"clientSite" json:"clientSite"`
      CustomerName string  `bson:"customerName" json:"customerName"`
      WorkingDays  float64 `bson:"workingDays" json:"workingDays"`
      OvertimeDays float64 `bson:"overtimeDays" json:"overtimeDays"`
  }
  ```
  `Income` is embedded (not aliased) so `IncomeForTimesheet` gets every field of `Income` without
  requiring any change to the `Income` struct itself; `bson:",inline"` flattens the embedded
  fields at the same document level instead of nesting them under an `income` key. `Sites` is the
  one field `Income` doesn't have.
- `timesheet_event.go` (new file, same as the original design): `TimesheetMonthlySummaryEvent`,
  `TimesheetEmployee`, `TimesheetSiteSummary` — the wire format, kept separate from `SiteWork`
  (the persisted shape) since they're different concerns.

**Usecase (`business/usecases/`)**
- `sync_income_for_timesheet.go` — the usecase. Given a decoded event: looks up the `User` by
  email (read-only — never writes back to `User`), sums `working_days`/`overtime_days` across
  `sites`, builds an `IncomeReq`, and calls `models.CreatePayroll`/`models.UpdatePayroll` (the
  existing, unmodified payroll calculation functions) to compute a `*models.Income` — then wraps
  it into `IncomeForTimesheet{Income: *result, Sites: sites}` and upserts it via a dedicated
  repository. No `SpecialIncome`/OT-rate data exists on the event, so it's hardcoded to `"0"`,
  same as the original design.
- `sync_income_for_timesheet_driven_ports.go`:
  ```go
  type ForGettingIncomeForTimesheet interface {
      GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeForTimesheet, error)
      Add(income *models.IncomeForTimesheet) error
      Update(income *models.IncomeForTimesheet) error
  }

  type ForGettingTimesheetUser interface {
      GetByEmail(email string) (*models.User, error)
  }
  ```
  `GetByUserYearMonth` must return `ErrIncomeForTimesheetNotFoundForPeriod` (a new sentinel) for
  "no record yet," not a raw driver error — this repeats the same not-found/real-error split that
  the original design initially got wrong and had to fix after review, so it's specified correctly
  from the start this time. `ForGettingTimesheetUser` only needs `GetByEmail` — no `Update`, since
  this design never writes to `User`.
- `sync_income_for_timesheet_driving_ports.go`: `ForSyncingIncomeForTimesheet.SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error`.

**Repository (`repositories/`)**
- `income_for_timesheet.go` — a **new, independent** repository type (does not share a struct
  with the existing `incomeRepository` in `repositories/income.go`), bound to its own collection
  constant `income_for_timesheet`. Being a structurally separate type makes it impossible to
  accidentally read/write the real `income` collection through this code path. Implements
  `GetByUserYearMonth`/`Add`/`Update` against that collection, and translates the driver's
  `mongo.ErrNoDocuments` into `ErrIncomeForTimesheetNotFoundForPeriod` inside `GetByUserYearMonth`
  (mirrors the adapter pattern the original design settled on after review). The `User` lookup
  reuses the existing `api/user` package's `Repository.GetByEmail` — read-only, no new user-side
  code needed.

**Consumer infra (`pkg/timesheetconsumer/`) — unchanged from the original (already-reviewed) design**
- `handler.go` — `HandleDelivery(d Acker, body []byte, uc usecases.ForSyncingIncomeForTimesheet)`,
  same five-branch ack/nack decision logic (success, malformed JSON, unmatched user, infra error,
  panic-recovery), same `Acker` interface satisfied structurally by `amqp.Delivery`.
- `consumer.go` — same connection-lifecycle logic as the original design's **final, fixed** state:
  exponential backoff (1s → 30s cap) that resets only after the full dial+declare/bind/consume
  pipeline succeeds (not merely after `Dial`), and watches **both** the connection's and the
  channel's `NotifyClose` to reconnect on either — both lessons the original design only reached
  after a review round; specified correctly here from the start.
- `pkg/config`/`business/models/config.go` — same `RABBITMQ_URL`/`RABBITMQ_EXCHANGE`/
  `RABBITMQ_QUEUE`/`RABBITMQ_ROUTING_KEY` env vars. Local dev default for the routing key is
  `timesheet.monthly_summary.published` (corrected from the original design's example value,
  which omitted the `.published` suffix).
- `main.go` — starts the consumer as a goroutine, same as before, wired to
  `usecases.NewSyncIncomeForTimesheetUsecase(...)` instead of the old usecase.
- `cmd/timesheetpublisher` — same local test-publisher CLI, unchanged.

## Data flow

1. Consumer decodes the message into `models.TimesheetMonthlySummaryEvent`.
2. `usecase.SyncFromEvent(evt)`:
   - `GetByEmail(evt.Employee.Email)` — not found → `ErrTimesheetUserNotFound` (ack); other error →
     propagate (nack+requeue).
   - Sum `working_days`/`overtime_days` across `sites`; map each entry to `models.SiteWork`.
   - `GetByUserYearMonth(user.ID, evt.Year, evt.Month)` — `ErrIncomeForTimesheetNotFoundForPeriod` →
     build via `CreatePayroll`, wrap into `IncomeForTimesheet`, `Add`; any other error → propagate;
     otherwise → build via `UpdatePayroll` (preserving the existing record's `Note`, same as the
     original design), wrap, `Update`.
   - Return. **No `User` write, no `Income` write, no Add/Update Income usecase touched.**
3. Same overwrite semantics as the original design: no dedup/ordering field, most recently
   successfully processed message wins.

## Error handling

Identical table to the original (already-reviewed) design:

| Situation | Handling |
|---|---|
| Malformed JSON | log + ack |
| Unmatched employee email | log + ack |
| Real error from `GetByEmail` or `GetByUserYearMonth`/`Add`/`Update` | nack + requeue |
| Panic during processing | recover, log, ack |
| RabbitMQ connection **or channel** closes | reconnect with backoff |

## Testing

- `sync_income_for_timesheet_test.go`: create-new, update-existing (note preserved), unmatched
  user → sentinel, real error from `GetByEmail` propagates (not masked), real error from
  `GetByUserYearMonth` propagates instead of being treated as not-found. Every subtest asserts
  **no call reaches any `User`-writing mock** (there is none to call, but the mock for
  `ForGettingTimesheetUser` should only ever expect `GetByEmail`, never `Update`, enforcing this
  at the interface level too).
- `pkg/timesheetconsumer` tests: same as the original design's (already-written) test suite,
  reused as-is since the package is unchanged.
- Manual verification: publish a sample event, confirm a row appears in `income_for_timesheet`
  with the right computed fields and `sites`, confirm the real `income` collection and the
  target user's document are both completely untouched.
