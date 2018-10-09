package user

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

type repository struct {
	session *mongo.Session
}

func newRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) createUser(u *models.User) (*models.User, error) {
	coll := r.session.GetCollection("user")
	err := coll.Insert(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}
