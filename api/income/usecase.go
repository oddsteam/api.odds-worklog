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
	Net string
	VAT string
	WHT string
}

func NewUsecase(r Repository, ur user.Repository) Usecase {
	return &usecase{r, ur}
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

func calIncomeSum(income string, vattype string) (*incomeSum, error) {
	var vat, wht string
	var vatf, whtf float64
	var ins = new(incomeSum)

	total, err := utils.StringToFloat64(income)
	if err != nil {
		return nil, err
	}
	wht, whtf, err = calWHT(income)
	if err != nil {
		return nil, err
	}

	ins.WHT = wht

	if vattype == "Y" {
		vat, vatf, err = calVAT(income)
		if err != nil {
			return nil, err
		}

		net := total + vatf - whtf

		ins.Net = utils.FloatToString(net)
		ins.VAT = vat
		return ins, nil
	}

	net := total - whtf
	ins.Net = utils.FloatToString(net)
	return ins, nil
}

func (u *usecase) GetIncomeStatusList(role string) ([]*models.IncomeStatus, error) {
	var incomeList []*models.IncomeStatus
	users, err := u.userRepo.GetUserByRole(role)
	if err != nil {
		return nil, err
	}

	year, month := utils.GetYearMonthNow()
	for index, element := range users {
		element.ThaiCitizenID = ""
		incomeUser, err := u.repo.GetIncomeUserByYearMonth(element.ID.Hex(), year, month)
		income := models.IncomeStatus{User: element}
		incomeList = append(incomeList, &income)
		if err != nil {
			incomeList[index].Status = "N"
		} else {
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

func (u *usecase) ExportIncome(role string) (string, error) {
	file, filename, err := utils.CreateCVSFile(role)
	defer file.Close()

	if err != nil {
		return "", err
	}

	users, err := u.userRepo.GetUserByRole(role)
	if err != nil {
		return "", err
	}

	year, month := utils.GetYearMonthNow()

	strWrite := make([][]string, 0)
	d := []string{"ชื่อ", "ชื่อบัญชี", "เลขบัญชี", "จำนวนเงินที่ต้องโอน", "วันที่กรอก"}
	strWrite = append(strWrite, d)
	for _, user := range users {
		income, err := u.repo.GetIncomeUserByYearMonth(user.ID.Hex(), year, month)
		if err == nil {
			t := income.SubmitDate
			tf := fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d", t.Day(), int(t.Month()), t.Year(), (t.Hour() + 7), t.Minute(), t.Second())
			// ชื่อ, ชื่อบัญชี, เลขบัญชี, จำนวนเงินที่ต้องโอน, วันที่กรอก
			d := []string{user.GetFullname(), user.BankAccountName, setValueCSV(user.BankAccountNumber), setValueCSV(utils.FormatCommas(income.NetIncome)), tf}
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
