package income

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/signintech/gopdf"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

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
