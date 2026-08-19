package timesheet_event_log

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

// The nested models.TimesheetEmployee and models.TimesheetSiteSummary are shared with the inbound
// RabbitMQ payload, so their json tags are the publisher's snake_case contract and cannot change.
// The web app reads camelCase, so the HTTP representation is mapped here instead.

type employeeResponse struct {
	Email       string `json:"email"`
	EnglishName string `json:"englishName"`
}

type siteResponse struct {
	ClientSite   string  `json:"clientSite"`
	CustomerName string  `json:"customerName"`
	WorkingDays  float64 `json:"workingDays"`
	OvertimeDays float64 `json:"overtimeDays"`
}

type eventLogResponse struct {
	ID         string           `json:"id"`
	EventType  string           `json:"eventType"`
	Year       int              `json:"year"`
	Month      int              `json:"month"`
	SummaryAt  time.Time        `json:"summaryAt"`
	Employee   employeeResponse `json:"employee"`
	Sites      []siteResponse   `json:"sites"`
	ReceivedAt time.Time        `json:"receivedAt"`
}

func toResponses(logs []*models.TimesheetEventLog) []eventLogResponse {
	responses := make([]eventLogResponse, 0, len(logs))
	for _, log := range logs {
		sites := make([]siteResponse, 0, len(log.Sites))
		for _, site := range log.Sites {
			sites = append(sites, siteResponse{
				ClientSite:   site.ClientSite,
				CustomerName: site.CustomerName,
				WorkingDays:  site.WorkingDays,
				OvertimeDays: site.OvertimeDays,
			})
		}
		responses = append(responses, eventLogResponse{
			ID:        log.ID.Hex(),
			EventType: log.EventType,
			Year:      log.Year,
			Month:     log.Month,
			SummaryAt: log.SummaryAt,
			Employee: employeeResponse{
				Email:       log.Employee.Email,
				EnglishName: log.Employee.EnglishName,
			},
			Sites:      sites,
			ReceivedAt: log.ReceivedAt,
		})
	}
	return responses
}
