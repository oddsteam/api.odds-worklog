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

func calVAT(income string) (string, float64, error) {
	num, err := utils.StringToFloat64(income)
	if err != nil {
		return "", 0.0, err
	}
	vat := num * 0.07
	return utils.FloatToString(vat), utils.RealFloat(vat), nil
}

func calWHT(income string) (string, float64, error) {
	num, err := utils.StringToFloat64(income)
	if err != nil {
		return "", 0.0, err
	}
	wht := num * 0.03
	return utils.FloatToString(wht), utils.RealFloat(wht), nil
}

func calIncomeSum(workAmount string, vattype string, incomes string) (*incomeSum, error) {
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
	wht, whtf, err = calWHT(totalIncomeStr)
	if err != nil {
		return nil, err
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

func (u *usecase) ExportIncome(role string, beforeMonth string) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	defer file.Close()

	if err != nil {
		return "", err
	}
	users, err := u.userRepo.GetByRole(role)
	if err != nil {
		return "", err
	}
	year, month := utils.GetYearMonthNow()
	beforemonth, err := utils.StringToInt(beforeMonth)
	if err != nil {
		return "", err
	}

	strWrite := make([][]string, 0)
	d := []string{"ชื่อ", "ชื่อบัญชี", "เลขบัญชี", "จำนวนเงินรายได้หลัก", "จำนวนรายได้พิเศษ", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก"}
	strWrite = append(strWrite, d)
	for _, user := range users {
		income, err := u.repo.GetIncomeUserByYearMonth(user.ID.Hex(), year, month-time.Month(beforemonth))
		if err == nil {
			if beforeMonth == "0" {
				u.repo.UpdateExportStatus(income.UserID)
			}
			t := income.SubmitDate
			summaryIncome, _ := calSummary(income.NetIncome, income.NetSpecialIncome)
			tf := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
			// ชื่อ, ชื่อบัญชี, เลขบัญชี, จำนวนเงินที่ต้องโอน, วันที่กรอก
			d := []string{user.GetName(), user.BankAccountName, setValueCSV(user.BankAccountNumber), setValueCSV(utils.FormatCommas(income.NetIncome)), setValueCSV(utils.FormatCommas(income.NetSpecialIncome)), setValueCSV(utils.FormatCommas(summaryIncome)), income.Note, tf}
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

func setValueCSV(s string) string {
	return `="` + s + `"`
}

func (u *usecase) ExportIncomeNotExport(role string) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	defer file.Close()

	if err != nil {
		return "", err
	}
	users, err := u.userRepo.GetByRole(role)
	if err != nil {
		return "", err
	}
	strWrite := make([][]string, 0)
	d := []string{"ชื่อ", "ชื่อบัญชี", "เลขบัญชี", "จำนวนเงินรายได้หลัก", "จำนวนรายได้พิเศษ", "รวมจำนวนที่ต้องโอน", "บันทึกรายการ", "วันที่กรอก"}
	strWrite = append(strWrite, d)
	for _, user := range users {
		income, err := u.repo.GetIncomeByUserID(user.ID.Hex())
		if err == nil {
			u.repo.UpdateExportStatus(income.UserID)
			t := income.SubmitDate
			summaryIncome, _ := calSummary(income.NetIncome, income.NetSpecialIncome)
			tf := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
			// ชื่อ, ชื่อบัญชี, เลขบัญชี, จำนวนเงินที่ต้องโอน, วันที่กรอก
			d := []string{user.GetName(), user.BankAccountName, setValueCSV(user.BankAccountNumber), setValueCSV(utils.FormatCommas(income.NetIncome)), setValueCSV(utils.FormatCommas(income.NetSpecialIncome)), setValueCSV(utils.FormatCommas(summaryIncome)), income.Note, tf}
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
