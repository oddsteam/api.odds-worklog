package po

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

const PoColl = "po"

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) Create(po *models.Po) (*models.Po, error){
	coll := r.session.GetCollection(PoColl)
	po.ID = bson.NewObjectId()
	err := coll.Insert(po)
	if err != nil {
		return nil, err
	}
	return po, nil
}