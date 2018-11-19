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

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) CreateUser(u *models.User) (*models.User, error) {
	coll := r.session.GetCollection(userColl)
	u.ID = bson.NewObjectId()
	err := coll.Insert(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) GetUser() ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	err := coll.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetUserByType(corporateFlag string) ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	err := coll.Find(bson.M{"corporateFlag": corporateFlag}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetUserByID(id string) (*models.User, error) {
	user := new(models.User)
	coll := r.session.GetCollection(userColl)
	err := coll.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) GetUserByEmail(email string) (*models.User, error) {
	user := new(models.User)
	coll := r.session.GetCollection(userColl)
	err := coll.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) UpdateUser(user *models.User) (*models.User, error) {
	coll := r.session.GetCollection(userColl)
	err := coll.UpdateId(user.ID, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) DeleteUser(id string) error {
	coll := r.session.GetCollection(userColl)
	return coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
