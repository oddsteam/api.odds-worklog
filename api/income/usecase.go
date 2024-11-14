package income

import (
	"encoding/csv"
	"errors"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
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
	beforemonth, err := utils.StringToInt(beforeMonth)
	if err != nil {
		return "", err
	}
	year, month := utils.GetYearMonthNow()
	getIncome := u.createFunctionGetIncomeByUserWithPeriod(year, month-time.Month(beforemonth))
	shouldUpdateExportStatus := beforeMonth == "0"

	return u.exportIncome(role, getIncome, shouldUpdateExportStatus)
}

func (u *usecase) ExportIncomeNew(role string, beforeMonth string) (string, error) {
	shouldUpdateExportStatus := beforeMonth == "0"

	return u.exportIncome_new(role, shouldUpdateExportStatus)
}

func (u *usecase) ExportIncomeNotExport(role string) (string, error) {
	year, month := utils.GetYearMonthNow()
	getIncome := u.createFunctionGetUnexportedIncomeByUserWithPeriod(year, month)
	shouldUpdateExportStatus := true
	return u.exportIncome(role, getIncome, shouldUpdateExportStatus)
}

func (u *usecase) exportCsvByInCome(role string, incomes []*models.Income) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	defer file.Close()

	if err != nil {
		return "", err
	}

	studentLoanList := u.repo.GetStudentLoans()

	strWrite := make([][]string, 0)
	strWrite = append(strWrite, createHeaders())

	for _, income := range incomes {
		user, err := u.GetUserByID(income.UserID)
		loan := studentLoanList.FindLoan(user.BankAccountName)
		if err == nil {
			d := createRow(*income, *user, loan)
			strWrite = append(strWrite, d)
		}
	}

	if len(strWrite) == 1 {
		return "", errors.New("No data for export to CSV file.")
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

func (u *usecase) exportIncome(role string, getIncome getIncomeFn, shouldUpdateExportStatus bool) (string, error) {
	return u.exportIncome_obsoleted(role, getIncome, shouldUpdateExportStatus)
}

func (u *usecase) exportIncome_new(role string, shouldUpdateExportStatus bool) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	defer file.Close()

	if err != nil {
		return "", err
	}

	startDate, endDate := utils.GetStartDateAndEndDate(time.Now())
	incomes, err := u.repo.GetAllIncomeByRoleStartDateAndEndDate(role, startDate, endDate)

	if err != nil {
		return "", err
	}

	studentLoanList := u.repo.GetStudentLoans()

	ics := NewIncomes(incomes, studentLoanList)
	strWrite, updatedIncomeIds := ics.toCSV()

	for _, id := range updatedIncomeIds {
		if shouldUpdateExportStatus {
			u.repo.UpdateExportStatus(id)
		}
	}
	if len(strWrite) == 1 {
		return "", errors.New("No data for export to CSV file.")
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

func (u *usecase) exportIncome_obsoleted(role string, getIncome getIncomeFn, shouldUpdateExportStatus bool) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	defer file.Close()

	if err != nil {
		return "", err
	}
	users, err := u.userRepo.GetByRole(role)
	if err != nil {
		return "", err
	}

	studentLoanList := u.repo.GetStudentLoans()

	strWrite := make([][]string, 0)
	strWrite = append(strWrite, createHeaders())
	for _, user := range users {
		income, err := getIncome(*user)
		loan := studentLoanList.FindLoan(user.BankAccountName)
		if err == nil {
			if shouldUpdateExportStatus {
				u.repo.UpdateExportStatus(income.ID.Hex())
			}
			d := createRow(*income, *user, loan)
			strWrite = append(strWrite, d)
		}
	}

	if len(strWrite) == 1 {
		return "", errors.New("No data for export to CSV file.")
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

type getIncomeFn = func(user models.User) (*models.Income, error)

func (u *usecase) createFunctionGetUnexportedIncomeByUserWithPeriod(year int, month time.Month) getIncomeFn {
	return func(user models.User) (*models.Income, error) {
		return u.repo.GetIncomeByUserID(user.ID.Hex(), year, month)
	}
}

func (u *usecase) createFunctionGetIncomeByUserWithPeriod(year int, month time.Month) getIncomeFn {
	return func(user models.User) (*models.Income, error) {
		return u.repo.GetIncomeUserByYearMonth(user.ID.Hex(), year, month)
	}
}

func createRow(record models.Income, user models.User, loan models.StudentLoan) []string {
	i := NewIncomeFromRecord(record)
	i.SetLoan(&loan)
	return i.export(user)
}

func (u *usecase) ExportIncomeByStartDateAndEndDate(role string, incomes []*models.Income) (string, error) {
	return u.exportCsvByInCome(role, incomes)
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
	ACCOUNT_NUMBER_INDEX
	PAYMENT_METHOD_INDEX
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
	return []string{"Vendor Code", "ชื่อบัญชี", "เลขบัญชี", "Payment method", "ชื่อ", "เลขบัตรประชาชน", "อีเมล", "จำนวนเงินรายได้หลัก", "จำนวนรายได้พิเศษ", "กยศและอื่น ๆ", "หัก ณ ที่จ่าย", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก"}
}
