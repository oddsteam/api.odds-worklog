package main

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func eventLog(email string, year, month int, receivedAt time.Time) *models.TimesheetEventLog {
	return &models.TimesheetEventLog{
		Year:       year,
		Month:      month,
		Employee:   models.TimesheetEmployee{Email: email},
		ReceivedAt: receivedAt,
	}
}

// legacyRow is a row written before the period was stored: submitDate is the moment it was
// written, and insertedAt (carried by the ObjectID) is the same moment.
func legacyRow(email string, insertedAt time.Time) *models.IncomeFromTimesheet {
	row := &models.IncomeFromTimesheet{}
	row.ID = primitive.NewObjectIDFromTimestamp(insertedAt)
	row.Email = email
	row.SubmitDate = insertedAt
	return row
}

func TestResolvePeriod(t *testing.T) {
	jul10 := time.Date(2026, 7, 10, 9, 0, 0, 0, time.UTC)

	// The sync writes the event log and then the row inside one delivery, so the row that a
	// timesheet event produced sits a moment after its log entry. That is what recovers the period
	// the row was really for — its own submitDate only says when it was written.
	t.Run("takes the period from the event the row was written from", func(t *testing.T) {
		ix := newPeriodIndex([]*models.TimesheetEventLog{eventLog("a@abc.com", 2026, 6, jul10)})

		year, month, fromEvent := ix.resolvePeriod(legacyRow("a@abc.com", jul10.Add(200*time.Millisecond)))

		assert.Equal(t, 2026, year)
		assert.Equal(t, time.June, month)
		assert.True(t, fromEvent)
	})

	t.Run("picks the closest event when the employee has several", func(t *testing.T) {
		ix := newPeriodIndex([]*models.TimesheetEventLog{
			eventLog("a@abc.com", 2026, 6, jul10),
			eventLog("a@abc.com", 2026, 7, jul10.Add(3*time.Hour)),
		})

		year, month, fromEvent := ix.resolvePeriod(legacyRow("a@abc.com", jul10.Add(3*time.Hour).Add(time.Second)))

		assert.Equal(t, 2026, year)
		assert.Equal(t, time.July, month)
		assert.True(t, fromEvent)
	})

	// Rows the income form mirrored into this collection have no event behind them. Their
	// submitDate month is the period the old lookup already treated them as being in, so reading
	// it off keeps them exactly where they have always been.
	t.Run("falls back to the row's own month when no event is close enough", func(t *testing.T) {
		ix := newPeriodIndex([]*models.TimesheetEventLog{eventLog("a@abc.com", 2026, 6, jul10)})

		year, month, fromEvent := ix.resolvePeriod(legacyRow("a@abc.com", jul10.Add(48*time.Hour)))

		assert.Equal(t, 2026, year)
		assert.Equal(t, time.July, month)
		assert.False(t, fromEvent)
	})

	t.Run("never matches an event belonging to another employee", func(t *testing.T) {
		ix := newPeriodIndex([]*models.TimesheetEventLog{eventLog("b@abc.com", 2026, 6, jul10)})

		_, month, fromEvent := ix.resolvePeriod(legacyRow("a@abc.com", jul10.Add(200*time.Millisecond)))

		assert.Equal(t, time.July, month)
		assert.False(t, fromEvent)
	})
}

func TestNewestByLastUpdate(t *testing.T) {
	// Duplicates for one period collapse to the row that was written last: it carries whatever the
	// most recent sync or manual edit left behind.
	t.Run("keeps the most recently written row", func(t *testing.T) {
		older := legacyRow("a@abc.com", time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC))
		older.LastUpdate = time.Date(2026, 7, 1, 0, 0, 0, 0, time.UTC)
		newer := legacyRow("a@abc.com", time.Date(2026, 7, 5, 0, 0, 0, 0, time.UTC))
		newer.LastUpdate = time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)

		keep := newestByLastUpdate([]*models.IncomeFromTimesheet{older, newer})

		assert.Equal(t, newer.ID, keep.ID)
	})
}
