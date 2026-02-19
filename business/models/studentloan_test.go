package models_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	mock_user "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func TestFindLoanForUserUsingBankAccountName(t *testing.T) {
	sll := models.StudentLoanList{
		List: []models.StudentLoan{
			{ID: "", Fullname: "คนอื่น ที่ไม่ใช่", Amount: 943},
			{ID: "", Fullname: "ชื่อ นามสกุล", Amount: 1579},
		}}
	acc := mock_user.IndividualUser1.BankAccountName
	actual := sll.FindLoan(acc)
	assert.Equal(t, acc, actual.Fullname)
	assert.Equal(t, "1,579.00", actual.CSVAmount())
}

func TestFindLoanForUserWhoseBankAccountNameContainsTitle(t *testing.T) {
	sll := models.StudentLoanList{
		List: []models.StudentLoan{
			{ID: "", Fullname: "ชื่อ นามสกุล", Amount: 1579},
		}}
	acc := "นายชื่อ นามสกุล"
	actual := sll.FindLoan(acc)
	assert.Equal(t, "1,579.00", actual.CSVAmount())
}

func TestAmountIs0WhenCannotFindLoanForUser(t *testing.T) {
	sll := models.StudentLoanList{
		List: []models.StudentLoan{
			{ID: "", Fullname: "คนอื่น ที่ไม่ใช่", Amount: 943},
		}}
	acc := mock_user.IndividualUser1.BankAccountName
	actual := sll.FindLoan(acc)
	assert.Equal(t, "0.00", actual.CSVAmount())
}

func TestSaveStudenLoanWithIDSoItCanBeSavedSucessfully(t *testing.T) {
	studentLoanResponse := []byte(getMockStudentLoanResponseInJanuary())
	loanlist, err := models.CreateStudentLoanList(studentLoanResponse)
	assert.Equal(t, nil, err)

	loanlist.CreateIDForLoans()

	for _, l := range loanlist.List {
		assert.NotEmpty(t, l.ID)
	}
}

func TestSaveStudenLoanWithMonthYearAsItIsUsedLaterForQuery(t *testing.T) {
	studentLoanResponse := []byte(getMockStudentLoanResponseInJanuary())
	loanlist, err := models.CreateStudentLoanList(studentLoanResponse)
	assert.Equal(t, nil, err)

	filterQuery := loanlist.GetFilterQuery(time.Date(2023, time.Month(1), 1, 13, 30, 29, 0, time.UTC))

	assert.Equal(t, filterQuery["monthYear"], "1/2566")
}

func getMockStudentLoanResponseInJanuary() string {
	return `[
		{
		   "no":1,
		   "month":"01",
		   "year":"2566",
		   "refNo":"1310500176065",
		   "customerName":"ลอยด์ ฟอเจอร์",
		   "totalAmount":1000,
		   "paidAmount":1000,
		   "deleteFlag":"",
		   "deleteCause":"",
		   "paymentDate":"",
		   "monthYear":"01/2566"
		},
		{
		   "no":2,
		   "month":"01",
		   "year":"2566",
		   "refNo":"1309901199595",
		   "customerName":"อาเนีย ฟอเจอร์",
		   "totalAmount":908,
		   "paidAmount":908,
		   "deleteFlag":"",
		   "deleteCause":"",
		   "paymentDate":"",
		   "monthYear":"01/2566"
		}
	 ]`
}
