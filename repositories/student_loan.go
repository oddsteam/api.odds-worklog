package repositories

import (
	"time"

	"github.com/globalsign/mgo"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

const (
	studentLoanColl = "studentloan"
)

func (r *incomeRepository) GetStudentLoans() models.StudentLoanList {
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

func (r *incomeRepository) studentLoanCollection() *mgo.Collection {
	return r.session.GetCollection(studentLoanColl)
}
