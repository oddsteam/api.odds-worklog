package repositories

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func NewIncomeFromTimesheetReader(session *mongo.Session) usecases.ForGettingIncomeFromTimesheetInTheMonth {
	return &incomeFromTimesheetRepository{session}
}

func NewIncomeFromTimesheetUserIncomeReader(session *mongo.Session) usecases.ForReadingIncomeFromTimesheetByUser {
	return &incomeFromTimesheetRepository{session}
}

// incomeFromTimesheetPeriodIndex guards the one-row-per-user-per-period invariant at the database
// level. The timesheet service delivers at-least-once, so the same summary can arrive more than
// once; the find-then-insert in the sync usecase cannot enforce uniqueness by itself.
func incomeFromTimesheetPeriodIndex() mongodriver.IndexModel {
	return mongodriver.IndexModel{
		Keys: bson.D{
			{Key: "userId", Value: 1},
			{Key: "year", Value: 1},
			{Key: "month", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("userId_year_month_unique"),
	}
}

// incomeFromTimesheetPeriodQuery keys on the period stored on the record rather than on a
// submitDate window like the hand-filled income collection does. A timesheet summary is synced
// after its month has closed, so the write timestamp is not in the period being reported and
// cannot identify the row to update — see models.IncomeFromTimesheet.
func incomeFromTimesheetPeriodQuery(userID string, year int, month time.Month) bson.M {
	return bson.M{
		"userId": userID,
		"year":   year,
		"month":  int(month),
	}
}

func (r *incomeFromTimesheetRepository) GetByUserYearMonth(userID string, year int, month time.Month) (*models.IncomeFromTimesheet, error) {
	query := incomeFromTimesheetPeriodQuery(userID, year, month)

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

// prepareIncomeFromTimesheetForInsert leaves SubmitDate alone: the sync usecase anchors it to the
// event's period, and overwriting it here would put every row in the month it happened to be
// synced instead of the month it reports on.
func prepareIncomeFromTimesheetForInsert(income *models.IncomeFromTimesheet) {
	income.LastUpdate = time.Now()
	income.ID = primitive.NewObjectID()
	income.ExportStatus = false
}

func (r *incomeFromTimesheetRepository) Add(income *models.IncomeFromTimesheet) error {
	prepareIncomeFromTimesheetForInsert(income)
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

func (r *incomeFromTimesheetRepository) GetByUserIdAllMonth(userId string) ([]*models.IncomeFromTimesheet, error) {
	records := make([]*models.IncomeFromTimesheet, 0)
	coll := r.session.GetCollection(incomeFromTimesheetColl)
	ctx := r.session.Ctx()
	cursor, err := coll.Find(ctx, bson.M{"userId": userId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}

func (r *incomeFromTimesheetRepository) GetAllByRoleStartDateAndEndDate(role string, startDate, endDate time.Time) ([]*models.IncomeFromTimesheet, error) {
	records := make([]*models.IncomeFromTimesheet, 0)
	coll := r.session.GetCollection(incomeFromTimesheetColl)
	ctx := r.session.Ctx()
	query := bson.M{
		"role": role,
		"submitDate": bson.M{
			"$gt": startDate,
			"$lt": endDate,
		},
	}
	cursor, err := coll.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &records); err != nil {
		return nil, err
	}
	return records, nil
}
