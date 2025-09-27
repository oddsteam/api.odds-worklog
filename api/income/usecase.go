package income

import (
	"encoding/csv"
	"errors"
	"strings"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

type usecase struct {
	repo     Repository
	userRepo user.Repository
}

func NewUsecase(r Repository, ur user.Repository) Usecase {
	return &usecase{r, ur}
}

func (u *usecase) AddIncome(req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
	year, month := utils.GetYearMonthNow()
	_, err := u.repo.GetIncomeUserByYearMonth(uid, year, month)
	if err == nil {
		return nil, errors.New("Sorry, has income data of user " + userDetail.GetName())
	}
	income, err := NewIncome(uid).prepareDataForAddIncome(*req, *userDetail)
	if err != nil {
		return nil, err
	}
	err = u.repo.AddIncome(income)
	if err != nil {
		return nil, err
	}

	return income, nil
}

func (u *usecase) UpdateIncome(id string, req *models.IncomeReq, uid string) (*models.Income, error) {
	userDetail, _ := u.userRepo.GetByID(uid)
	income, err := u.repo.GetIncomeByID(id, uid)
	if err != nil {
		return nil, err
	}

	err = NewIncome(uid).prepareDataForUpdateIncome(*req, *userDetail, income)
	if err != nil {
		return nil, err
	}
	u.repo.UpdateIncome(income)

	return income, nil
}

func (u *usecase) GetIncomeStatusList(role string, isAdmin bool) ([]*models.IncomeStatus, error) {
	var incomeList []*models.IncomeStatus
	users, err := u.userRepo.GetByRole(role)
	if err != nil {
		return nil, err
	}

	year, month := utils.GetYearMonthNow()
	for index, element := range users {
		element.ThaiCitizenID = ""
		element.DailyIncome = ""

		incomeUser, err := u.repo.GetIncomeUserByYearMonth(element.ID.Hex(), year, month)
		income := models.IncomeStatus{User: element}
		incomeList = append(incomeList, &income)
		if !isAdmin {
			element.ID = ""
		}
		if err != nil {
			incomeList[index].Status = "N"
		} else {
			incomeList[index].WorkDate = incomeUser.WorkDate
			incomeList[index].WorkingHours = incomeUser.WorkingHours
			incomeList[index].SubmitDate = incomeUser.SubmitDate.Format(time.RFC3339)
			incomeList[index].Status = "Y"
		}
	}
	return incomeList, nil
}

func (u *usecase) GetIncomeByUserIdAndCurrentMonth(userId string) (*models.Income, error) {
	year, month := utils.GetYearMonthNow()
	return u.repo.GetIncomeUserByYearMonth(userId, year, month)
}

func (u *usecase) GetIncomeByUserIdAllMonth(userId string) ([]*models.Income, error) {
	listIncome, err := u.repo.GetIncomeByUserIdAllMonth(userId)
	if err != nil {
		return nil, err
	}
	if len(listIncome) == 0 {
		return nil, nil
	}
	for index := range listIncome {
		if listIncome[index].NetSpecialIncome != "" && listIncome[index].NetDailyIncome != "" {
			listIncome[index].NetIncome, err = calTotal(listIncome[index].NetDailyIncome, listIncome[index].NetSpecialIncome)
			if err != nil {
				return nil, err
			}
		}
	}
	return listIncome, nil
}

func (u *usecase) ExportIncome(role string, beforeMonth string) (string, error) {
	return u.exportIncome(role, beforeMonth)
}

func (u *usecase) exportIncome(role string, beforeMonth string) (string, error) {
	var t time.Time
	if beforeMonth == "0" {
		t = time.Now()
	} else {
		t = time.Now().AddDate(0, -1, 0)
	}
	startDate, endDate := utils.GetStartDateAndEndDate(t)
	return u.ExportIncomeByStartDateAndEndDate(role, startDate, endDate)
}

func (u *usecase) ExportIncomeByStartDateAndEndDate(role string, startDate, endDate time.Time) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	defer file.Close()

	if err != nil {
		return "", err
	}

	incomes, err := u.repo.GetAllIncomeByRoleStartDateAndEndDate(role, startDate, endDate)

	if err != nil {
		return "", err
	}

	studentLoanList := u.repo.GetStudentLoans()

	ics := NewIncomes(incomes, studentLoanList)
	strWrite, _ := ics.toCSV()

	if len(strWrite) == 1 {
		return "", errors.New("no data for export to CSV file")
	}

	csvWriter := csv.NewWriter(file)
	csvWriter.WriteAll(strWrite)
	csvWriter.Flush()

	ep := models.Export{
		Filename: filename,
		Date:     time.Now(),
	}
	err = u.repo.AddExport(&ep)
	if err != nil {
		return "", err
	}

	return filename, nil
}

