# Rename to `income_from_timesheet`, Raw Event Log, and Export Endpoint — Design

Amends `docs/superpowers/specs/2026-08-08-income-for-timesheet-design.md` (branch
`feature/income-for-timesheet`, already implemented and committed). Everything in that spec
still applies except where this document overrides it. Three changes, bundled because the first
two touch the same already-built pipeline and the third builds directly on top of it:

1. Rename `income_for_timesheet` → `income_from_timesheet`, everywhere (Go identifiers, file
   names, Mongo collection name).
2. Add a new `timesheet_event_log` collection that captures the **raw, as-received event
   payload** before any processing — an audit trail independent of whether the event was
   successfully processed.
3. Add a CSV export endpoint reading from `income_from_timesheet`, in its own new module,
   matching the real Income export's format and calculations exactly.

## 1. Rename (mechanical)

Every Go identifier, file name, and the Mongo collection constant changes from
`*IncomeForTimesheet*`/`income_for_timesheet` to `*IncomeFromTimesheet*`/`income_from_timesheet`.
This touches every file created for the original spec:
`business/models/income_for_timesheet.go` → `income_from_timesheet.go`,
`business/usecases/sync_income_for_timesheet*.go` → `sync_income_from_timesheet*.go` (and their
generated mocks), `repositories/income_for_timesheet.go` → `income_from_timesheet.go`, plus every
reference to these names in `main.go` and `pkg/timesheetconsumer` (which only reference the
usecase's driving-port type name, not the collection). No behavior changes in this step — pure
rename, same as an IDE "rename symbol" operation, done file-by-file with `go build`/`go test`
checked after each rename to catch anything missed.

## 2. Raw event log (`timesheet_event_log`)

**Why:** Today, if an event is discarded (unmatched user) or fails at some later step, there is
no record of what was actually received — only log lines, which aren't queryable. Saving the raw
payload first, unconditionally, means every event that ever arrived can be inspected later
regardless of what happened to it downstream.

**Model** (`business/models/timesheet_event_log.go`, new file — does not modify
`business/models/timesheet_event.go` from the original spec):

```go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TimesheetEventLog is the raw, as-received payload from a timesheet.monthly_summary.published
// event, saved before any processing or calculation — an audit trail independent of whether the
// event was ultimately processed successfully.
type TimesheetEventLog struct {
	ID         primitive.ObjectID     `bson:"_id" json:"id"`
	EventType  string                 `bson:"eventType" json:"eventType"`
	Year       int                    `bson:"year" json:"year"`
	Month      int                    `bson:"month" json:"month"`
	SummaryAt  time.Time              `bson:"summaryAt" json:"summaryAt"`
	Employee   TimesheetEmployee      `bson:"employee" json:"employee"`
	Sites      []TimesheetSiteSummary `bson:"sites" json:"sites"`
	ReceivedAt time.Time              `bson:"receivedAt" json:"receivedAt"`
}
```

`TimesheetEmployee`/`TimesheetSiteSummary` are reused from the existing `timesheet_event.go`
(unchanged) — this new struct is intentionally independent of `TimesheetMonthlySummaryEvent`
(not embedded), since the wire-decode struct only has `json` tags today and this is a persisted
shape that needs `bson` tags; keeping them separate avoids retrofitting persistence tags onto a
struct whose only job so far was JSON decoding.

**Repository** (`repositories/timesheet_event_log.go`, new file) — insert-only, no read/update:
```go
const timesheetEventLogColl = "timesheet_event_log"

type timesheetEventLogRepository struct {
	session *mongo.Session
}

func NewTimesheetEventLogRepository(session *mongo.Session) usecases.ForLoggingTimesheetEvent {
	return &timesheetEventLogRepository{session}
}

func (r *timesheetEventLogRepository) Save(evt models.TimesheetMonthlySummaryEvent) error {
	log := models.TimesheetEventLog{
		ID:         primitive.NewObjectID(),
		EventType:  evt.EventType,
		Year:       evt.Year,
		Month:      evt.Month,
		SummaryAt:  evt.SummaryAt,
		Employee:   evt.Employee,
		Sites:      evt.Sites,
		ReceivedAt: time.Now(),
	}
	coll := r.session.GetCollection(timesheetEventLogColl)
	ctx := r.session.Ctx()
	_, err := coll.InsertOne(ctx, log)
	return err
}
```

**Usecase change:** `SyncFromEvent` gains a new dependency `ForLoggingTimesheetEvent{Save(evt
models.TimesheetMonthlySummaryEvent) error}`, and calls `Save` as the very first statement in the
method — before the user lookup, before any calculation. If `Save` fails, `SyncFromEvent` returns
that error immediately (same as any other infra error: `HandleDelivery` nacks and requeues,
per the existing error-handling table — a message we couldn't even log should be retried, not
dropped). The constructor becomes `NewSyncIncomeFromTimesheetUsecase(incomeRepo, userRepo,
eventLogRepo)`.

## 3. Export endpoint

Unchanged from what was already discussed and agreed in this session (approved before this
document was written):

- New, fully independent module `api/income_from_timesheet` (handler + interface), route `GET
  /income-from-timesheet/export/individual/:month`. Zero modifications to `api/income/*`.
- New usecase `business/usecases/export_income_from_timesheet.go` (+ driven ports): fetches
  `[]*models.IncomeFromTimesheet` for the role+month via a new repository method
  `GetAllByRoleStartDateAndEndDate(role string, startDate, endDate time.Time)
  ([]*models.IncomeFromTimesheet, error)` (added to `repositories/income_from_timesheet.go`),
  converts each to `*models.Income` (`&record.Income`), enriches `SiteName` from the user's
  `SiteID` (logic duplicated from `business/usecases/export_income.go`'s private
  `enrichSiteNames` — cannot be called directly since it's a method on that file's unexported
  `usecase` struct; ~15 lines of duplication, isolated in the new file, in exchange for zero
  changes to `export_income.go`), fetches the current student loan list, then reuses
  `models.NewPayrollCycle` and the **existing, unmodified** `pkg/file/csv_writer.go`
  (`WriteFile`/`ToCSV`/`export`) — no new CSV-writing code at all. Output matches the real
  Income export's columns exactly (Site, VAT, WHT, loan deduction, etc.) — only the data source
  collection differs.
- Auth: reuses the existing exported `income.IsIncomeExportAllowed(c)` directly (cross-package
  call, no duplication needed since it's already exported).
- No entry written to the `export` audit-log collection that the real export writes to — this
  keeps the two export histories from mixing in one collection; this endpoint has no equivalent
  audit trail requirement today.
- Error handling matches the real export handler exactly: repo/CSV errors → HTTP 500, missing
  month param → HTTP 400, unauthorized → HTTP 401.

**Design rationale for full isolation (all three changes):** this feature may be temporary — the
`income_from_timesheet`/`timesheet_event_log` pipeline could be removed later once the real
Income-writing design (already implemented on `feature/timesheet-sync`) is adopted instead. Every
file in this spec is new; nothing here modifies a file that existed before this feature. Deleting
the feature later means deleting these files and the two `main.go`/route-registration lines that
reference them — nothing else needs to be touched or untangled.

## Testing

- Rename: existing tests move with their files unchanged in substance — only names change.
  Confirm `go test ./...` is green after the rename before adding new code.
- `timesheet_event_log`: no dedicated unit test for the repository (insert-only, no branching
  logic — consistent with this repo's convention of not unit-testing simple Mongo repositories).
  The usecase test suite gains a mock for `ForLoggingTimesheetEvent` and a new subtest: "returns
  the error when saving the raw event log fails, without looking up the user" (asserts
  `GetByEmail` is never called).
- Export usecase: unit tests mocking the new repo method, asserting the conversion + SiteName
  enrichment produce the right `models.Income` shape passed into `PayrollCycle`/`csvWriter`.
- Manual verification: publish a sample event, confirm a `timesheet_event_log` row appears
  immediately (even before checking `income_from_timesheet`), then hit the new export endpoint
  and confirm the downloaded CSV contains the expected row with the same Site/VAT/WHT figures the
  real export would compute for equivalent inputs.
