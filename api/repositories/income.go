package repositories

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/api/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const (
	incomeColl = "income"
	exportColl = "export"
)

type incomeRepository struct {
	session *mongo.Session
}

func NewIncomeReader(session *mongo.Session) usecases.ForGettingIncomeData {
	return &incomeRepository{session}
}

func NewIncomeWriter(session *mongo.Session) usecases.ForControllingIncomeData {
	return &incomeRepository{session}
}

func (r *incomeRepository) GetAllIncomeByRoleStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) ([]*models.Income, error) {
	query := createQueryIncomeByRoleStartDateAndEndDate(role, startDate, endDate)
	return getAllInComeByQuery(r, query)
}

func createQueryIncomeByRoleStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) bson.M {
	query := bson.M{
		"role": role,
		"submitDate": bson.M{
			"$gt": startDate,
			"$lt": endDate,
		},
	}
	return query
}

func getAllInComeByQuery(r *incomeRepository, query bson.M) ([]*models.Income, error) {
	incomes := make([]*models.Income, 0)

	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(query).All(&incomes)
	if err != nil {
		return nil, err
	}
	return incomes, nil
}

func (r *incomeRepository) AddExport(ep *models.Export) error {
	coll := r.session.GetCollection(exportColl)
	ep.ID = bson.NewObjectId()
	return coll.Insert(ep)
}
