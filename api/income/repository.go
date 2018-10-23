package income

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

const income = "income"

type repository struct {
	session *mongo.Session
}

func newRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) AddIncome(u *models.Income) (*models.Income, error) {
	// user := new(models.User)
	// colluser := r.session.GetCollection("user")

	// colluser.Find(bson.M{"email": u.CreateBy}).One(&user)
	coll := r.session.GetCollection(income)
	u.ID = bson.NewObjectId()

	err := coll.Insert(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
