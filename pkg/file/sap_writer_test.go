package file

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/entity"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestSAPWriter(t *testing.T) {
	const (
		SAP_TXN_INDEX             = 0
		SAP_PAYER_NAME_INDEX      = 1
		SAP_PAYEE_NAME_INDEX      = 2
		SAP_MALE_TO_NAME_INDEX    = 3
		SAP_BENEFICIARY1_INDEX    = 4
		SAP_BENEFICIARY2_INDEX    = 5
		SAP_BENEFICIARY3_INDEX    = 6
		SAP_BENEFICIARY4_INDEX    = 7
		SAP_ZIPCODE_INDEX         = 8
		SAP_CUSTOMER_REF_INDEX    = 9
		SAP_DATE_EFFECTIVE_INDEX  = 10
		SAP_DATE_PICKUP_INDEX     = 11
		SAP_CURRENCY_INDEX        = 12
		SAP_EMPTY_1_INDEX         = 13
		SAP_COMPANY_ACCNO_INDEX   = 14
		SAP_AMOUNT_INDEX          = 15
		SAP_PAYEE_BANK_CODE_INDEX = 16
		SAP_ACCOUNTNO_INDEX       = 17
		SAP_UNKNOW_1_INDEX        = 18
		SAP_EMPTY_2_INDEX         = 19
		SAP_EMPTY_3_INDEX         = 20
		SAP_ADVICEMODE2_INDEX     = 21
		SAP_FAXNO_INDEX           = 22
		SAP_EMAIL_INDEX           = 23
		SAP_SMSNO_INDEX           = 24
		SAP_CHARGE_ON_INDEX       = 25
		SAP_PRODUCT_INDEX         = 26
		SAP_SCHEDULE_INDEX        = 27
		SAP_EMPTY_4_INDEX         = 28
		SAP_DOCREQ_INDEX          = 29
		SAP_EMPTY_5_INDEX         = 30
		SAP_END_INDEX             = 31

		SAP_WHT_WHT_INDEX      = 0
		SAP_WHT_EMPTY_1_INDEX  = 1
		SAP_WHT_TAX_ID_INDEX   = 2
		SAP_WHT_EMPTY_3_INDEX  = 3
		SAP_WHT_EMPTY_4_INDEX  = 4
		SAP_WHT_EMPTY_5_INDEX  = 5
		SAP_WHT_EMPTY_6_INDEX  = 6
		SAP_WHT_EMPTY_7_INDEX  = 7
		SAP_WHT_EMPTY_8_INDEX  = 8
		SAP_WHT_EMPTY_9_INDEX  = 9
		SAP_WHT_EMPTY_10_INDEX = 10
		SAP_WHT_EMPTY_11_INDEX = 11
		SAP_WHT_EMPTY_12_INDEX = 12
		SAP_WHT_EMPTY_13_INDEX = 13
		SAP_WHT_EMPTY_14_INDEX = 14
		SAP_WHT_COM_NAME       = 15
		SAP_WHT_ADDRESS        = 16
		SAP_WHT_EMPTY_17       = 17
		SAP_WHT_EMPTY_18       = 18
		SAP_WHT_EMPTY_19       = 19
		SAP_WHT_EMPTY_20       = 20
	)

	t.Run("test export to SAP transaction should format correctly", func(t *testing.T) {
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		i := entity.MockSoloCorporateIncome

		txn, _ := exportSAP(*entity.NewPayrollFromIncome(i), dateEff)

		assert.Equal(t, "TXN", txn[SAP_TXN_INDEX])
		assert.Equal(t, "บจก. ออด-อี (ประเทศไทย) จำกัด                                                                                           ", txn[SAP_PAYER_NAME_INDEX])
		assert.Equal(t, "บจก. โซโล่ เลเวลลิ่ง                                                                                                              ", txn[SAP_PAYEE_NAME_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_MALE_TO_NAME_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_BENEFICIARY1_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_BENEFICIARY2_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_BENEFICIARY3_INDEX])
		assert.Equal(t, "                                        ", txn[SAP_BENEFICIARY4_INDEX])
		assert.Equal(t, "          ", txn[SAP_ZIPCODE_INDEX])
		assert.Equal(t, "                ", txn[SAP_CUSTOMER_REF_INDEX])
		assert.Equal(t, "29092025", txn[SAP_DATE_EFFECTIVE_INDEX])
		assert.Equal(t, "29092025", txn[SAP_DATE_PICKUP_INDEX])
		assert.Equal(t, "THB", txn[SAP_CURRENCY_INDEX])
		assert.Equal(t, "                                                  ", txn[SAP_EMPTY_1_INDEX])
		assert.Equal(t, "00011595873         ", txn[SAP_COMPANY_ACCNO_INDEX])
		assert.Equal(t, "000000054704.00", txn[SAP_AMOUNT_INDEX])
		assert.Equal(t, "0110246         ", txn[SAP_PAYEE_BANK_CODE_INDEX])
		assert.Equal(t, "02462737202         ", txn[SAP_ACCOUNTNO_INDEX])
		assert.Equal(t, "0400", txn[SAP_UNKNOW_1_INDEX])
		assert.Equal(t, "  ", txn[SAP_EMPTY_2_INDEX])
		assert.Equal(t, "                    ", txn[SAP_EMPTY_3_INDEX])
		assert.Equal(t, "     ", txn[SAP_ADVICEMODE2_INDEX])
		assert.Equal(t, "                                                  ", txn[SAP_FAXNO_INDEX])
		assert.Equal(t, "                                                  ", txn[SAP_EMAIL_INDEX])
		assert.Equal(t, "                                                  ", txn[SAP_SMSNO_INDEX])
		assert.Equal(t, "OUR          ", txn[SAP_CHARGE_ON_INDEX])
		assert.Equal(t, "DCR", txn[SAP_PRODUCT_INDEX])
		assert.Equal(t, "     ", txn[SAP_SCHEDULE_INDEX])
		assert.Equal(t, "                                  ", txn[SAP_EMPTY_4_INDEX])
		assert.Equal(t, "                                                                                                         ", txn[SAP_DOCREQ_INDEX])
		assert.Equal(t, "                                                                                                                                                                                                                                                                                                       ", txn[SAP_EMPTY_5_INDEX])
		assert.Equal(t, "END", txn[SAP_END_INDEX])

	})

	// t.Run("export จำนวนเงินที่ต้องโอนสำหรับ individual in fomat CSV และ SAP ควรตรงกัน", func(t *testing.T) {
	// 	uidFromSession := "5bbcf2f90fd2df527bc39539"
	// 	dailyIncome := "2000"
	// 	workDate := "16.5"
	// 	specialIncome := "250"
	// 	workingHours := "128.45"
	// 	u := entity.GivenIndividualUser(uidFromSession, dailyIncome)
	// 	req := entity.IncomeReq{
	// 		WorkDate:      workDate,
	// 		SpecialIncome: specialIncome,
	// 		WorkingHours:  workingHours,
	// 	}
	// 	record := entity.CreateIncome(u, req, "note")
	// 	i := entity.NewIncomeFromRecord(*record)

	// 	csvColumns := i.export()
	// 	dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
	// 	txn, _ := exportSAP(*i, dateEff)

	// 	assert.Equal(t, "1953.38", csvColumns[WITHHOLDING_TAX_INDEX])
	// 	assert.Equal(t, "63,159.12", csvColumns[TRANSFER_AMOUNT_INDEX])
	// 	assert.Equal(t, "000000063159.12", txn[SAP_AMOUNT_INDEX])
	// })

	t.Run("test export to SAP wht should format correctly", func(t *testing.T) {
		dateEff := time.Date(2025, 9, 29, 0, 0, 0, 0, time.UTC)
		i := entity.MockSoloCorporateIncome

		_, wht := exportSAP(*entity.NewPayrollFromIncome(i), dateEff)

		assert.Equal(t, "WHT", wht[SAP_WHT_WHT_INDEX])
		assert.Equal(t, "             ", wht[SAP_WHT_EMPTY_1_INDEX])
		assert.Equal(t, "0105556110718", wht[SAP_WHT_TAX_ID_INDEX])
		assert.Equal(t, "  ", wht[SAP_WHT_EMPTY_3_INDEX])
		assert.Equal(t, "000000000000.00", wht[SAP_WHT_EMPTY_4_INDEX])
		assert.Equal(t, "  ", wht[SAP_WHT_EMPTY_5_INDEX])
		assert.Equal(t, "                                   ", wht[SAP_WHT_EMPTY_6_INDEX])
		assert.Equal(t, "     ", wht[SAP_WHT_EMPTY_7_INDEX])
		assert.Equal(t, "000000000000.00", wht[SAP_WHT_EMPTY_8_INDEX])
		assert.Equal(t, "000000000000.00", wht[SAP_WHT_EMPTY_9_INDEX])
		assert.Equal(t, "  ", wht[SAP_WHT_EMPTY_10_INDEX])
		assert.Equal(t, "                                   ", wht[SAP_WHT_EMPTY_11_INDEX])
		assert.Equal(t, "     ", wht[SAP_WHT_EMPTY_12_INDEX])
		assert.Equal(t, "000000000000.00", wht[SAP_WHT_EMPTY_13_INDEX])
		assert.Equal(t, "                                                                                                                                                ", wht[SAP_WHT_EMPTY_14_INDEX])
		assert.Equal(t, "บจก. ออด-อี (ประเทศไทย) จำกัด                                                                                           ", wht[SAP_WHT_COM_NAME])
		assert.Equal(t, "2549/41-43 พหลโยธิน ลาดยาว จตุจักร กรุงเทพ 10900                                                                                                                ", wht[SAP_WHT_ADDRESS])
		assert.Equal(t, "                                                                                                                        ", wht[SAP_WHT_EMPTY_17])
		assert.Equal(t, "                                                                                                                                                                ", wht[SAP_WHT_EMPTY_18])
		assert.Equal(t, "                    ", wht[SAP_WHT_EMPTY_19])
		assert.Equal(t, utils.AddBlank("", 938), wht[SAP_WHT_EMPTY_20])

	})

}
