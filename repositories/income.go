package repositories

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const (
	incomeColl = "income"
	exportColl = "export"
)

type incomeRepository struct {
	session *mongo.Session
}

func NewUserIncomeReader(session *mongo.Session) usecases.ForReadingUserIncome {
	return &incomeRepository{session}
}

func NewUserIncomeWriter(session *mongo.Session) usecases.ForControllingUserIncome {
	return &incomeRepository{session}
}

func NewUserIncomeUpdater(session *mongo.Session) usecases.ForUpdatingUserIncome {
	return &incomeRepository{session}
}

func NewIncomeReader(session *mongo.Session) usecases.ForGettingIncomeData {
	return &incomeRepository{session}
}

func NewIncomeWriter(session *mongo.Session) usecases.ForControllingIncomeData {
	return &incomeRepository{session}
}

func (r *incomeRepository) AddIncome(income *models.Income) error {
	t := time.Now()
	return r.AddIncomeOnSpecificTime(income, t)
}

func (r *incomeRepository) AddIncomeOnSpecificTime(income *models.Income, t time.Time) error {
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

func (r *incomeRepository) GetIncomeUserByYearMonth(id string, fromYear int, fromMonth time.Month) (*models.Income, error) {
	query := createQueryByIdAndPeriod(fromYear, fromMonth, id)
	return getIncomeByUserIDWithQuery(r, id, fromYear, fromMonth, query)
}

func (r *incomeRepository) GetIncomeByStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) (*models.Income, error) {
	query := createQueryByPeriod(startDate, endDate)
	return getIncomeByQuery(r, query)
}

func (r *incomeRepository) GetIncomeByUserIdAllMonth(id string) ([]*models.Income, error) {
	income := make([]*models.Income, 0)

	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(bson.M{"userId": id}).All(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (r *incomeRepository) GetIncomeByID(incID, uID string) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(incID), "userId": uID}).One(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (r *incomeRepository) GetIncomeByUserID(uID string, fromYear int, fromMonth time.Month) (*models.Income, error) {
	query := createQueryByIdAndPeriod(fromYear, fromMonth, uID)
	query["exportStatus"] = false
	return getIncomeByUserIDWithQuery(r, uID, fromYear, fromMonth, query)
}

func createQueryByIdAndPeriod(fromYear int, fromMonth time.Month, uID string) bson.M {
	fromDate := time.Date(fromYear, fromMonth, 1, 0, 0, 0, 0, time.UTC)
	toDate := fromDate.AddDate(0, 1, 0)

	query := bson.M{
		"userId": uID,
		"submitDate": bson.M{
			"$gt": fromDate,
			"$lt": toDate,
		},
	}
	return query
}

func createQueryByPeriod(startDate time.Time, endDate time.Time) bson.M {
	query := bson.M{
		"exportStatus": true,
		"submitDate": bson.M{
			"$gt": startDate,
			"$lt": endDate,
		},
	}
	return query
}

func getIncomeByQuery(r *incomeRepository, query bson.M) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(query).One(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func getIncomeByUserIDWithQuery(r *incomeRepository, uID string, fromYear int, fromMonth time.Month, query bson.M) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(query).One(&income)
	if err != nil {
		return nil, err
	}
	return income, nil
}

func (r *incomeRepository) UpdateIncome(income *models.Income) error {
	income.LastUpdate = time.Now()
	income.ExportStatus = false
	coll := r.session.GetCollection(incomeColl)
	err := coll.UpdateId(income.ID, &income)
	if err != nil {
		return err
	}
	return nil
}

func (r *incomeRepository) DropIncome() error {
	return r.session.GetCollection(incomeColl).DropCollection()
}

func (r *incomeRepository) UpdateExportStatus(id string) error {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(bson.M{"_id": bson.ObjectIdHex(id)}).One(&income)

	if err != nil {
		return err
	}

	err = coll.Update(bson.M{"_id": income.ID}, bson.M{"$set": bson.M{"exportStatus": true}})
	if err != nil {
		return err
	}
	return nil
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
