package consumer

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

const collection = "consumer"

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) Create(u *models.Consumer) (*models.Consumer, error) {
	u.ID = bson.NewObjectId()

	coll := r.session.GetCollection(collection)
	err := coll.Insert(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) GetByClientID(cid string) (*models.Consumer, error) {
	consumer := new(models.Consumer)

	coll := r.session.GetCollection(collection)
	err := coll.Find(bson.M{"clientId": cid}).One(&consumer)
	if err != nil {
		return nil, utils.ErrInvalidConsumer
	}
	return consumer, nil
}
