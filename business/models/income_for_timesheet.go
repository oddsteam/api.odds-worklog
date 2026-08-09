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
