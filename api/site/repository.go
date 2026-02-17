package site

import (
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"github.com/globalsign/mgo/bson"
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
	site.ID = bson.NewObjectId()
	err := coll.Insert(site)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *repository) UpdateSiteGroup(site *models.Site) (*models.Site, error) {
	coll := r.session.GetCollection(siteColl)
	err := coll.UpdateId(site.ID, &site)
	if err != nil {
		return nil, err
	}
	return site, nil
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
	site := new(models.Site)
	coll := r.session.GetCollection(siteColl)
	err := coll.FindId(bson.ObjectIdHex(id)).One(&site)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *repository) GetSiteGroupByName(name string) (*models.Site, error) {
	site := new(models.Site)
	coll := r.session.GetCollection(siteColl)
	err := coll.Find(bson.M{"name": name}).One(&site)
	if err != nil {
		return nil, err
	}
	return site, nil
}

func (r *repository) DeleteSiteGroup(id string) error {
	coll := r.session.GetCollection(siteColl)
	return coll.Remove(bson.M{"_id": bson.ObjectIdHex(id)})
}
