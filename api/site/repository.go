package site

import (
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const siteColl = "site"

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) CreateSiteGroup(site *models.Site) (*models.Site, error) {
	coll := r.session.GetCollection(siteColl)
	site.ID = primitive.NewObjectID()
	ctx := r.session.Ctx()
	_, err := coll.InsertOne(ctx, site)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *repository) UpdateSiteGroup(site *models.Site) (*models.Site, error) {
	coll := r.session.GetCollection(siteColl)
	ctx := r.session.Ctx()
	_, err := coll.UpdateOne(ctx, bson.M{"_id": site.ID}, bson.M{"$set": site})
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *repository) GetSiteGroup() ([]*models.Site, error) {
	sites := make([]*models.Site, 0)
	coll := r.session.GetCollection(siteColl)
	ctx := r.session.Ctx()
	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &sites); err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *repository) GetSiteGroupByID(id string) (*models.Site, error) {
	site := new(models.Site)
	coll := r.session.GetCollection(siteColl)
	ctx := r.session.Ctx()
	err := coll.FindOne(ctx, bson.M{"_id": bsonutil.MustObjectIDFromHex(id)}).Decode(site)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *repository) GetSiteGroupByName(name string) (*models.Site, error) {
	site := new(models.Site)
	coll := r.session.GetCollection(siteColl)
	ctx := r.session.Ctx()
	err := coll.FindOne(ctx, bson.M{"name": name}).Decode(site)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *repository) DeleteSiteGroup(id string) error {
	coll := r.session.GetCollection(siteColl)
	ctx := r.session.Ctx()
	_, err := coll.DeleteOne(ctx, bson.M{"_id": bsonutil.MustObjectIDFromHex(id)})
	return err
}
