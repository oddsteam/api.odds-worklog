package income

import (
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

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
