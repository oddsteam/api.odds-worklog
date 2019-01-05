package income

import (
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
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

func (r *repository) GetIncomeByID(incID, uID string) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(incID), "userId": uID}).One(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (r *repository) UpdateIncome(income *models.Income) error {
	income.LastUpdate = time.Now()
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
