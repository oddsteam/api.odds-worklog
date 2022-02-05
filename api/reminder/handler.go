package reminder

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/smtp"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/file"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"gitlab.odds.team/worklog/api.odds-worklog/worker"
)

// NewHttpHandler for reminder resource godoc
func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	reminderRepo := NewRepository(session)
	userRepo := user.NewRepository(session)
	usecaseFile := file.NewUsecase(userRepo)
	r = r.Group("/reminder")
	r.GET("/setting", func(c echo.Context) error {
		return GetReminder(c, reminderRepo)
	})

	r.POST("/setting", func(c echo.Context) error {
		return SaveReminder(c, reminderRepo)
	})
	r.POST("/mail/:id", func(c echo.Context) error {
		return SendMail(c, userRepo, usecaseFile)
	})
}

func validateReminderRequest(reminder *models.Reminder) error {
	if reminder.Name != "reminder" {
		return errors.New("Request Name is not reminder")
	}
	if reminder.Setting.Date != "25" && reminder.Setting.Date != "26" && reminder.Setting.Date != "27" {
		return errors.New("Request Setting Date is not between 25-27")
	}
	if reminder.Setting.Message == "" {
		return errors.New("Request Setting Message is empty")
	}
	return nil
}

// SaveReminder Setting godoc
// @Summary Save Reminder Setting
// @Description Save Reminder Setting
// @Tags reminder
// @Accept  json
// @Produce  json
// @Param reminder body models.Reminder true  "line, slack, facebook can empty"
// @Success 200 {object} models.Reminder
// @Failure 400 {object} utils.HTTPError
// @Router /setting/reminder [post]
func SaveReminder(c echo.Context, reminderRepo Repository) error {
	reminder := new(models.Reminder)
	if err := c.Bind(&reminder); err != nil {
		return utils.NewError(c, 400, utils.ErrBadRequest)
	}
	if err := validateReminderRequest(reminder); err != nil {
		return utils.NewError(c, 400, err)
	}
	r, err := reminderRepo.SaveReminder(reminder)
	if err != nil {
		return utils.NewError(c, 500, errors.New("Can not insert data into DB"))
	}
	worker.StartWorker(reminder)
	return c.JSON(http.StatusOK, r)
}

// GetReminder Setting godoc
// @Summary Get Reminder Setting
// @Description Get Reminder Setting
// @Tags reminder
// @Produce  json
// @Success 200 {object} models.Reminder
// @Failure 500 {object} utils.HTTPError
// @Router /setting/reminder [get]
func GetReminder(c echo.Context, reminderRepo Repository) error {
	r, err := reminderRepo.GetReminder()
	if err != nil {
		return utils.NewError(c, 500, errors.New("Data not found in DB"))
	}
	return c.JSON(http.StatusOK, r)
}

func ListEmailUserIncomeStatusIsNo(incomeUsecase income.Usecase) ([]string, error) {
	emails := []string{}
	incomeIndividualStatusList, err := incomeUsecase.GetIncomeStatusList("individual", false)
	if err != nil {
		return nil, err
	}
	incomeCorpStatusList, err := incomeUsecase.GetIncomeStatusList("corporate", false)
	if err != nil {
		return nil, err
	}
	incomeStatusList := append(incomeIndividualStatusList, incomeCorpStatusList...)
	for _, incomeStatus := range incomeStatusList {
		// if incomeStatus.Status == "N" {
		emails = append(emails, incomeStatus.User.Email)
		// }
	}
	return emails, nil
}

func SendMail(c echo.Context, userRepo user.Repository, usecaseFile file.Usecase) error {
	id := c.Param("id")
	user, err := userRepo.GetByID(id)
	if err != nil {
		return err
	}
	admins, err := userRepo.GetByRole("admin")
	if err != nil {
		return err
	}
	receive := []string{}
	for i := 0; i < len(admins); i++ {
		receive = append(receive, admins[i].Email)
	}
	fileName, err := usecaseFile.GetPathIDCard(id)
	if err != nil {
		return err
	}
	sender := New()
	m := NewMessage("[ODDS] แจ้ง User ใหม่เข้าใช้งานระบบ Worklog", "File PDF ID Card ของคุณ "+user.FirstName+" "+user.LastName+" เป็น User ใหม่ที่เข้าใช้งานในระบบ worklog.odds.team \n สามารถติดต่อได้ที่ Email : "+user.Email)
	m.To = receive
	m.AttachFile(fileName)
	fmt.Println(sender.Send(m))
	return c.JSON(http.StatusOK, "Send Mail Success ")
}

func New() *Sender {
	auth := smtp.PlainAuth("", "oddsnotify@gmail.com", "@abcd12345", "smtp.gmail.com")
	return &Sender{auth}
}

func (s *Sender) Send(m *Message) error {
	return smtp.SendMail(fmt.Sprintf("%s:%s", "smtp.gmail.com", "587"), s.auth, "oddsnotify@gmail.com", m.To, m.ToBytes())
}

func NewMessage(s, b string) *Message {
	return &Message{Subject: s, Body: b, Attachments: make(map[string][]byte)}
}

func (m *Message) AttachFile(src string) error {
	b, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	_, fileName := filepath.Split(src)
	m.Attachments[fileName] = b
	return nil
}

func (m *Message) ToBytes() []byte {
	buf := bytes.NewBuffer(nil)
	withAttachments := len(m.Attachments) > 0
	buf.WriteString(fmt.Sprintf("Subject: %s\n", m.Subject))
	buf.WriteString(fmt.Sprintf("To: %s\n", strings.Join(m.To, ",")))
	if len(m.CC) > 0 {
		buf.WriteString(fmt.Sprintf("Cc: %s\n", strings.Join(m.CC, ",")))
	}

	if len(m.BCC) > 0 {
		buf.WriteString(fmt.Sprintf("Bcc: %s\n", strings.Join(m.BCC, ",")))
	}

	buf.WriteString("MIME-Version: 1.0\n")
	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()
	if withAttachments {
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\n", boundary))
		buf.WriteString(fmt.Sprintf("--%s\n", boundary))
	} else {
		buf.WriteString("Content-Type: text/plain; charset=utf-8\n")
	}

	buf.WriteString(m.Body)
	if withAttachments {
		for k, v := range m.Attachments {
			buf.WriteString(fmt.Sprintf("\n\n--%s\n", boundary))
			buf.WriteString(fmt.Sprintf("Content-Type: %s\n", http.DetectContentType(v)))
			buf.WriteString("Content-Transfer-Encoding: base64\n")
			buf.WriteString(fmt.Sprintf("Content-Disposition: attachment; filename=%s\n", k))

			b := make([]byte, base64.StdEncoding.EncodedLen(len(v)))
			base64.StdEncoding.Encode(b, v)
			buf.Write(b)
			buf.WriteString(fmt.Sprintf("\n--%s", boundary))
		}

		buf.WriteString("--")
	}

	return buf.Bytes()
}

type Sender struct {
	auth smtp.Auth
}

type Message struct {
	To          []string
	CC          []string
	BCC         []string
	Subject     string
	Body        string
	Attachments map[string][]byte
}
