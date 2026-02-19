package income

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const (
	incomeColl      = "income"
	studentLoanColl = "studentloan"
)

type repository struct {
	session *mongo.Session
}

func NewRepository(session *mongo.Session) Repository {
	return &repository{session}
}

func (r *repository) AddIncome(income *models.Income) error {
	t := time.Now()
	return r.AddIncomeOnSpecificTime(income, t)
}

func (r *repository) AddIncomeOnSpecificTime(income *models.Income, t time.Time) error {
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
	query := createQueryByIdAndPeriod(fromYear, fromMonth, id)
	return getIncomeByUserIDWithQuery(r, id, fromYear, fromMonth, query)
}

func (r *repository) GetIncomeByStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) (*models.Income, error) {

	query := createQueryByPeriod(startDate, endDate)
	return getIncomeByQuery(r, query)
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

func getIncomeByQuery(r *repository, query bson.M) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(query).One(&income)
	log.Println("getIncomeByQuery")
	log.Println(query)
	if err != nil {
		return nil, err
	}

	log.Println(income.SubmitDate)

	return income, nil
}

func getIncomeByUserIDWithQuery(r *repository, uID string, fromYear int, fromMonth time.Month, query bson.M) (*models.Income, error) {
	income := new(models.Income)
	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(query).One(&income)
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

func (r *repository) DropIncome() error {
	return r.session.GetCollection(incomeColl).DropCollection()
}

func (r *repository) UpdateExportStatus(id string) error {
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

func (r *repository) SaveStudentLoans(loanlist models.StudentLoanList) int {
	coll := r.studentLoanCollection()
	filter := loanlist.GetFilterQuery(time.Now())
	update := loanlist.GetUpdateQuery()
	changed, err := coll.Upsert(filter, update)
	if err != nil {
		panic(err.Error())
	}
	return changed.Matched
}

func (r *repository) studentLoanCollection() *mgo.Collection {
	return r.session.GetCollection(studentLoanColl)
}
