package user

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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

func (h *httpHandler) getUser(c echo.Context) error {
	users, err := h.usecase.getUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *httpHandler) getUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.usecase.getUserByID(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *httpHandler) updateUser(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&u); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.usecase.updateUser(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *httpHandler) deleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.usecase.deleteUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *httpHandler) login(c echo.Context) error {
	var u models.Login
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	user, err := h.usecase.login(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

func NewHttpHandler(e *echo.Echo, session *mongo.Session) {
	ur := newRepository(session)
	uc := newUsecase(ur)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	handler := &httpHandler{uc}
	config := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	e.POST("/login", handler.login)

	r := e.Group("/")
	r.Use(middleware.JWTWithConfig(config))
	r.GET("user", handler.getUser)
	r.POST("user", handler.createUser)
	r.GET("user/:id", handler.getUserByID)
	r.PUT("user", handler.updateUser)
	r.DELETE("user/:id", handler.deleteUser)

}
