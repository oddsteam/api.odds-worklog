package site

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

const siteColl = "site"

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) CreateSiteGroup(sites *models.Site) (*models.Site, error) {
	coll := r.session.GetCollection(siteColl)
	sites.ID = bson.NewObjectId()
	err := coll.Insert(sites)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *repository) UpdateSiteGroup(sites *models.Site) (*models.Site, error) {
	coll := r.session.GetCollection(siteColl)
	err := coll.UpdateId(sites.ID, &sites)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *repository) GetSiteGroup() ([]*models.Site, error) {
	sites := make([]*models.Site, 0)

	coll := r.session.GetCollection(siteColl)
	err := coll.Find(bson.M{}).All(&sites)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *repository) GetSiteGroupByID(id string) (*models.Site, error) {
	sites := new(models.Site)
	coll := r.session.GetCollection(siteColl)
	err := coll.FindId(bson.ObjectIdHex(id)).One(&sites)
	if err != nil {
		return nil, err
	}
	return sites, nil
}

func (r *repository) DeleteSiteGroup(id string) error {
	coll := r.session.GetCollection(siteColl)
	return coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
