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
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/file"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

// NewHttpHandler for reminder resource godoc
func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	userRepo := user.NewRepository(session)
	usecaseFile := file.NewUsecase(userRepo)
	r = r.Group("/reminder")
	r.POST("/mail/:id", func(c echo.Context) error {
		return SendMail(c, userRepo, usecaseFile)
	})
}

func SendMail(c echo.Context, userRepo user.Repository, usecaseFile file.Usecase) error {
	id := c.Param("id")
	user, err := userRepo.GetByID(id)
	if err != nil {
		return err
	}
	fileName, err := usecaseFile.GetPathIDCard(id)
	if err != nil {
		return err
	}
	m := CreateMailMessage(*user, fileName)
	if len(m.To) == 0 {
		return utils.NewError(c, http.StatusInternalServerError, errors.New("SMTP_REMINDER_RECEIVERS is required for reminder emails"))
	}
	sender := New()
	fmt.Println(sender.Send(m))
	return c.JSON(http.StatusOK, "Send Mail Success")
}

func CreateMailMessage(user models.User, fileName string) *Message {
	m := NewMessage("[ODDS] แจ้ง User ใหม่เข้าใช้งานระบบ Worklog", "File PDF ID Card ของคุณ "+user.BankAccountName+" ("+user.FirstName+" "+user.LastName+") เป็น User ใหม่ที่เข้าใช้งานในระบบ worklog.odds.team \n สามารถติดต่อได้ที่ Email : "+user.Email)
	receive := getReminderReceivers()
	m.To = receive
	m.AttachFile(fileName)
	return m
}

func getReminderReceivers() []string {
	receivers := os.Getenv("SMTP_REMINDER_RECEIVERS")
	if receivers == "" {
		return []string{}
	}
	return strings.Split(strings.ReplaceAll(receivers, " ", ""), ",")
}

func New() *Sender {
	auth := smtp.PlainAuth("", "oddsnotify@gmail.com", "wvquzrqfmiuckurw", "smtp.gmail.com")
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
