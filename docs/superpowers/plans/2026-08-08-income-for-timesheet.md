# Income For Timesheet (Phase 1, Revised) Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking. **Note for this specific plan run:** the human partner is executing this manually, typing every file themselves, with the controller checking each step by reading files back — not dispatching implementer subagents. The task/step structure below still applies as the walkthrough script.

**Goal:** Consume `timesheet.monthly_summary.published` RabbitMQ events and upsert the same payroll-calculated figures Income would have, into a brand-new `income_for_timesheet` collection — with zero changes to the real `Income` collection, `User` model, or existing Add/Update Income usecases.

**Architecture:** `IncomeForTimesheet` embeds `models.Income` (`bson:",inline"`) plus a `Sites` field, so it reuses `models.CreatePayroll`/`UpdatePayroll` unmodified. A new, structurally independent repository (`repositories/income_for_timesheet.go`) owns its own collection — it does not share a struct with the existing income repository, so it is impossible for this code path to touch the real `income` collection. A thin adapter on `api/user`'s existing repository translates "user not found" into a usecase-level sentinel without modifying the shared method other callers depend on. `pkg/timesheetconsumer` is the same design already validated on `feature/timesheet-sync` (message handling + connection lifecycle, including the two fixes that design needed after review), written correctly from the start here.

**Tech Stack:** Go 1.25.1, `go.mongodb.org/mongo-driver`, `github.com/golang/mock` (gomock) + `github.com/stretchr/testify/assert`, `github.com/rabbitmq/amqp091-go`.

## Global Constraints

- Module path: `gitlab.odds.team/worklog/api.odds-worklog` — all internal imports use this prefix.
- TDD red→green→refactor for every behavioral change; run `go test ./...` after each step (per `CLAUDE.md`).
- `business/models` has zero imports from other internal packages (ADR-0001).
- `business/usecases` may only import `business/models` directly; infrastructure is reached through `*_driven_ports.go` interfaces implemented by outer layers.
- **Never touch** `business/models/income.go`, `business/models/user.go`, `business/usecases/add_income.go`, `business/usecases/update_income.go`, or `repositories/income.go` in this plan — every task in this plan only adds new files, or adds a new function to `api/user/repository.go` without changing any existing function there.
- `SpecialIncome` is hardcoded to `"0"` (timesheet sends no OT-rate data) — do not compute or guess a rate.
- No dead-letter queue; no ordering/dedup field — most recently *successfully processed* message wins via plain overwrite.
- RabbitMQ routing key for local dev and this implementation: `timesheet.monthly_summary.published` (confirmed value — do not use `timesheet.monthly_summary` without the suffix).
- Mocks generated with `golang/mock`, never hand-edited; regenerate via the `mockgen` command given in each task.

---

### Task 1: `IncomeForTimesheet` and `SiteWork` models

**Files:**
- Create: `business/models/income_for_timesheet.go`

**Interfaces:**
- Produces: `models.IncomeForTimesheet{Income; Sites []SiteWork}` (embeds `Income`, so every field of `Income` is accessible directly, e.g. `record.NetIncome`), `models.SiteWork{ClientSite, CustomerName string; WorkingDays, OvertimeDays float64}`. Every later task depends on these exact names.

