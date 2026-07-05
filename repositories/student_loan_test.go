package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
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
	return mongo.ErrNoDocuments
}
