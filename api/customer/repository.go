package customer

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	customerCol = "customer"
	exportColl  = "export"
)

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) CreateCustomer(custmoer *models.Customer) (*models.Customer, error) {
	coll := r.session.GetCollection(customerCol)
	custmoer.ID = bson.NewObjectId()

	err := coll.Insert(custmoer)
	if err != nil {
		return nil, err
	}
	return custmoer, nil
}

func (r *repository) GetCustomers() ([]*models.Customer, error) {
	custmoers := make([]*models.Customer, 0)

	coll := r.session.GetCollection(customerCol)
	err := coll.Find(bson.M{}).All(&custmoers)
	if err != nil {
		return nil, err
	}
	return custmoers, nil
}

func (r *repository) GetCustomerByID(id string) (*models.Customer, error) {
	customer := new(models.Customer)
	coll := r.session.GetCollection(customerCol)
	err := coll.FindId(bson.ObjectIdHex(id)).One(&customer)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (r *repository) UpdateCustomer(custmoer *models.Customer) (*models.Customer, error) {
	coll := r.session.GetCollection(customerCol)
	err := coll.UpdateId(custmoer.ID, &custmoer)
	if err != nil {
		return nil, err
	}
	return custmoer, nil
}

func (r *repository) DeleteCustomer(id string) error {
	coll := r.session.GetCollection(customerCol)
	return coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
