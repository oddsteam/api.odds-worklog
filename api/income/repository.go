package income

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const (
	incomeColl = "income"
	exportColl = "export"
	userColl   = "user"
)

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) AddIncome(income *models.Income) error {
	t := time.Now()
	income.SubmitDate = t
	income.LastUpdate = t
	income.ID = bson.NewObjectId()
	income.ExportStatus = false
	coll := r.session.GetCollection(incomeColl)
	err := coll.Insert(income)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) GetIncomeUserByYearMonth(id string, fromYear int, fromMonth time.Month) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)

	fromDate := time.Date(fromYear, fromMonth, 1, 0, 0, 0, 0, time.UTC)
	toDate := fromDate.AddDate(0, 1, 0)

	err := coll.Find(
		bson.M{
			"userId": id,
			"submitDate": bson.M{
				"$gt": fromDate,
				"$lt": toDate,
			},
		}).One(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}
func (r *repository) GetIncomeByUserIdAllMonth(id string) ([]*models.Income, error) {
	income := make([]*models.Income, 0)

	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(bson.M{"userId": id}).All(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (r *repository) GetIncomeByID(incID, uID string) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(incID), "userId": uID}).One(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (r *repository) GetIncomeByUserID(uID string, fromYear int, fromMonth time.Month) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)

	fromDate := time.Date(fromYear, fromMonth, 1, 0, 0, 0, 0, time.UTC)
	toDate := fromDate.AddDate(0, 1, 0)

	err := coll.Find(bson.M{
		"userId": uID,
		"submitDate": bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		},
		"exportStatus": false,
	}).One(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (r *repository) UpdateIncome(income *models.Income) error {
	income.LastUpdate = time.Now()
	income.ExportStatus = false
	coll := r.session.GetCollection(incomeColl)
	err := coll.UpdateId(income.ID, &income)
	if err != nil {
		return err
	}
	return nil
}

func (r *repository) AddExport(ep *models.Export) error {
	coll := r.session.GetCollection(exportColl)
	ep.ID = bson.NewObjectId()
	return coll.Insert(ep)
}

func (r *repository) DropIncome() error {
	return r.session.GetCollection(incomeColl).DropCollection()
}

func (r *repository) UpdateExportStatus(id string) error {
	coll := r.session.GetCollection(incomeColl)
	err := coll.Update(bson.M{"_id": bson.ObjectIdHex(id)}, bson.M{"$set": bson.M{"exportStatus": true}})
	if err != nil {
		return err
	}
	return nil
}
