package income

import (
	"testing"
	"time"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

func TestGetStudentLoansQueryWithCurrentMonth(t *testing.T) {
	n := time.Date(2022, time.Month(11), 1, 13, 30, 29, 0, time.UTC)
	query := loanQuery(n)
	assert.Equal(t, "11/2565", query["list.monthYear"])
}

func TestAdminCanExportIndividualIncomeWithoutStudentLoans(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("should not panic or it will block export income")
		}
	}()
	getStudentLoans(getLoanOfThisMonthWhichDoesNotExist)
}

func TestExportIncomeIgnoresStudentLoansCalculationWhenLoansDoesNotExist(t *testing.T) {
	loans := getStudentLoans(getLoanOfThisMonthWhichDoesNotExist)
	assert.Equal(t, 0, len(loans.List))
}

func getLoanOfThisMonthWhichDoesNotExist(result interface{}) error {
	return mgo.ErrNotFound
}
