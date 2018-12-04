package reminder

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/income"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"gitlab.odds.team/worklog/api.odds-worklog/worker"
)

// NewHTTPHandler for reminder resource godoc
func NewHTTPHandler(r *echo.Group, session *mongo.Session) {
	reminderRepo := NewRepository(session)

	r = r.Group("/reminder")
	r.GET("/setting", func(c echo.Context) error {
		return GetReminder(c, reminderRepo)
	})

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
		// if incomeStatus.Status == "N" {
		emails = append(emails, incomeStatus.User.Email)
		// }
	}
	return emails, nil
}
