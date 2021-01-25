package backoffice

import (
	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const userColl = "user"

const backofficeColl = "backoffice"

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) Get() ([]*models.UserIncome, error) {
	users := make([]*models.UserIncome, 0)

	coll := r.session.GetCollection(userColl)

	o1 := bson.M{
		"$addFields": bson.M{"_userId": bson.M{"$toString": "$_id"}},
	}

	o2 := bson.M{
		"$lookup": bson.M{
			"from":         "income",
			"localField":   "_userId",
			"foreignField": "userId",
			"as":           "incomes",
		},
	}
	operations := []bson.M{o1, o2}

	pipe := coll.Pipe(operations)

	// Run the queries and capture the results
	err := pipe.All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetKey() (*models.BackOfficeKey, error) {

	key := new(models.BackOfficeKey)

	coll := r.session.GetCollection(backofficeColl)
	err := coll.Find(bson.M{"key": bson.M{"$ne": nil}}).One(&key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
