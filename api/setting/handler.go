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
}

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
