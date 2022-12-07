package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	mock_user "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestFindLoanForUserUsingBankAccountName(t *testing.T) {
	sll := models.StudentLoanList{
		List: []models.StudentLoan{
			{ID: "", Fullname: "คนอื่น ที่ไม่ใช่", Amount: 943},
			{ID: "", Fullname: "ชื่อ นามสกุล", Amount: 1579},
		}}
	u := mock_user.IndividualUser1
	actual := sll.FindLoan(u)
	assert.Equal(t, u.BankAccountName, actual.Fullname)
	assert.Equal(t, `="1,579.00"`, actual.CSVAmount())
}

func TestAmountIs0WhenCannotFindLoanForUser(t *testing.T) {
	sll := models.StudentLoanList{
		List: []models.StudentLoan{
			{ID: "", Fullname: "คนอื่น ที่ไม่ใช่", Amount: 943},
		}}
	u := mock_user.IndividualUser1
	actual := sll.FindLoan(u)
	assert.Equal(t, `="0.00"`, actual.CSVAmount())
}
