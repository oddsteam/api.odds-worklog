package models

// IncomeFromTimesheet mirrors Income (embedded, so every Income field is
// available directly) plus the per-site breakdown from the timesheet event.
// It is persisted to its own collection — see repositories/income_from_timesheet.go —
// entirely separate from the real Income collection.
//
// Year and Month carry the period the timesheet event was for. The embedded Income has no
// period of its own — in the hand-filled income flow the period is inferred from SubmitDate,
// which works there because people submit during the month they are reporting. A timesheet
// summary arrives after its month has closed, so the period has to be stored explicitly;
// it is the key one row per user per period is looked up and deduplicated by.
type IncomeFromTimesheet struct {
	Income `bson:",inline"`
	Year   int        `bson:"year" json:"year"`
	Month  int        `bson:"month" json:"month"`
	Sites  []SiteWork `bson:"sites" json:"sites,omitempty"`
}

type SiteWork struct {
	ClientSite   string  `bson:"clientSite" json:"clientSite"`
	CustomerName string  `bson:"customerName" json:"customerName"`
	WorkingDays  float64 `bson:"workingDays" json:"workingDays"`
	OvertimeDays float64 `bson:"overtimeDays" json:"overtimeDays"`
}
