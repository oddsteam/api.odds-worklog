package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

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
