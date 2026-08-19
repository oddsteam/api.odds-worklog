package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// NewTimesheetEventLogLister returns the same repository type for read-only listing.
func NewTimesheetEventLogLister(session *mongo.Session) usecases.ForListingTimesheetEventLogs {
	return &timesheetEventLogRepository{session}
}

func (r *timesheetEventLogRepository) List(limit int) ([]*models.TimesheetEventLog, error) {
	coll := r.session.GetCollection(timesheetEventLogColl)
	ctx := r.session.Ctx()
	opts := options.Find().
		SetSort(bson.D{{Key: "receivedAt", Value: -1}}).
		SetLimit(int64(limit))
	cursor, err := coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []*models.TimesheetEventLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
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
