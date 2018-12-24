package po

import (
	"time"

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

func (r *repository) Create(po *models.Po) (*models.Po, error) {
	t := time.Now()
	po.ID = bson.NewObjectId()
	po.Create = t
	po.LastUpdate = t

	coll := r.session.GetCollection(PoColl)
	err := coll.Insert(po)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (r *repository) Update(po *models.Po) (*models.Po, error) {
	po.LastUpdate = time.Now()
	coll := r.session.GetCollection(PoColl)
	err := coll.UpdateId(po.ID, &po)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (r *repository) Get() ([]*models.Po, error) {
	po := make([]*models.Po, 0)
	coll := r.session.GetCollection(PoColl)
	err := coll.Find(bson.M{}).All(&po)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (r *repository) GetByID(id string) (*models.Po, error) {
	po := new(models.Po)
	coll := r.session.GetCollection(PoColl)
	err := coll.FindId(bson.ObjectIdHex(id)).One(&po)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (r *repository) GetByCusID(id string) ([]*models.Po, error) {
	po := make([]*models.Po, 0)

	coll := r.session.GetCollection(PoColl)
	err := coll.Find(bson.M{"customerId": id}).All(&po)
	if err != nil {
		return nil, err
	}
	return po, nil
}
