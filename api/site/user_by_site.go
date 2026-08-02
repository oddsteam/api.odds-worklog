package site

import (
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

const userColl = "user"

type userBySiteReader struct {
	session *mongo.Session
}

func newUserBySiteReader(session *mongo.Session) ForGettingUsersBySiteID {
	return &userBySiteReader{session}
}

func (r *userBySiteReader) GetBySiteID(id string) ([]*models.User, error) {
	users := make([]*models.User, 0)
	coll := r.session.GetCollection(userColl)
	ctx := r.session.Ctx()
	cursor, err := coll.Find(ctx, bson.M{"siteId": id})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
