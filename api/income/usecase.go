package income

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/signintech/gopdf"
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

func strDelimit(str string, sepstr string, sepcount int) string {
	pos := len(str) - sepcount
	for pos > 0 {
		str = str[:pos] + sepstr + str[pos:]
		pos = pos - sepcount
	}
	return str
}

var gofpdfDir string

func ImageFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "image", fileStr)
}

func (u *usecase) ExportPdf() (string, error) {
	// pdf := gofpdf.New("P", "mm", "A4", "")

	userId := "5bde4e2e1a044b8c9ce44fe4"
	year, month := utils.GetYearMonthNow()

	sd, err_ := u.userRepo.GetUserByID(userId)

	if err_ != nil {
		return "", err_
	}

	str := "1451003242123"
	strPosition := "0105556110718"
	strArray := splitCitizen(str)
	strArrayPosition := splitCitizen(strPosition)

	d := time.Now()
	dm := int(d.Month())
	dd := "27"
	dmn := converseMonthtoThaiName(dm)
	dy := setDy((int(d.Year()) + 543), dm)

	companyName := "บริษัท ออด-อี (ประเทศไทย) จํากัด"
	companyAddress := "2549/41-43 พหลโยธิน ลาดยาว จตุจักร กรุงเทพ 10900"
	employeeName := sd.GetFullname()
	employeeAddress := "265/28 อ.เมือง ต.ในเมือง จ.ชัยภูมิ 36000"
	salaryString := "ห้าร้อยบาทถ้วน"

	fmt.Sprintf("%s", companyName)

	rs, _err := u.repo.GetIncomeUserByYearMonth(userId, year, month)
	// utf8, erro := tis620.ToUTF8("สวัสดีครับ")

	if _err != nil {
		return "", _err
	}

	t1 := rs.UserID
	fmt.Sprintf("%s", t1)

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: 930, H: 1350}}) //595.28, 841.89 = A4
	pdf.AddPage()
	var err error

	err = pdf.AddTTFFont("THSarabunBold", "font/THSarabun-Bold.ttf")
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	err = pdf.SetFont("THSarabunBold", "", 18)
	if err != nil {
		log.Print(err.Error())
		return "", err
	}

	pdf.Image("image/tavi50.png", 0, 0, nil) //print image

	pdf.SetX(97.5)
	pdf.SetY(173.25)
	pdf.Text(companyName)

	pdf.SetX(97.5)
	pdf.SetY(207.25)
	pdf.Text(companyAddress)

	pdf.SetX(97.5)
	pdf.SetY(287.5)
	pdf.Text(employeeName)

	pdf.SetX(97.5)
	pdf.SetY(327.5)
	pdf.Text(employeeAddress)

	pdf.SetX(315)
	pdf.SetY(1060)
	pdf.Text(salaryString)

	pdf.SetX(547.25)
	pdf.SetY(1195)
	pdf.Text(dd)

	pdf.SetX(590)
	pdf.SetY(1195)
	pdf.Text(dmn)

	pdf.SetX(685)
	pdf.SetY(1195)
	pdf.Text(dy)

	positionUserX := 590.75
	positionUserY := 253.75

	for i, r := range strArray {
		if i == 0 {
			pdf.SetX(positionUserX)
			pdf.SetY(positionUserY)
			pdf.Text(r)
		} else if i == 1 || i == 5 || i == 10 || i == 12 {
			positionUserX += 29
			pdf.SetX(positionUserX)
			pdf.SetY(positionUserY)
			pdf.Text(r)
		} else {
			positionUserX += 19
			pdf.SetX(positionUserX)
			pdf.SetY(positionUserY)
			pdf.Text(r)
		}
	}

	positionX := 590.75
	positionY := 147.25

	for j, r := range strArrayPosition {
		if j == 0 {
			pdf.SetX(positionX)
			pdf.SetY(positionY)
			pdf.Text(r)
		} else if j == 1 || j == 5 || j == 10 || j == 12 {
			positionX += 29
			pdf.SetX(positionX)
			pdf.SetY(positionY)
			pdf.Text(r)
		} else {
			positionX += 19
			pdf.SetX(positionX)
			pdf.SetY(positionY)
			pdf.Text(r)
		}
	}

	t := time.Now()
	tf := fmt.Sprintf("%d_%02d_%02d_%02d_%02d_%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second())
	filename := fmt.Sprintf("files/tavi50/%s_%s.pdf", "tavi50", tf)

	error := pdf.WritePdf(filename)

	if error != nil {
		return "", error
	}

	return filename, nil
}

func converseMonthtoThaiName(dm int) string {
	dmt := [12]string{"มกราคม", "กุมภาพันธ์", "มีนาคม", "เมษายน", "พฤษภาคม", "มิถุนายน", "กรกฎาคม", "สิงหาคม", "กันยายน", "ตุลาคม", "พฤศจิกายน", "ธันวาคม"}
	monthThaiName := ""
	for i, v := range dmt {
		if dm == 1 {
			monthThaiName = dmt[len(dmt)-1]
		}
		if i+1 == dm-1 {
			monthThaiName = v
		}
	}
	return monthThaiName
}

func setDy(dy int, dm int) string {
	year := strconv.Itoa(dy)

	if dm == 1 {
		year = strconv.Itoa(dy - 1)
	}

	return year
}

func splitCitizen(citizen string) []string {
	citizenArrey := []string{}
	for _, r := range citizen {
		c := string(r)
		citizenArrey = append(citizenArrey, c)
	}
	return citizenArrey
}

func getImageBytes() []byte {
	b, err := ioutil.ReadFile("image/tavi50.png")
	if err != nil {
		panic(err)
	}
	return b
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

func (u *usecase) DropIncome() error {
	return u.repo.DropIncome()
}

func setValueCSV(s string) string {
	return `="` + s + `"`
}
