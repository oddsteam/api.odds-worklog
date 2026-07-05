package repositories

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	pkgmongo "gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

const (
	studentLoanColl = "studentloan"
)

type studentLoanRepository struct {
	session *pkgmongo.Session
}

func NewStudentLoanRepository(session *pkgmongo.Session) *studentLoanRepository {
	return &studentLoanRepository{session}
}

func (r *studentLoanRepository) GetStudentLoans() models.StudentLoanList {
	sll := models.StudentLoanList{}
	loanQuery := sll.GetFilterQuery(time.Now())
	ctx := r.session.Ctx()
	coll := r.studentLoanCollection()
	return getStudentLoans(func(result interface{}) error {
		return coll.FindOne(ctx, loanQuery).Decode(result)
	})
}

func (r *studentLoanRepository) SaveStudentLoans(loanlist models.StudentLoanList) int {
	coll := r.studentLoanCollection()
	ctx := r.session.Ctx()
	filter := loanlist.GetFilterQuery(time.Now())
	update := loanlist.GetUpdateQuery()
	opts := options.Update().SetUpsert(true)
	result, err := coll.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		panic(err.Error())
	}
	return int(result.MatchedCount)
}

type getOneFn = func(result interface{}) error

func getStudentLoans(getOneLoan getOneFn) models.StudentLoanList {
	loans := new(models.StudentLoanList)
	err := getOneLoan(loans)
	if err == mongo.ErrNoDocuments {
		return models.StudentLoanList{}
	}
	if err != nil {
		panic(err.Error())
	}
	return *loans
}

func (r *studentLoanRepository) studentLoanCollection() *mongo.Collection {
	return r.session.GetCollection(studentLoanColl)
}
