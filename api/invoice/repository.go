package invoice

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

type repository struct {
	session *mongo.Session
}

const invoiceColl = "invoice"

func NewRepository(s *mongo.Session) Repository {
	return &repository{s}
}

func (r *repository) Create(i *models.Invoice) (*models.Invoice, error) {
	coll := r.session.GetCollection(invoiceColl)
	i.ID = bson.NewObjectId()
	err := coll.Insert(i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (r *repository) Get() ([]*models.Invoice, error) {
	invoices := make([]*models.Invoice, 0)
	coll := r.session.GetCollection(invoiceColl)
	err := coll.Find(bson.M{}).All(&invoices)
	if err != nil {
		return nil, err
	}
	return invoices, nil
}
