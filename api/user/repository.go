package user

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

const userColl = "user"

type repository struct {
	session *mongo.Session
}

func newRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) createUser(u *models.User) (*models.User, error) {
	coll := r.session.GetCollection(userColl)
	u.ID = bson.NewObjectId()
	err := coll.Insert(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) getUser() ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	err := coll.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) getUserByID(id string) (*models.User, error) {
	user := new(models.User)
	coll := r.session.GetCollection(userColl)
	err := coll.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) updateUser(user *models.User) (*models.User, error) {
	coll := r.session.GetCollection(userColl)
	err := coll.UpdateId(user.ID, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) deleteUser(id string) error {
	coll := r.session.GetCollection(userColl)
	return coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
