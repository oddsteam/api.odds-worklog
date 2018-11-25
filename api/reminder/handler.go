package reminder

import (
	"errors"
	"net/http"
	"os/exec"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/slack"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

// NewHTTPHandler for reminder resource godoc
func NewHTTPHandler(r *echo.Group, session *mongo.Session, m middleware.JWTConfig) {
	userRepo := user.NewRepository(session)
	incomeRepo := income.NewRepository(session)
	reminderRepo := NewRepository(session)
	incomeUsecase := income.NewUsecase(incomeRepo, userRepo)

	r = r.Group("/reminder")
	r.GET("/setting", func(c echo.Context) error {
		return GetReminder(c, reminderRepo)
	})
	r.GET("/send", func(c echo.Context) error {
		return send(c, incomeUsecase, reminderRepo)
	})

	r.Use(middleware.JWTWithConfig(m))
	r.POST("/setting", func(c echo.Context) error {
		return SaveReminder(c, reminderRepo)
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
	go setCronJobReminder()
	return c.JSON(http.StatusOK, r)
}

func setCronJobReminder() {
	exec.Command("/bin/sh", "/app/updateCrontab.sh")
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

// Send Reminder godoc
// @Summary Send Reminder
// @Description Send Notification Reminder
// @Tags reminder
// @Produce  json
// @Success 200 {object} string
// @Failure 500 {object} utils.HTTPError
// @Router /reminder/send [get]
func send(c echo.Context, incomeUsecase income.Usecase, reminder Repository) error {
	var token string
	isDev := c.QueryParam("isDev")
	var emails []string
	if isDev == "false" {
		token = "xoxb-293071900534-486896062132-2RMbUSdX6DqoOKsVMCSXQoiM" // Odds workspace
		user, err := ListEmailUserIncomeStatusIsNo(incomeUsecase)
		if err != nil {
			return utils.NewError(c, 500, err)
		}
		emails = user
	} else {
		token = "xoxb-484294901968-485201164352-IC904vZ6Bxwx2xkI2qzWgy5J" // Reminder workspace
		emails = []string{
			"tong@odds.team",
			"saharat@odds.team",
			"thanundorn@odds.team",
			"santi@odds.team",
		}
	}

	s, err := reminder.GetReminder()
	if err != nil {
		return utils.NewError(c, 500, err)
	}

	err = sendNotification(token, emails, s.Setting.Message)
	if err != nil {
		return utils.NewError(c, 500, err)
	}

	return c.JSON(http.StatusOK, true)
}

func ListEmailUserIncomeStatusIsNo(incomeUsecase income.Usecase) ([]string, error) {
	emails := []string{}
	incomeIndividualStatusList, err := incomeUsecase.GetIncomeStatusList("N")
	if err != nil {
		return nil, err
	}
	incomeCorpStatusList, err := incomeUsecase.GetIncomeStatusList("Y")
	if err != nil {
		return nil, err
	}
	incomeStatusList := append(incomeIndividualStatusList, incomeCorpStatusList...)
	for _, incomeStatus := range incomeStatusList {
		if incomeStatus.Status == "N" {
			emails = append(emails, incomeStatus.User.Email)
		}
	}
	return emails, nil
}

func sendNotification(token string, emails []string, message string) error {
	client := slack.Client{
		Token: token,
	}
	slackUsers, err := client.GetUserList()
	if err != nil {
		return err
	}
	for _, email := range emails {
		for _, member := range slackUsers.Members {
			if member.Profile.Email == email {
				im, err := client.OpenIMChannel(member.ID)
				if err != nil {
					return err
				}
				channelID := im.Channel.ID
				client.PostMessage(channelID, message)
			}
		}
	}
	return nil
}
