package income

import (
	"encoding/csv"
	"errors"
	"fmt"
	"time"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	repo     Repository
	userRepo user.Repository
}

type incomeSum struct {
	Net         string
	VAT         string
	WHT         string
	TotalIncome string
}

func NewUsecase(r Repository, ur user.Repository) Usecase {
	return &usecase{r, ur}
}

func calIncomeSum(workAmount string, vattype string, incomes string, role string) (*incomeSum, error) {
	var vat, wht string
	var vatf, whtf float64
	var ins = new(incomeSum)

	amount, _ := utils.StringToFloat64(workAmount)
	income, _ := utils.StringToFloat64(incomes)
	totalIncomeStr := utils.FloatToString(amount * income)
	totalIncome, err := utils.StringToFloat64(totalIncomeStr)
	if err != nil {
		return nil, err
	}
	if role == "corporate" {
		wht, whtf, err = calWHTCorporate(totalIncomeStr)
		if err != nil {
			return nil, err
		}
	}
	if role == "individual" {
		wht, whtf, err = calWHT(totalIncomeStr)
		if err != nil {
			return nil, err
		}
	}
	ins.WHT = wht
	ins.TotalIncome = totalIncomeStr

	if vattype == "Y" {
		vat, vatf, err = calVAT(totalIncomeStr)
		if err != nil {
			return nil, err
		}

		net := totalIncome + vatf - whtf

		ins.Net = utils.FloatToString(net)
		ins.VAT = vat
		return ins, nil
	}
	net := totalIncome - whtf
	ins.Net = utils.FloatToString(net)
	return ins, nil
}

func calWHTCorporate(income string) (string, float64, error) {
	return calWHT(income)
}

func calWHT(income string) (string, float64, error) {
	return multiply(income, 0.03)
}

func calVAT(income string) (string, float64, error) {
	return multiply(income, 0.07)
}

func multiply(a string, b float64) (string, float64, error) {
	num, err := utils.StringToFloat64(a)
	if err != nil {
		return "", 0.0, err
	}
	vat := num * b
	return utils.FloatToString(vat), utils.RealFloat(vat), nil
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
			listIncome[index].NetIncome, err = calSummary(listIncome[index].NetDailyIncome, listIncome[index].NetSpecialIncome)
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

func (u *usecase) ExportIncomeNotExport(role string) (string, error) {
	year, month := utils.GetYearMonthNow()
	getIncome := u.createFunctionGetUnexportedIncomeByUserWithPeriod(year, month)
	shouldUpdateExportStatus := true
	return u.exportIncome(role, getIncome, shouldUpdateExportStatus)
}

func (u *usecase) exportIncome(role string, getIncome getIncomeFn, shouldUpdateExportStatus bool) (string, error) {
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
	fmt.Printf("%#v", studentLoanList)

	strWrite := make([][]string, 0)
	strWrite = append(strWrite, createHeaders())
	for _, user := range users {
		income, err := getIncome(*user)
		if err == nil {
			if shouldUpdateExportStatus {
				u.repo.UpdateExportStatus(income.ID.Hex())
			}
			d := createRow(*income, *user)
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

func createRow(income models.Income, user models.User) []string {
	t := income.SubmitDate
	summaryIncome, _ := calSummary(income.NetDailyIncome, income.NetSpecialIncome)
	tf := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
	d := []string{
		user.GetName(),
		user.BankAccountName,
		setValueCSV(user.BankAccountNumber),
		user.Email,
		setValueCSV(utils.FormatCommas(income.NetDailyIncome)),
		setValueCSV(utils.FormatCommas(income.NetSpecialIncome)),
		setValueCSV("0"),
		setValueCSV(utils.FormatCommas(summaryIncome)),
		income.Note,
		tf,
	}
	return d
}

func calSummary(main string, special string) (string, error) {
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

func setValueCSV(s string) string {
	return utils.SetValueCSV(s)
}

func createHeaders() []string {
	return []string{"ชื่อ", "ชื่อบัญชี", "เลขบัญชี", "อีเมล", "จำนวนเงินรายได้หลัก", "จำนวนรายได้พิเศษ", "กยศและอื่น ๆ", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก"}
}
