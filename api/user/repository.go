package user

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
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

func (r *repository) login(authen *models.Login) (*models.Token, error) {
	username := authen.Username
	password := authen.Password
	if username == "admin" && password == "admin" {
		// Set custom claims
		claims := &models.JwtCustomClaims{
			"Admin!",
			true,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return nil, err
		}
		TK := &models.Token{
			Token: t,
		}
		return TK, nil
	}

	return nil, echo.ErrUnauthorized
}
