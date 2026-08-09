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
