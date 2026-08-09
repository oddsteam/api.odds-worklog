package repositories

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const incomeForTimesheetColl = "income_for_timesheet"

type incomeForTimesheetRepository struct {
	session *mongo.Session
}

func NewIncomeForTimesheetRepository(session *mongo.Session) usecases.ForGettingIncomeForTimesheet {
	return &incomeForTimesheetRepository{session}
}

func (r *incomeForTimesheetRepository) GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeForTimesheet, error) {
	fromDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	toDate := fromDate.AddDate(0, 1, 0)
	query := bson.M{
		"userId": userID,
		"submitDate": bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		},
	}

	record := new(models.IncomeForTimesheet)
	coll := r.session.GetCollection(incomeForTimesheetColl)
	ctx := r.session.Ctx()
	err := coll.FindOne(ctx, query).Decode(record)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, usecases.ErrIncomeForTimesheetNotFoundForPeriod
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *incomeForTimesheetRepository) Add(income *models.IncomeForTimesheet) error {
	t := time.Now()
	income.SubmitDate = t
	income.LastUpdate = t
	income.ID = primitive.NewObjectID()
	income.ExportStatus = false
	coll := r.session.GetCollection(incomeForTimesheetColl)
	ctx := r.session.Ctx()
	_, err := coll.InsertOne(ctx, income)
	return err
}

func (r *incomeForTimesheetRepository) Update(income *models.IncomeForTimesheet) error {
	income.LastUpdate = time.Now()
	coll := r.session.GetCollection(incomeForTimesheetColl)
	ctx := r.session.Ctx()
	_, err := coll.UpdateOne(ctx, bson.M{"_id": income.ID}, bson.M{"$set": income})
	return err
}
