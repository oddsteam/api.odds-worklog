package customer

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"github.com/globalsign/mgo/bson"
)

const customerColl = "customer"

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) Create(customer *models.Customer) (*models.Customer, error) {
	t := time.Now()
	customer.ID = bson.NewObjectId()
	customer.Create = t
	customer.LastUpdate = t

	coll := r.session.GetCollection(customerColl)
	err := coll.Insert(customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *repository) Get() ([]*models.Customer, error) {
	custmoers := make([]*models.Customer, 0)
	coll := r.session.GetCollection(customerColl)
	err := coll.Find(bson.M{}).All(&custmoers)
	if err != nil {
		return nil, err
	}
	return custmoers, nil
}

func (r *repository) GetByID(id string) (*models.Customer, error) {
	customer := new(models.Customer)
	coll := r.session.GetCollection(customerColl)
	err := coll.FindId(bson.ObjectIdHex(id)).One(&customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *repository) Update(custmoer *models.Customer) (*models.Customer, error) {
	custmoer.LastUpdate = time.Now()
	coll := r.session.GetCollection(customerColl)
	err := coll.UpdateId(custmoer.ID, &custmoer)
	if err != nil {
		return nil, err
	}
	return custmoer, nil
}

func (r *repository) Delete(id string) error {
	coll := r.session.GetCollection(customerColl)
	return coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