func (u *usecase) ExportIncomeSAP(role string, beforeMonth string, dateEff time.Time) (string, error) {
	var t time.Time
	if beforeMonth == "0" {
		t = time.Now()
	} else {
		t = time.Now().AddDate(0, -1, 0)
	}
	startDate, endDate := utils.GetStartDateAndEndDate(t)
	return u.ExportIncomeSAPByStartDateAndEndDate(role, startDate, endDate, dateEff)
}

func (u *usecase) ExportIncomeSAPByStartDateAndEndDate(role string, startDate, endDate time.Time, dateEff time.Time) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	encoder := charmap.Windows874.NewEncoder()
	writer := transform.NewWriter(file, encoder)
	defer file.Close()
	defer writer.Close()

	if err != nil {
		return "", err
	}

	incomes, err := u.repo.GetAllIncomeByRoleStartDateAndEndDate(role, startDate, endDate)

	if err != nil {
		return "", err
	}

	studentLoanList := u.repo.GetStudentLoans()

	ics := NewIncomes(incomes, studentLoanList)

	strWrite, _ := ics.toSAP(dateEff)

	if len(strWrite) == 0 {
		return "", errors.New("no data for export to SAP file")
	}

	for _, record := range strWrite {
		_, err := writer.Write([]byte(strings.Join(record, "") + "\n"))
		if err != nil {
			return "", err
		}
	}

	ep := models.Export{
		Filename: filename,
		Date:     time.Now(),
	}
	err = u.repo.AddExport(&ep)
	if err != nil {
		return "", err
	}

	return filename, nil
}

/** deprecated **/
func createRow(record models.Income, user models.User, loan models.StudentLoan) []string {
	i := NewIncomeFromRecord(record)
	i.SetLoan(&loan)
	return i.export(user)
}

func (u *usecase) GetAllInComeByStartDateAndEndDate(userIds []string, startDate time.Time, endDate time.Time) ([]*models.Income, error) {
	return u.repo.GetAllIncomeByStartDateAndEndDate(userIds, startDate, endDate)
}

func (u *usecase) GetByRole(role string) ([]*models.User, error) {
	return u.userRepo.GetByRole(role)
}

func (u *usecase) GetUserByID(userId string) (*models.User, error) {
	return u.userRepo.GetByID(userId)
}

func calTotal(main string, special string) (string, error) {
	ma, err := utils.StringToFloat64(main)
	if err != nil {
		return "", err
	}
	sp, err := utils.StringToFloat64(special)
	if err != nil {
		return "", err
	}
	return utils.FloatToString(ma + sp), nil
}

const (
	VENDOR_CODE_INDEX = iota
	ACCOUNT_NAME_INDEX
	PAYMENT_METHOD_INDEX
	ACCOUNT_NUMBER_INDEX
	NAME_INDEX
	ID_CARD_INDEX
	EMAIL_INDEX
	NET_DAILY_INCOME_INDEX
	NET_SPECIAL_INCOME_INDEX
	LOAN_DEDUCTION_INDEX
	WITHHOLDING_TAX_INDEX
	TRANSFER_AMOUNT_INDEX
	NOTE_INDEX
	SUBMIT_DATE_INDEX
)

func createHeaders() []string {
	return []string{"Vendor Code", "ชื่อบัญชี", "Payment method", "เลขบัญชี", "ชื่อ", "เลขบัตรประชาชน", "อีเมล", "จำนวนเงินรายได้หลัก", "จำนวนรายได้พิเศษ", "กยศและอื่น ๆ", "หัก ณ ที่จ่าย", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก"}
}
