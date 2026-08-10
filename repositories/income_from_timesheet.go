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

const incomeFromTimesheetColl = "income_from_timesheet"

type incomeFromTimesheetRepository struct {
	session *mongo.Session
}

func NewIncomeFromTimesheetRepository(session *mongo.Session) usecases.ForGettingIncomeFromTimesheet {
	return &incomeFromTimesheetRepository{session}
}

func (r *incomeFromTimesheetRepository) GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeFromTimesheet, error) {
	fromDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	toDate := fromDate.AddDate(0, 1, 0)
	query := bson.M{
		"userId": userID,
		"submitDate": bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		},
	}

	record := new(models.IncomeFromTimesheet)
	coll := r.session.GetCollection(incomeFromTimesheetColl)
	ctx := r.session.Ctx()
	err := coll.FindOne(ctx, query).Decode(record)
	if errors.Is(err, mongodriver.ErrNoDocuments) {
		return nil, usecases.ErrIncomeFromTimesheetNotFoundForPeriod
	}
	if err != nil {
		return nil, err
	}
	return record, nil
}

func (r *incomeFromTimesheetRepository) Add(income *models.IncomeFromTimesheet) error {
	t := time.Now()
	income.SubmitDate = t
	income.LastUpdate = t
	income.ID = primitive.NewObjectID()
	income.ExportStatus = false
	coll := r.session.GetCollection(incomeFromTimesheetColl)
	ctx := r.session.Ctx()
	_, err := coll.InsertOne(ctx, income)
	return err
}

func (r *incomeFromTimesheetRepository) Update(income *models.IncomeFromTimesheet) error {
	income.LastUpdate = time.Now()
	coll := r.session.GetCollection(incomeFromTimesheetColl)
	ctx := r.session.Ctx()
	_, err := coll.UpdateOne(ctx, bson.M{"_id": income.ID}, bson.M{"$set": income})
	return err
}
