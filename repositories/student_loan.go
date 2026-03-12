package repositories

import (
	"time"

	"github.com/globalsign/mgo"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const (
	studentLoanColl = "studentloan"
)

type studentLoanRepository struct {
	session *mongo.Session
}

func NewStudentLoanRepository(session *mongo.Session) *studentLoanRepository {
	return &studentLoanRepository{session}
}

func (r *studentLoanRepository) GetStudentLoans() models.StudentLoanList {
	sll := models.StudentLoanList{}
	loanQuery := sll.GetFilterQuery(time.Now())
	return getStudentLoans(r.studentLoanCollection().Find(loanQuery).One)
}

func (r *studentLoanRepository) SaveStudentLoans(loanlist models.StudentLoanList) int {
	coll := r.studentLoanCollection()
	filter := loanlist.GetFilterQuery(time.Now())
	update := loanlist.GetUpdateQuery()
	changed, err := coll.Upsert(filter, update)
	if err != nil {
		panic(err.Error())
	}
	return changed.Matched
}

type getOneFn = func(result interface{}) (err error)

func getStudentLoans(getOneLoan getOneFn) models.StudentLoanList {
	loans := new(models.StudentLoanList)
	getOneLoan(loans)
	return *loans
}

func (r *studentLoanRepository) studentLoanCollection() *mgo.Collection {
	return r.session.GetCollection(studentLoanColl)
}
