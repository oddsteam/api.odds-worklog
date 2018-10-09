package user

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	validator "gopkg.in/go-playground/validator.v9"
)

type httpHandler struct {
	usecase Usecase
}

func isRequestValid(m *models.User) (bool, error) {
	if err := validator.New().Struct(m); err != nil {
		return false, err
	}
	return true, nil
}

func (h *httpHandler) createUser(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&u); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.usecase.createUser(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

// func (h *httpHandler) GetUserByID(c echo.Context) (*models.User, error) {

// }

// func (h *httpHandler) Update(c echo.Context) (*models.User, error) {

// }

// func (h *httpHandler) Delete(c echo.Context) (bool, error) {

// }

func NewHttpHandler(e *echo.Echo, session *mongo.Session) {
	ur := newRepository(session)
	uc := newUsecase(ur)

	handler := &httpHandler{uc}
	// e.GET("v1/user", handler.getUser)
	e.POST("v1/user", handler.createUser)
	// e.GET("v1/user/:id", handler.getUserByID)
	// e.DELETE("v1/user/:id", handler.delete)
}
