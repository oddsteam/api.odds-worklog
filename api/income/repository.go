package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

const incomeColl = "income"

type repository struct {
	session *mongo.Session
}

func newRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) AddIncome(income *models.Income) error {
	coll := r.session.GetCollection(incomeColl)
	income.ID = bson.NewObjectId()

	err := coll.Insert(income)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetIncomeUserNow(id, month string) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)

	err := coll.Find(bson.M{"userId": id, "submitDate": bson.RegEx{month, ""}}).One(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}
