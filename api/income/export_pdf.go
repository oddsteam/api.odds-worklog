package income

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/signintech/gopdf"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

var gofpdfDir string

func ImageFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "image", fileStr)
}

func (u *usecase) ExportPdf(id string) (string, error) {
	// pdf := gofpdf.New("P", "mm", "A4", "")

	userId := id
	year, month := utils.GetYearMonthNow()
	var fileNames []string
	var months []int
	var days []string
	var incomes []string
	var whts []string
	sd, err_ := u.userRepo.GetByID(userId)

	if err_ != nil {
		return "", err_
	}

	str := sd.GetThaiCitizenID()
	strPosition := "0105556110718"
	strArray := splitCitizen(str)
	strArrayPosition := splitCitizen(strPosition)

	companyName := "บริษัท ออด-อี (ประเทศไทย) จํากัด"
	companyAddress := "2549/41-43 พหลโยธิน ลาดยาว จตุจักร กรุงเทพ 10900"
	employeeName := sd.GetName()
	employeeAddress := sd.GetAddress()

	fmt.Sprintf("%s", companyName)

	for i := 1; i <= int(month); i++ {
		rs, _ := u.repo.GetIncomeUserByYearMonth(userId, year, time.Month(i))
		income, _ := u.repo.GetIncomeByUserID(userId)
		if rs != nil {
			months = append(months, int(rs.SubmitDate.Month()))
			days = append(days, fmt.Sprintf("%02d", rs.SubmitDate.Day()))
		}

		if income != nil {
			incomes = append(incomes, income.TotalIncome)
			whts = append(whts, income.WHT)
		}
	}
	for i := 0; i <= len(months)-1; i++ {
		d := time.Now()
		dd := days[i]
		dmn := converseMonthtoThaiName(months[i])
		dy := setDy((int(d.Year()) + 543), months[i])
		ti := incomes[i]
		wht := whts[i]
		// utf8, erro := tis620.ToUTF8("สวัสดีครับ")
		salaryString := ConvertIntToThaiBath(wht)

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

		pdf.SetX(830)
		pdf.SetY(93)
		pdf.Text(dy)

		pdf.SetX(530)
		pdf.SetY(500)
		pdf.Text(dd + "/" + fmt.Sprintf("%02d", months[i]) + "/" + dy)

		pdf.SetX(650)
		pdf.SetY(500)
		pdf.Text(ti)

		pdf.SetX(780)
		pdf.SetY(500)
		pdf.Text(wht)

		pdf.SetX(650)
		pdf.SetY(1030)
		pdf.Text(ti)

		pdf.SetX(780)
		pdf.SetY(1030)
		pdf.Text(wht)

		pdf.Image("image/check.jpg", 130, 1115, nil)
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
		filename := fmt.Sprintf("tavi50/%s_%d_%s.pdf", "tavi50", months[i], tf)

		error := pdf.WritePdf(filename)
		fileNames = append(fileNames, filename)

		if error != nil {
			return "", error
		}
	}
	if err := ZipFiles("tavi50.zip", fileNames); err != nil {
		panic(err)
	}
	return "tavi50.zip", nil
}

func converseMonthtoThaiName(dm int) string {
	dmt := [12]string{"มกราคม", "กุมภาพันธ์", "มีนาคม", "เมษายน", "พฤษภาคม", "มิถุนายน", "กรกฎาคม", "สิงหาคม", "กันยายน", "ตุลาคม", "พฤศจิกายน", "ธันวาคม"}
	monthThaiName := ""
	for i, v := range dmt {
		if dm == 1 {
			monthThaiName = dmt[len(dmt)-12]
		}
		if i+1 == dm {
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
func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

func ZipFiles(filename string, files []string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {
		if err = AddFileToZip(zipWriter, file); err != nil {
			return err
		}
	}
	return nil
}

func ConvertIntToThaiBath(txt string) string {
	var bahtTH = ""
	amonth, err := utils.StringToFloat64(txt)
	if err != nil {
		amonth = 0
	}
	bahtTxt := utils.FloatToString(amonth)
	num := [11]string{"ศูนย์", "หนึ่ง", "สอง", "สาม", "สี่", "ห้า", "หก", "เจ็ด", "แปด", "เก้า", "สิบ"}
	rank := [7]string{"", "สิบ", "ร้อย", "พัน", "หมื่น", "แสน", "ล้าน"}
	var temp []string
	temp = strings.Split(bahtTxt, ".")
	intVal := temp[0]
	decVal := temp[1]
	bathThai, _ := utils.StringToFloat64(bahtTxt)
	if bathThai == 0 {
		bahtTH = "ศูนย์บาทถ้วน"
	} else {
		for i := 0; i < len(intVal); i++ {
			n := intVal[i : i+1]
			if n != "0" {
				if i == (len(intVal)-1) && (n == "1") {
					bahtTH += "เอ็ด"
				} else if i == (len(intVal)-2) && (n == "2") {
					bahtTH += "ยี่"
				} else if i == (len(intVal)-2) && (n == "1") {
					bahtTH += ""
				} else {
					position, _ := strconv.Atoi(n)
					bahtTH += num[position]
					bahtTH += rank[(len(intVal)-i)-1]
				}
			}
		}
		bahtTH += "บาท"
		if decVal == "00" {
			bahtTH += "ถ้วน"
		} else {
			for i := 0; i < len(decVal); i++ {
				n := decVal[i : i+1]
				if n != "0" {
					if i == (len(decVal)-1) && (n == "1") {
						bahtTH += "เอ็ด"
					} else if i == (len(decVal)-2) && (n == "2") {
						bahtTH += "ยี่"
					} else if i == (len(decVal)-2) && (n == "ๅ") {
						bahtTH += ""
					} else {
						position, _ := strconv.Atoi(n)
						bahtTH += num[position]
						bahtTH += rank[(len(decVal)-i)-1]
					}
				}
			}
			bahtTH += "สตางค์"
			bahtTH += "ถ้วน"

		}

	}
	return bahtTH

}
