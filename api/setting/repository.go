package setting

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

const settingColl = "setting"

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) SaveReminder(reminder *models.Reminder) (*models.Reminder, error) {
	coll := r.session.GetCollection(settingColl)
	reminder.Setting.Time = "23:59"
	selector := bson.M{"name": "reminder"}
	_, err := coll.Upsert(selector, reminder)
	if err != nil {
		return nil, err
	}
	return reminder, nil
}
