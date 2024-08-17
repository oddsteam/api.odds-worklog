package income

import (
	"log"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const (
	incomeColl      = "income"
	exportColl      = "export"
	studentLoanColl = "studentloan"
	userColl        = "user"
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
	query := createQueryByIdAndPeriod(fromYear, fromMonth, id)
	return getIncomeByUserIDWithQuery(r, id, fromYear, fromMonth, query)
}

func (r *repository) GetIncomeByStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) (*models.Income, error) {

	query := createQueryByPeriod(startDate, endDate)
	return getIncomeByQuery(r, query)
}

func (r *repository) GetAllIncomeByStartDateAndEndDate(userIds []string, startDate time.Time, endDate time.Time) ([]*models.Income, error) {

	query := createQueryIncomeByStartDateAndEndDate(userIds, startDate, endDate)
	return getAllInComeByQuery(r, query)
}

func (r *repository) GetAllIncomeByRoleStartDateAndEndDate(role string, startDate time.Time, endDate time.Time) ([]*models.Income, error) {
	query := createQueryIncomeByRoleStartDateAndEndDate(role, startDate, endDate)
	return getAllInComeByQuery(r, query)
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

func createQueryGetUserByRole(role string) bson.M {

	query := bson.M{
		"role": role,
	}
	return query
}

func createQueryIncomeByStartDateAndEndDate(userIds []string, startDate time.Time, endDate time.Time) bson.M {
	query := bson.M{
		"userId": bson.M{
			"$in": userIds,
		},
		"submitDate": bson.M{
			"$gt": startDate,
			"$lt": endDate,
		},
	}
	return query
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

func getAllInComeByQuery(r *repository, query bson.M) ([]*models.Income, error) {
	incomes := make([]*models.Income, 0)

	coll := r.session.GetCollection(incomeColl)
	err := coll.Find(query).All(&incomes)
	if err != nil {
		return nil, err
	}
	return incomes, nil
}

//func getAllUserByQuery(r *repository, query bson.M) ([]string,error){
//
//
//}
//GetAllUserIdByRole

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

func (r *repository) AddExport(ep *models.Export) error {
	coll := r.session.GetCollection(exportColl)
	ep.ID = bson.NewObjectId()
	return coll.Insert(ep)
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

func (r *repository) GetStudentLoans() models.StudentLoanList {
	sll := models.StudentLoanList{}
	loanQuery := sll.GetFilterQuery(time.Now())
	return getStudentLoans(r.studentLoanCollection().Find(loanQuery).One)
}

type getOneFn = func(result interface{}) (err error)

func getStudentLoans(getOneLoan getOneFn) models.StudentLoanList {
	loans := new(models.StudentLoanList)
	getOneLoan(loans)
	return *loans
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
