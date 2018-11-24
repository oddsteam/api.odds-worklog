package setting

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

//NewHTTPHandler for setting
func NewHTTPHandler(r *echo.Group, session *mongo.Session, m middleware.JWTConfig) {
	reminderRepo := NewRepository(session)

	r = r.Group("/setting")
	r.GET("/reminder", func(c echo.Context) error {
		return Get(c, reminderRepo)
	})
	r.Use(middleware.JWTWithConfig(m))
	r.POST("/reminder", func(c echo.Context) error {
		return Save(c, reminderRepo)
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

// Save Reminder Setting godoc
// @Summary Save Reminder Setting
// @Description Save Reminder Setting
// @Tags reminder
// @Accept  json
// @Produce  json
// @Param reminder body models.Reminder true  "line, slack, facebook can empty"
// @Success 200 {object} models.Reminder
// @Failure 400 {object} utils.HTTPError
// @Router /setting/reminder [post]
func Save(c echo.Context, reminderRepo Repository) error {
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
	return c.JSON(http.StatusOK, r)
}

// Get Reminder Setting godoc
// @Summary Get Reminder Setting
// @Description Get Reminder Setting
// @Tags reminder
// @Produce  json
// @Success 200 {object} models.Reminder
// @Failure 500 {object} utils.HTTPError
// @Router /setting/reminder [get]
func Get(c echo.Context, reminderRepo Repository) error {
	r, err := reminderRepo.GetReminder()
	if err != nil {
		return utils.NewError(c, 500, errors.New("Data not found in DB"))
	}
	return c.JSON(http.StatusOK, r)
}
