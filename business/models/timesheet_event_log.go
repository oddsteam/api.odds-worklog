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
