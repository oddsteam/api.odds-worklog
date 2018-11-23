package setting

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

//NewHTTPHandler for setting
func NewHTTPHandler(r *echo.Group, session *mongo.Session) {
	reminderRepo := NewRepository(session)

	r = r.Group("/setting")
	r.POST("/reminder", func(c echo.Context) error {
		return save(c, reminderRepo)
	})
	r.GET("/reminder", func(c echo.Context) error {
		return get(c, reminderRepo)
	})
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
func save(c echo.Context, reminderRepo Repository) error {
	reminder := new(models.Reminder)
	if err := c.Bind(&reminder); err != nil {
		return utils.NewError(c, 400, utils.ErrBadRequest)
	}
	r, err := reminderRepo.SaveReminder(reminder)
	if err != nil {
		return utils.NewError(c, 400, utils.ErrBadRequest)
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
func get(c echo.Context, reminderRepo Repository) error {
	r, err := reminderRepo.GetReminder()
	if err != nil {
		return utils.NewError(c, 500, utils.ErrNotFound)
	}
	return c.JSON(http.StatusOK, r)
}