Pure data-shape addition, no behavior to test yet (exercised by Task 3's tests) — just add the file and confirm the module builds.

- [ ] **Step 1: Write the file**

Create `business/models/income_for_timesheet.go`:

```go
package models

// IncomeForTimesheet mirrors Income (embedded, so every Income field is
// available directly) plus the per-site breakdown from the timesheet event.
// It is persisted to its own collection — see repositories/income_for_timesheet.go —
// entirely separate from the real Income collection.
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

- [ ] **Step 2: Confirm it builds**

Run: `go build ./...`
Expected: exits 0, no output.

- [ ] **Step 3: Commit**

```bash
git add business/models/income_for_timesheet.go
git commit -m "Add IncomeForTimesheet and SiteWork models"
```

---

### Task 2: `TimesheetMonthlySummaryEvent` model

**Files:**
- Create: `business/models/timesheet_event.go`
- Create: `business/models/timesheet_event_test.go`

**Interfaces:**
- Produces: `models.TimesheetMonthlySummaryEvent{EventType string; Year, Month int; SummaryAt time.Time; Employee TimesheetEmployee; Sites []TimesheetSiteSummary}`, `models.TimesheetEmployee{Email, EnglishName string}`, `models.TimesheetSiteSummary{ClientSite, CustomerName string; WorkingDays, OvertimeDays float64}`. Task 3 decodes into and reads these exact types/fields.

This is the wire format from timesheet (snake_case JSON keys) — kept separate from `models.SiteWork` (the persisted shape) since they're different concerns that happen to share a shape today.

- [ ] **Step 1: Write the failing test**

Create `business/models/timesheet_event_test.go`:

```go
package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimesheetMonthlySummaryEventUnmarshal(t *testing.T) {
	t.Run("decodes the event catalog's sample payload", func(t *testing.T) {
		raw := []byte(`{
			"event_type": "timesheet.monthly_summary",
			"year": 2026, "month": 6,
			"summary_at": "2026-07-10T15:31:10+07:00",
			"employee": { "email": "employee@odds.team", "english_name": "Jane Doe" },
			"sites": [
				{ "client_site": "SITE-A", "customer_name": "Site A Customer",
				  "working_days": 12.5, "overtime_days": 2.0 }
			]
		}`)

		var evt TimesheetMonthlySummaryEvent
		err := json.Unmarshal(raw, &evt)

		assert.NoError(t, err)
		assert.Equal(t, "timesheet.monthly_summary", evt.EventType)
		assert.Equal(t, 2026, evt.Year)
		assert.Equal(t, 6, evt.Month)
		expectedSummaryAt, _ := time.Parse(time.RFC3339, "2026-07-10T15:31:10+07:00")
		assert.True(t, expectedSummaryAt.Equal(evt.SummaryAt))
		assert.Equal(t, "employee@odds.team", evt.Employee.Email)
		assert.Equal(t, "Jane Doe", evt.Employee.EnglishName)
		assert.Len(t, evt.Sites, 1)
		assert.Equal(t, "SITE-A", evt.Sites[0].ClientSite)
		assert.Equal(t, "Site A Customer", evt.Sites[0].CustomerName)
		assert.Equal(t, 12.5, evt.Sites[0].WorkingDays)
		assert.Equal(t, 2.0, evt.Sites[0].OvertimeDays)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./business/models/... -run TestTimesheetMonthlySummaryEventUnmarshal -v`
Expected: FAIL — `undefined: TimesheetMonthlySummaryEvent`

- [ ] **Step 3: Write the implementation**

Create `business/models/timesheet_event.go`:

```go
package models

import "time"

type TimesheetMonthlySummaryEvent struct {
	EventType string                 `json:"event_type"`
	Year      int                    `json:"year"`
	Month     int                    `json:"month"`
	SummaryAt time.Time              `json:"summary_at"`
	Employee  TimesheetEmployee      `json:"employee"`
	Sites     []TimesheetSiteSummary `json:"sites"`
}

type TimesheetEmployee struct {
	Email       string `json:"email"`
	EnglishName string `json:"english_name"`
}

type TimesheetSiteSummary struct {
	ClientSite   string  `json:"client_site"`
	CustomerName string  `json:"customer_name"`
	WorkingDays  float64 `json:"working_days"`
	OvertimeDays float64 `json:"overtime_days"`
}
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./business/models/... -run TestTimesheetMonthlySummaryEventUnmarshal -v`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add business/models/timesheet_event.go business/models/timesheet_event_test.go
git commit -m "Add TimesheetMonthlySummaryEvent model"
```

---

### Task 3: `sync_income_for_timesheet` usecase

**Files:**
- Create: `business/usecases/sync_income_for_timesheet_driven_ports.go`
- Create: `business/usecases/sync_income_for_timesheet_driving_ports.go`
- Create: `business/usecases/sync_income_for_timesheet.go`
- Create: `business/usecases/sync_income_for_timesheet_test.go`
- Create (generated): `business/usecases/mock/sync_income_for_timesheet_adaptors_mock.go`
- Create (generated): `business/usecases/mock/sync_income_for_timesheet_driving_mock.go`

**Interfaces:**
- Consumes: `models.TimesheetMonthlySummaryEvent`/`TimesheetSiteSummary` (Task 2), `models.SiteWork`, `models.IncomeForTimesheet` (Task 1), `models.CreatePayroll(user models.User, req models.IncomeReq, note string) *models.Income`, `models.UpdatePayroll(user models.User, req models.IncomeReq, note string, record *models.Income) *models.Income`, `models.FloatToString(f float64) string` (all pre-existing, unmodified).
- Produces: `usecases.ForSyncingIncomeForTimesheet` (`SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error`), constructor `usecases.NewSyncIncomeForTimesheetUsecase(incomeRepo ForGettingIncomeForTimesheet, userRepo ForGettingTimesheetUser) ForSyncingIncomeForTimesheet`, sentinels `usecases.ErrTimesheetUserNotFound` and `usecases.ErrIncomeForTimesheetNotFoundForPeriod`, driven ports `ForGettingIncomeForTimesheet{GetByUserYearMonth, Add, Update}` and `ForGettingTimesheetUser{GetByEmail}`. Task 4 (repository wiring) implements both driven ports and must return the two sentinels above for their respective "not found" cases — any other error must propagate unchanged. Task 5 (consumer handler) calls `SyncFromEvent` and branches on `errors.Is(err, ErrTimesheetUserNotFound)`.

- [ ] **Step 1: Write the driven ports**

Create `business/usecases/sync_income_for_timesheet_driven_ports.go`:

```go
package usecases

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// GetByUserYearMonth must return ErrIncomeForTimesheetNotFoundForPeriod (not a raw driver
// error) when no record exists yet for the given user+year+month — any other non-nil error
// is a real failure and must be propagated, not treated as "not found."
type ForGettingIncomeForTimesheet interface {
	GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeForTimesheet, error)
	Add(income *models.IncomeForTimesheet) error
	Update(income *models.IncomeForTimesheet) error
}

// GetByEmail must return ErrTimesheetUserNotFound (not a raw driver error) when no user
// matches the email — any other non-nil error is a real failure and must be propagated.
type ForGettingTimesheetUser interface {
	GetByEmail(email string) (*models.User, error)
}
```

- [ ] **Step 2: Write the driving port**

Create `business/usecases/sync_income_for_timesheet_driving_ports.go`:

```go
package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type ForSyncingIncomeForTimesheet interface {
	SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error
}
```

- [ ] **Step 3: Generate mocks**

Run: `mockgen -source="business/usecases/sync_income_for_timesheet_driven_ports.go" -destination="business/usecases/mock/sync_income_for_timesheet_adaptors_mock.go"`

Run: `mockgen -source="business/usecases/sync_income_for_timesheet_driving_ports.go" -destination="business/usecases/mock/sync_income_for_timesheet_driving_mock.go"`

This produces `mock_usecases.MockForGettingIncomeForTimesheet`, `mock_usecases.MockForGettingTimesheetUser`, and `mock_usecases.MockForSyncingIncomeForTimesheet` (package `mock_usecases`, same as every other file under `business/usecases/mock/`).

- [ ] **Step 4: Write the failing tests**

Create `business/usecases/sync_income_for_timesheet_test.go`:

```go
package usecases

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
)

func timesheetSyncUser() models.User {
	return models.User{
		ID:    bsonutil.MustObjectIDFromHex("5bbcf2f90fd2df527bc39539"),
		Email: "test@abc.com",
		Role:  "individual",
	}
}

func timesheetSyncEvent() models.TimesheetMonthlySummaryEvent {
	return models.TimesheetMonthlySummaryEvent{
		EventType: "timesheet.monthly_summary",
		Year:      2026,
		Month:     6,
		SummaryAt: time.Date(2026, 7, 10, 15, 31, 10, 0, time.UTC),
		Employee:  models.TimesheetEmployee{Email: "test@abc.com", EnglishName: "Tester Super"},
		Sites: []models.TimesheetSiteSummary{
			{ClientSite: "SITE-A", CustomerName: "Site A Customer", WorkingDays: 10, OvertimeDays: 1},
			{ClientSite: "SITE-B", CustomerName: "Site B Customer", WorkingDays: 2.5, OvertimeDays: 1},
		},
	}
}

func TestSyncIncomeForTimesheet(t *testing.T) {
	t.Run("creates a new income_for_timesheet record when none exists, summing days across sites", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		user := timesheetSyncUser()
		evt := timesheetSyncEvent()

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(nil, ErrIncomeForTimesheetNotFoundForPeriod)
		incomeRepo.EXPECT().Add(gomock.Any()).DoAndReturn(func(record *models.IncomeForTimesheet) error {
			assert.Equal(t, "12.50", record.WorkDate)
			assert.Equal(t, "2.00", record.WorkingHours)
			assert.Equal(t, "0", record.SpecialIncome)
			assert.Equal(t, []models.SiteWork{
				{ClientSite: "SITE-A", CustomerName: "Site A Customer", WorkingDays: 10, OvertimeDays: 1},
				{ClientSite: "SITE-B", CustomerName: "Site B Customer", WorkingDays: 2.5, OvertimeDays: 1},
			}, record.Sites)
			return nil
		})

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(evt)

		assert.NoError(t, err)
	})

	t.Run("updates the existing record for the month and preserves its note", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		user := timesheetSyncUser()
		evt := timesheetSyncEvent()
		existing := &models.IncomeForTimesheet{Income: models.MockIncome}
		existing.Note = "existing remark"

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(existing, nil)
		incomeRepo.EXPECT().Update(gomock.Any()).DoAndReturn(func(record *models.IncomeForTimesheet) error {
			assert.Equal(t, "existing remark", record.Note)
			assert.Equal(t, "12.50", record.WorkDate)
			assert.Equal(t, "2.00", record.WorkingHours)
			return nil
		})

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(evt)

		assert.NoError(t, err)
	})

	t.Run("returns ErrTimesheetUserNotFound when GetByEmail returns it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(nil, ErrTimesheetUserNotFound)

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(timesheetSyncEvent())

		assert.ErrorIs(t, err, ErrTimesheetUserNotFound)
	})

	t.Run("propagates a real error from GetByEmail instead of masking it", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(nil, assert.AnError)
		incomeRepo.EXPECT().GetByUserYearMonth(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
		incomeRepo.EXPECT().Add(gomock.Any()).Times(0)
		incomeRepo.EXPECT().Update(gomock.Any()).Times(0)

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(timesheetSyncEvent())

		assert.ErrorIs(t, err, assert.AnError)
	})

	t.Run("propagates a real error from GetByUserYearMonth instead of treating it as not-found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		userRepo := mock_usecases.NewMockForGettingTimesheetUser(ctrl)
		incomeRepo := mock_usecases.NewMockForGettingIncomeForTimesheet(ctrl)

		user := timesheetSyncUser()

		userRepo.EXPECT().GetByEmail("test@abc.com").Return(&user, nil)
		incomeRepo.EXPECT().GetByUserYearMonth(user.ID.Hex(), 2026, time.Month(6)).
			Return(nil, assert.AnError)
		incomeRepo.EXPECT().Add(gomock.Any()).Times(0)
		incomeRepo.EXPECT().Update(gomock.Any()).Times(0)

		uc := NewSyncIncomeForTimesheetUsecase(incomeRepo, userRepo)
		err := uc.SyncFromEvent(timesheetSyncEvent())

		assert.ErrorIs(t, err, assert.AnError)
	})
}
```

- [ ] **Step 5: Run tests to verify they fail**

Run: `go test ./business/usecases/... -run TestSyncIncomeForTimesheet -v`
Expected: FAIL — `undefined: NewSyncIncomeForTimesheetUsecase` (and related undefined symbols)

- [ ] **Step 6: Write the implementation**

Create `business/usecases/sync_income_for_timesheet.go`:

```go
package usecases

import (
	"errors"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

var ErrTimesheetUserNotFound = errors.New("timesheet event: no matching user for employee email")
var ErrIncomeForTimesheetNotFoundForPeriod = errors.New("income_for_timesheet: no record for this user and period")

type syncIncomeForTimesheetUsecase struct {
	incomeRepo ForGettingIncomeForTimesheet
	userRepo   ForGettingTimesheetUser
}

func NewSyncIncomeForTimesheetUsecase(incomeRepo ForGettingIncomeForTimesheet, userRepo ForGettingTimesheetUser) ForSyncingIncomeForTimesheet {
	return &syncIncomeForTimesheetUsecase{incomeRepo, userRepo}
}

func (u *syncIncomeForTimesheetUsecase) SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error {
	user, err := u.userRepo.GetByEmail(evt.Employee.Email)
	if err != nil {
		return err
	}

	var workingDays, overtimeDays float64
	sites := make([]models.SiteWork, 0, len(evt.Sites))
	for _, s := range evt.Sites {
		workingDays += s.WorkingDays
		overtimeDays += s.OvertimeDays
		sites = append(sites, models.SiteWork{
			ClientSite:   s.ClientSite,
			CustomerName: s.CustomerName,
			WorkingDays:  s.WorkingDays,
			OvertimeDays: s.OvertimeDays,
		})
	}

	req := models.IncomeReq{
		WorkDate:      models.FloatToString(workingDays),
		WorkingHours:  models.FloatToString(overtimeDays),
		SpecialIncome: "0",
	}

	existing, err := u.incomeRepo.GetByUserYearMonth(user.ID.Hex(), evt.Year, time.Month(evt.Month))
	switch {
	case errors.Is(err, ErrIncomeForTimesheetNotFoundForPeriod):
		income := models.CreatePayroll(*user, req, "")
		record := &models.IncomeForTimesheet{Income: *income, Sites: sites}
		if err := u.incomeRepo.Add(record); err != nil {
			return err
		}
	case err != nil:
		return err
	default:
		models.UpdatePayroll(*user, req, existing.Note, &existing.Income)
		existing.Sites = sites
		if err := u.incomeRepo.Update(existing); err != nil {
			return err
		}
	}

	return nil
}
```

- [ ] **Step 7: Run tests to verify they pass**

Run: `go test ./business/usecases/... -run TestSyncIncomeForTimesheet -v`
Expected: PASS (all 5 subtests)

- [ ] **Step 8: Run the full test suite**

Run: `go test ./...`
Expected: PASS, no regressions

- [ ] **Step 9: Commit**

```bash
git add business/usecases/sync_income_for_timesheet*.go business/usecases/mock/sync_income_for_timesheet_adaptors_mock.go business/usecases/mock/sync_income_for_timesheet_driving_mock.go
git commit -m "Add sync_income_for_timesheet usecase"
```

---

### Task 4: Repository wiring

**Files:**
- Create: `repositories/income_for_timesheet.go`
- Modify: `api/user/repository.go`

**Interfaces:**
- Consumes: `usecases.ForGettingIncomeForTimesheet`, `usecases.ForGettingTimesheetUser`, `usecases.ErrIncomeForTimesheetNotFoundForPeriod`, `usecases.ErrTimesheetUserNotFound` (Task 3).
- Produces: `repositories.NewIncomeForTimesheetRepository(session *mongo.Session) usecases.ForGettingIncomeForTimesheet`, `user.NewTimesheetUserRepository(session *mongo.Session) usecases.ForGettingTimesheetUser`. Task 7 (`main.go` wiring) calls both.

`repositories/income_for_timesheet.go` is a **brand-new, independent** repository — it does not embed or share a struct with the existing `incomeRepository` in `repositories/income.go`, and that file is not touched at all. `api/user/repository.go` gains one new constructor and one new small adapter type; the existing `repository` struct, `NewRepository`, and `GetByEmail` are unchanged — the adapter only wraps the not-found translation around the existing method.

Both new files need `errors` and `go.mongodb.org/mongo-driver/mongo` imported under the alias `mongodriver` — `pkg/mongo` (this repo's own Mongo session package) is already imported as the bare identifier `mongo` in both files, and the official driver's error-sentinel package (needed for `mongo.ErrNoDocuments`) has the same default identifier, so it must be aliased to avoid a collision.

- [ ] **Step 1: Write the income_for_timesheet repository**

Create `repositories/income_for_timesheet.go`:

```go
package repositories

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const incomeForTimesheetColl = "income_for_timesheet"

type incomeForTimesheetRepository struct {
	session *mongo.Session
}

func NewIncomeForTimesheetRepository(session *mongo.Session) usecases.ForGettingIncomeForTimesheet {
	return &incomeForTimesheetRepository{session}
}

func (r *incomeForTimesheetRepository) GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeForTimesheet, error) {
	fromDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	toDate := fromDate.AddDate(0, 1, 0)
	query := bson.M{
		"userId": userID,
		"submitDate": bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		},
	}

	record := new(models.IncomeForTimesheet)
	coll := r.session.GetCollection(incomeForTimesheetColl)
	ctx := r.session.Ctx()
	err := coll.FindOne(ctx, query).Decode(record)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, usecases.ErrIncomeForTimesheetNotFoundForPeriod
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *incomeForTimesheetRepository) Add(income *models.IncomeForTimesheet) error {
	t := time.Now()
	income.SubmitDate = t
	income.LastUpdate = t
	income.ID = primitive.NewObjectID()
	income.ExportStatus = false
	coll := r.session.GetCollection(incomeForTimesheetColl)
	ctx := r.session.Ctx()
	_, err := coll.InsertOne(ctx, income)
	return err
}

func (r *incomeForTimesheetRepository) Update(income *models.IncomeForTimesheet) error {
	income.LastUpdate = time.Now()
	coll := r.session.GetCollection(incomeForTimesheetColl)
	ctx := r.session.Ctx()
	_, err := coll.UpdateOne(ctx, bson.M{"_id": income.ID}, bson.M{"$set": income})
	return err
}
```

- [ ] **Step 2: Add the user not-found adapter**

In `api/user/repository.go`, add `"errors"` and `mongodriver "go.mongodb.org/mongo-driver/mongo"` to the imports, plus `"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"`. Then add this near `NewRepository`:

```go
func NewTimesheetUserRepository(session *mongo.Session) usecases.ForGettingTimesheetUser {
	return &timesheetUserRepository{&repository{session}}
}

// timesheetUserRepository adapts repository's GetByEmail to usecases.ForGettingTimesheetUser's
// contract: callers expect usecases.ErrTimesheetUserNotFound for "no matching user," not the
// raw mongodriver.ErrNoDocuments the shared method returns.
type timesheetUserRepository struct {
	*repository
}

func (r *timesheetUserRepository) GetByEmail(email string) (*models.User, error) {
	user, err := r.repository.GetByEmail(email)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, usecases.ErrTimesheetUserNotFound
	}
	return user, err
}
```

- [ ] **Step 3: Confirm it builds**

Run: `go build ./...`
Expected: exits 0, no output

- [ ] **Step 4: Run the full test suite**

Run: `go test ./...`
Expected: PASS, no regressions

- [ ] **Step 5: Commit**

```bash
git add repositories/income_for_timesheet.go api/user/repository.go
git commit -m "Wire repositories for sync_income_for_timesheet"
```

---

### Task 5: `pkg/timesheetconsumer` — message decoding and ack/nack decisions

**Files:**
- Create: `pkg/timesheetconsumer/handler.go`
- Create: `pkg/timesheetconsumer/handler_test.go`

**Interfaces:**
- Consumes: `usecases.ForSyncingIncomeForTimesheet.SyncFromEvent(evt models.TimesheetMonthlySummaryEvent) error`, `usecases.ErrTimesheetUserNotFound` (Task 3).
- Produces: `timesheetconsumer.HandleDelivery(d Acker, body []byte, uc usecases.ForSyncingIncomeForTimesheet)`, `timesheetconsumer.Acker{Ack(multiple bool) error; Nack(multiple, requeue bool) error}`. Task 6's connection loop calls `HandleDelivery` for every delivery; `amqp.Delivery` satisfies `Acker` structurally, so Task 6 needs no adapter type.

- [ ] **Step 1: Add the amqp091-go dependency**

Run: `go get github.com/rabbitmq/amqp091-go@v1.10.0`
Run: `go mod tidy`

(This dependency may get pruned back out by `go mod tidy` if nothing in this task's code imports it yet — `Acker` is satisfied structurally, no import needed here. That's expected; Task 6 re-adds it for real when `consumer.go` imports `amqp091-go` directly.)

- [ ] **Step 2: Write the failing tests**

Create `pkg/timesheetconsumer/handler_test.go`:

```go
package timesheetconsumer

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

type fakeAcker struct {
	acked       bool
	ackMultiple bool
	nacked      bool
	nackRequeue bool
}

func (f *fakeAcker) Ack(multiple bool) error {
	f.acked = true
	f.ackMultiple = multiple
	return nil
}

func (f *fakeAcker) Nack(multiple, requeue bool) error {
	f.nacked = true
	f.nackRequeue = requeue
	return nil
}

func validEventBody() []byte {
	return []byte(`{
		"event_type": "timesheet.monthly_summary",
		"year": 2026, "month": 6,
		"summary_at": "2026-07-10T15:31:10+07:00",
		"employee": { "email": "employee@odds.team", "english_name": "Jane Doe" },
		"sites": []
	}`)
}

func TestHandleDelivery(t *testing.T) {
	t.Run("acks on successful sync", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		uc.EXPECT().SyncFromEvent(gomock.Any()).Return(nil)
		acker := &fakeAcker{}

		HandleDelivery(acker, validEventBody(), uc)

		assert.True(t, acker.acked)
		assert.False(t, acker.nacked)
	})

	t.Run("acks and drops on malformed JSON", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		acker := &fakeAcker{}

		HandleDelivery(acker, []byte("not json"), uc)

		assert.True(t, acker.acked)
		assert.False(t, acker.nacked)
	})

	t.Run("acks and drops on unmatched user", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		uc.EXPECT().SyncFromEvent(gomock.Any()).Return(usecases.ErrTimesheetUserNotFound)
		acker := &fakeAcker{}

		HandleDelivery(acker, validEventBody(), uc)

		assert.True(t, acker.acked)
		assert.False(t, acker.nacked)
	})

	t.Run("nacks and requeues on infra error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		uc.EXPECT().SyncFromEvent(gomock.Any()).Return(errors.New("mongo write failed"))
		acker := &fakeAcker{}

		HandleDelivery(acker, validEventBody(), uc)

		assert.False(t, acker.acked)
		assert.True(t, acker.nacked)
		assert.True(t, acker.nackRequeue)
	})

	t.Run("recovers a panic and acks", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		uc := mock_usecases.NewMockForSyncingIncomeForTimesheet(ctrl)
		uc.EXPECT().SyncFromEvent(gomock.Any()).DoAndReturn(func(models.TimesheetMonthlySummaryEvent) error {
			panic("boom")
		})
		acker := &fakeAcker{}

		assert.NotPanics(t, func() {
			HandleDelivery(acker, validEventBody(), uc)
		})
		assert.True(t, acker.acked)
		assert.False(t, acker.nacked)
	})
}
```

- [ ] **Step 3: Run tests to verify they fail**

Run: `go test ./pkg/timesheetconsumer/... -v`
Expected: FAIL — `undefined: HandleDelivery`

- [ ] **Step 4: Write the implementation**

Create `pkg/timesheetconsumer/handler.go`:

```go
package timesheetconsumer

import (
	"encoding/json"
	"errors"
	"log"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
)

// Acker is satisfied by github.com/rabbitmq/amqp091-go's Delivery.
type Acker interface {
	Ack(multiple bool) error
	Nack(multiple, requeue bool) error
}

func HandleDelivery(d Acker, body []byte, uc usecases.ForSyncingIncomeForTimesheet) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("timesheetconsumer: panic recovered, dropping message: %v", r)
			d.Ack(false)
		}
	}()

	var evt models.TimesheetMonthlySummaryEvent
	if err := json.Unmarshal(body, &evt); err != nil {
		log.Printf("timesheetconsumer: malformed event body, dropping: %v", err)
		d.Ack(false)
		return
	}

	err := uc.SyncFromEvent(evt)
	switch {
	case err == nil:
		d.Ack(false)
	case errors.Is(err, usecases.ErrTimesheetUserNotFound):
		log.Printf("timesheetconsumer: %v (email=%s, year=%d, month=%d), dropping", err, evt.Employee.Email, evt.Year, evt.Month)
		d.Ack(false)
	default:
		log.Printf("timesheetconsumer: sync failed (email=%s, year=%d, month=%d), requeueing: %v", evt.Employee.Email, evt.Year, evt.Month, err)
		d.Nack(false, true)
	}
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `go test ./pkg/timesheetconsumer/... -v`
Expected: PASS (all 5 subtests)

- [ ] **Step 6: Commit**

```bash
git add go.mod go.sum pkg/timesheetconsumer/handler.go pkg/timesheetconsumer/handler_test.go
git commit -m "Add timesheetconsumer message handling with ack/nack decisions"
```

---

### Task 6: `pkg/timesheetconsumer` — RabbitMQ connection lifecycle

**Files:**
- Create: `pkg/timesheetconsumer/consumer.go`

**Interfaces:**
- Consumes: `timesheetconsumer.HandleDelivery` (Task 5), `usecases.ForSyncingIncomeForTimesheet` (Task 3).
- Produces: `timesheetconsumer.Config{URL, Exchange, Queue, RoutingKey string}`, `timesheetconsumer.Start(cfg Config, uc usecases.ForSyncingIncomeForTimesheet)`. Task 7's `main.go` calls `Start` in a goroutine.

No unit tests for this file (connection plumbing, no broker available in this environment) — correctness is `go build` here, and Task 8's manual verification against a real local RabbitMQ. This version bakes in, from the start, two things the original `feature/timesheet-sync` design only got right after a review round — read carefully, don't drop either:

1. **Backoff resets only after the full pipeline is up**, not merely after `amqp.Dial` succeeds — otherwise a broker that's reachable but misconfigured (wrong exchange/queue) would retry at ~1s forever instead of backing off toward the 30s cap.
2. **Both the connection's and the channel's `NotifyClose` are watched** — if only the connection is watched, a channel-level fault (failed ack/nack, broker-side channel exception) leaves the delivery loop dead while `Start` blocks forever with no reconnect.

- [ ] **Step 1: Add the amqp091-go dependency for real**

Run: `go get github.com/rabbitmq/amqp091-go@v1.10.0`
Run: `go mod tidy`

This time it should stick — `consumer.go` genuinely imports the package.

- [ ] **Step 2: Write the implementation**

Create `pkg/timesheetconsumer/consumer.go`:

```go
package timesheetconsumer

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
)

type Config struct {
	URL        string
	Exchange   string
	Queue      string
	RoutingKey string
}

const (
	initialBackoff = time.Second
	maxBackoff     = 30 * time.Second
)

// Start connects to RabbitMQ and consumes timesheet.monthly_summary.published events
// until the process exits. It reconnects with exponential backoff on any
// connection or channel failure and never returns.
func Start(cfg Config, uc usecases.ForSyncingIncomeForTimesheet) {
	backoff := initialBackoff
	for {
		conn, err := amqp.Dial(cfg.URL)
		if err != nil {
			log.Printf("timesheetconsumer: connect failed, retrying in %s: %v", backoff, err)
			time.Sleep(backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		closed := make(chan *amqp.Error, 1)
		conn.NotifyClose(closed)

		chClosed, err := consume(conn, cfg, uc)
		if err != nil {
			log.Printf("timesheetconsumer: consume setup failed: %v", err)
			conn.Close()
			time.Sleep(backoff)
			backoff = nextBackoff(backoff)
			continue
		}

		// Reset only after the full pipeline (dial + declare/bind/consume) is
		// up, not merely after a successful Dial — otherwise a healthy
		// connection with a broken config (e.g. wrong exchange/queue name)
		// would retry at ~1s forever instead of backing off toward the cap.
		backoff = initialBackoff

		select {
		case reason := <-closed:
			log.Printf("timesheetconsumer: connection closed, reconnecting: %v", reason)
		case reason := <-chClosed:
			log.Printf("timesheetconsumer: channel closed, reconnecting: %v", reason)
			conn.Close()
		}
	}
}

func consume(conn *amqp.Connection, cfg Config, uc usecases.ForSyncingIncomeForTimesheet) (chan *amqp.Error, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	chClosed := make(chan *amqp.Error, 1)
	ch.NotifyClose(chClosed)

	if err := ch.ExchangeDeclare(cfg.Exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(cfg.Queue, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	if err := ch.QueueBind(q.Name, cfg.RoutingKey, cfg.Exchange, false, nil); err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		for d := range msgs {
			HandleDelivery(d, d.Body, uc)
		}
	}()

	return chClosed, nil
}

func nextBackoff(current time.Duration) time.Duration {
	next := current * 2
	if next > maxBackoff {
		return maxBackoff
	}
	return next
}
```

- [ ] **Step 3: Confirm it builds**

Run: `go build ./...`
Expected: exits 0, no output

- [ ] **Step 4: Run the full test suite**

Run: `go test ./...`
Expected: PASS, no regressions

- [ ] **Step 5: Commit**

```bash
git add go.mod go.sum pkg/timesheetconsumer/consumer.go
git commit -m "Add timesheetconsumer RabbitMQ connection lifecycle"
```

---

### Task 7: Config, `main.go` wiring, and docker-compose

**Files:**
- Modify: `business/models/config.go`
- Modify: `pkg/config/config.go`
- Modify: `main.go`
- Modify: `deployment/local/docker-compose.yaml`

**Interfaces:**
- Consumes: `timesheetconsumer.Config`, `timesheetconsumer.Start` (Task 6), `repositories.NewIncomeForTimesheetRepository`, `user.NewTimesheetUserRepository` (Task 4), `usecases.NewSyncIncomeForTimesheetUsecase` (Task 3).
- Produces: `models.Config.RabbitMQURL/RabbitMQExchange/RabbitMQQueue/RabbitMQRoutingKey` (Task 8's publisher CLI reads the same four fields via `config.Config()`).

`models.Config` is populated via a **positional** struct literal in `pkg/config/config.go` — new fields must be appended in the same order in both the struct definition and the positional literal, or values silently land in the wrong field (this would not be caught by `go build`).

- [ ] **Step 1: Add RabbitMQ fields to `models.Config`**

In `business/models/config.go`:

```go
type Config struct {
	MongoDBHost          string
	MongoDBName          string
	MongoDBConectionPool int
	APIPort              string
	Username             string
	Password             string
	RabbitMQURL          string
	RabbitMQExchange     string
	RabbitMQQueue        string
	RabbitMQRoutingKey   string
}
```

- [ ] **Step 2: Read the new env vars**

In `pkg/config/config.go`:

```go
func Config() *models.Config {
	godotenv.Load()

	cp, _ := strconv.Atoi(os.Getenv("MONGO_DB_CONECTION_POOL"))
	config := models.Config{
		os.Getenv("MONGO_DB_HOST"),
		os.Getenv("MONGO_DB_NAME"),
		cp,
		os.Getenv("API_PORT"),
		os.Getenv("MONGO_DB_USERNAME"),
		os.Getenv("MONGO_DB_PASSWORD"),
		os.Getenv("RABBITMQ_URL"),
		os.Getenv("RABBITMQ_EXCHANGE"),
		os.Getenv("RABBITMQ_QUEUE"),
		os.Getenv("RABBITMQ_ROUTING_KEY"),
	}
	return &config
}
```

- [ ] **Step 3: Start the consumer from `main.go`**

Add these imports to `main.go`:

```go
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/timesheetconsumer"
	"gitlab.odds.team/worklog/api.odds-worklog/repositories"
```

(`api/user` is already imported as `user` — no change needed there.)

After `session := mongo.Setup()` / `defer session.Close()` and before the Echo instance is created, add:

```go
	c := config.Config()
	go timesheetconsumer.Start(timesheetconsumer.Config{
		URL:        c.RabbitMQURL,
		Exchange:   c.RabbitMQExchange,
		Queue:      c.RabbitMQQueue,
		RoutingKey: c.RabbitMQRoutingKey,
	}, usecases.NewSyncIncomeForTimesheetUsecase(
		repositories.NewIncomeForTimesheetRepository(session),
		user.NewTimesheetUserRepository(session),
	))
```

`main.go` already has a later `c := config.Config()` right before `e.Start(c.APIPort)` — delete that later duplicate declaration and let `e.Start(c.APIPort)` reuse this earlier `c` instead. There must be exactly one `config.Config()` call left in the file.

- [ ] **Step 4: Add the RabbitMQ service to local docker-compose**

In `deployment/local/docker-compose.yaml`, add under `services:`:

```yaml
  rabbitmq:
    image: rabbitmq:3-management
    container_name: odds-worklog-rabbitmq
    ports:
      - "5672:5672"
      - "15672:15672"
```

- [ ] **Step 5: Confirm it builds**

Run: `go build ./...`
Expected: exits 0, no output

- [ ] **Step 6: Run the full test suite**

Run: `go test ./...`
Expected: PASS, no regressions

- [ ] **Step 7: Commit**

```bash
git add business/models/config.go pkg/config/config.go main.go deployment/local/docker-compose.yaml
git commit -m "Wire timesheetconsumer into main.go and local docker-compose"
```

---

### Task 8: Local test publisher CLI and manual verification

**Files:**
- Create: `cmd/timesheetpublisher/main.go`
- Create: `cmd/timesheetpublisher/sample_event.json`

**Interfaces:**
- Consumes: `config.Config()` (Task 7) for the same four `RABBITMQ_*` env vars the consumer uses.

- [ ] **Step 1: Write the sample event fixture**

Create `cmd/timesheetpublisher/sample_event.json` (email must match a real user in your local Mongo — adjust `employee.email` to a user you've created locally before running this):

```json
{
  "event_type": "timesheet.monthly_summary",
  "year": 2026,
  "month": 8,
  "summary_at": "2026-08-08T12:00:00+07:00",
  "employee": { "email": "test@abc.com", "english_name": "Tester Super" },
  "sites": [
    { "client_site": "SITE-A", "customer_name": "Site A Customer", "working_days": 12.5, "overtime_days": 2.0 },
    { "client_site": "SITE-B", "customer_name": "Site B Customer", "working_days": 5, "overtime_days": 0 }
  ]
}
```

- [ ] **Step 2: Write the publisher**

Create `cmd/timesheetpublisher/main.go`:

```go
package main

import (
	"context"
	"flag"
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/config"
)

func main() {
	file := flag.String("file", "", "path to a JSON file containing a timesheet.monthly_summary event")
	flag.Parse()
	if *file == "" {
		log.Fatal("missing required -file flag")
	}

	body, err := os.ReadFile(*file)
	if err != nil {
		log.Fatalf("read %s: %v", *file, err)
	}

	c := config.Config()
	conn, err := amqp.Dial(c.RabbitMQURL)
	if err != nil {
		log.Fatalf("dial %s: %v", c.RabbitMQURL, err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("open channel: %v", err)
	}
	defer ch.Close()

	if err := ch.ExchangeDeclare(c.RabbitMQExchange, "topic", true, false, false, false, nil); err != nil {
		log.Fatalf("declare exchange %s: %v", c.RabbitMQExchange, err)
	}

	err = ch.PublishWithContext(context.Background(), c.RabbitMQExchange, c.RabbitMQRoutingKey, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		log.Fatalf("publish: %v", err)
	}

	log.Printf("published %s to exchange=%s routing_key=%s", *file, c.RabbitMQExchange, c.RabbitMQRoutingKey)
}
```

- [ ] **Step 3: Confirm it builds**

Run: `go build ./...`
Expected: exits 0, no output

- [ ] **Step 4: Add local env vars**

Add to your local `.env` (not committed):

```
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_EXCHANGE=worklog.timesheet
RABBITMQ_QUEUE=worklog.income_for_timesheet
RABBITMQ_ROUTING_KEY=timesheet.monthly_summary.published
```

- [ ] **Step 5: Manual end-to-end verification**

Run: `make up` (brings up the new `rabbitmq` service alongside Mongo/Keycloak), then `make run` in a second terminal.

Run: `go run cmd/timesheetpublisher/main.go -file cmd/timesheetpublisher/sample_event.json`

Expected: publisher logs `published ... to exchange=...`; the API process logs show the delivery being consumed and acked.

Check in Mongo (`docker exec -it odds-worklog-mongo mongosh -u admin -p admin`, then `use odds_worklog_db;`):
- `db.income_for_timesheet.findOne({email:"test@abc.com"})` — `sites` has 2 entries, `workDate` is `"17.50"` (12.5 + 5), `workingHours` is `"2.00"` (2.0 + 0).
- `db.income.findOne({email:"test@abc.com"})` — **must be exactly whatever it was before this test, completely untouched.**
- `db.user.findOne({email:"test@abc.com"})` — **must be exactly whatever it was before this test, completely untouched** (no new fields, no flag).

Run the publisher a second time with the same file: expect no error, and `income_for_timesheet` to still show exactly one document for that user+month (idempotent upsert, no duplicate).

Edit `cmd/timesheetpublisher/sample_event.json`'s `employee.email` to an address with no matching user, republish: expect the API log to show the "dropping" log line and no crash.

- [ ] **Step 6: Commit**

```bash
git add cmd/timesheetpublisher/
git commit -m "Add local timesheet event publisher for manual verification"
```

---

## Self-Review Notes

- **Spec coverage:** `IncomeForTimesheet` embedding (Task 1), event model (Task 2), usecase with the not-found/real-error split specified correctly from the start for *both* `GetByEmail` and `GetByUserYearMonth` (Task 3), fully independent repository + user adapter (Task 4), consumer message handling (Task 5), consumer connection lifecycle with both the backoff-reset fix and the dual-NotifyClose fix baked in from the start (Task 6), config/main.go/docker-compose wiring with the correct `.published` routing key (Task 7), local publisher + manual verification that explicitly checks the real `income`/`user` collections are untouched (Task 8) — all covered.
- **Type consistency:** `ForGettingIncomeForTimesheet`/`ForGettingTimesheetUser` (Task 3) match the repository constructors in Task 4 exactly; `HandleDelivery`/`Acker` (Task 5) match Task 6's `consume` usage; `Config` field names match between Task 6's `timesheetconsumer.Config` and Task 7's `main.go` wiring.
- **No placeholders:** every step has literal code. The two bugs the original `feature/timesheet-sync` design needed a review round to catch (not-found/real-error conflation, connection-only NotifyClose) are specified correctly here from Task 3 and Task 6's first draft — this plan does not silently repeat either mistake.
