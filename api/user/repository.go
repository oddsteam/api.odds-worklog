package user

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"github.com/globalsign/mgo/bson"
)

const userColl = "user"

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) Create(u *models.User) (*models.User, error) {
	t := time.Now()
	u.ID = bson.NewObjectId()
	u.Create = t
	u.LastUpdate = t

	coll := r.session.GetCollection(userColl)
	err := coll.Insert(u)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *repository) Get() ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	err := coll.Find(bson.M{}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetByRole(role string) ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	err := coll.Find(bson.M{"role": role}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetByID(id string) (*models.User, error) {
	user := new(models.User)
	coll := r.session.GetCollection(userColl)
	err := coll.FindId(bson.ObjectIdHex(id)).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) GetBySiteID(id string) ([]*models.User, error) {
	users := make([]*models.User, 0)

	coll := r.session.GetCollection(userColl)
	err := coll.Find(bson.M{"siteId": id}).All(&users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *repository) GetByEmail(email string) (*models.User, error) {
	user := new(models.User)
	coll := r.session.GetCollection(userColl)
	err := coll.Find(bson.M{"email": email}).One(&user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) Update(user *models.User) (*models.User, error) {
	user.LastUpdate = time.Now()
	coll := r.session.GetCollection(userColl)
	err := coll.UpdateId(user.ID, &user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repository) Delete(id string) error {
	coll := r.session.GetCollection(userColl)
	return coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
