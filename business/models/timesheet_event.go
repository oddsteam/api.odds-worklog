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
