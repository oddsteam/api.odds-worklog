package repositories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func TestIncomeFromTimesheetPeriodQuery(t *testing.T) {
	// A timesheet summary for June is synced in July or later, so submitDate (the write
	// timestamp) is not in June and cannot be used to find the June row. Matching on the period
	// stored on the record is what keeps a re-sync an update instead of a duplicate insert.
	t.Run("matches on the period stored on the record", func(t *testing.T) {
		query := incomeFromTimesheetPeriodQuery("5bbcf2f90fd2df527bc39539", 2026, time.June)

		assert.Equal(t, bson.M{
			"userId": "5bbcf2f90fd2df527bc39539",
			"year":   2026,
			"month":  6,
		}, query)
	})

	t.Run("does not match on the write timestamp", func(t *testing.T) {
		query := incomeFromTimesheetPeriodQuery("5bbcf2f90fd2df527bc39539", 2026, time.June)

		assert.NotContains(t, query, "submitDate")
	})
}

func TestPrepareIncomeFromTimesheetForInsert(t *testing.T) {
	// The sync usecase anchors SubmitDate to the event's period. The repository used to overwrite
	// it with time.Now() on insert, which put every row in the month it was synced.
	t.Run("keeps the period-anchored submitDate set by the caller", func(t *testing.T) {
		period := time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
		record := &models.IncomeFromTimesheet{Year: 2026, Month: 6}
		record.SubmitDate = period

		prepareIncomeFromTimesheetForInsert(record)

		assert.Equal(t, period, record.SubmitDate)
	})

	t.Run("assigns an id, a fresh lastUpdate and an unexported status", func(t *testing.T) {
		record := &models.IncomeFromTimesheet{Year: 2026, Month: 6}
		record.SubmitDate = time.Date(2026, 6, 1, 12, 0, 0, 0, time.UTC)
		record.ExportStatus = true

		prepareIncomeFromTimesheetForInsert(record)

		assert.False(t, record.ID.IsZero())
		assert.False(t, record.ExportStatus)
		assert.WithinDuration(t, time.Now(), record.LastUpdate, time.Minute)
	})
}

func TestIncomeFromTimesheetPeriodIndex(t *testing.T) {
	// One row per user per period is the invariant the sync relies on. The find-then-insert in the
	// sync usecase cannot enforce it on its own under the at-least-once delivery the timesheet
	// service gives us, so the collection carries a unique index as the real guard.
	t.Run("is unique on the user and the period", func(t *testing.T) {
		index := incomeFromTimesheetPeriodIndex()

		assert.Equal(t, bson.D{
			{Key: "userId", Value: 1},
			{Key: "year", Value: 1},
			{Key: "month", Value: 1},
		}, index.Keys)
		assert.NotNil(t, index.Options)
		assert.True(t, *index.Options.Unique)
	})
}
