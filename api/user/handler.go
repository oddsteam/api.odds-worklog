package user

import (
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	validator "gopkg.in/go-playground/validator.v9"
)

type HttpHandler struct {
	Usecase Usecase
}

func isRequestValid(m *models.User) (bool, error) {
	if err := validator.New().Struct(m); err != nil {
		return false, err
	}
	return true, nil
}

func (h *HttpHandler) CreateUser(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&u); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.Usecase.CreateUser(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

func (h *HttpHandler) GetUser(c echo.Context) error {
	users, err := h.Usecase.GetUser()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

func (h *HttpHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.Usecase.GetUserByID(id)
	if err != nil {
		return c.JSON(http.StatusNoContent, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HttpHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	u := models.User{
		ID: bson.ObjectIdHex(id),
	}
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&u); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	user, err := h.Usecase.UpdateUser(&u)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, user)
}

func (h *HttpHandler) DeleteUser(c echo.Context) error {
	id := c.Param("id")
	err := h.Usecase.DeleteUser(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, models.ResponseError{Message: err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *HttpHandler) Login(c echo.Context) error {
	var u models.Login
	if err := c.Bind(&u); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}
	username := u.Username
	password := u.Password
	if username == "admin" && password == "admin" {
		// Set custom claims
		claims := &models.JwtCustomClaims{
			"Admin!",
			true,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}

		// Create token with claims
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return c.JSON(http.StatusUnauthorized, err.Error())
		}
		TK := &models.Token{
			Token: t,
		}
		return c.JSON(http.StatusOK, TK)
	}

	return c.JSON(http.StatusUnauthorized, nil)
}
func NewHttpHandler(e *echo.Echo, session *mongo.Session) {
	ur := newRepository(session)
	uc := newUsecase(ur)
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	handler := &HttpHandler{uc}
	config := middleware.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		SigningKey: []byte("secret"),
	}
	e.POST("/login", handler.Login)

	r := e.Group("/")
	r.Use(middleware.JWTWithConfig(config))
	r.GET("user", handler.GetUser)
	r.POST("user", handler.CreateUser)
	r.GET("user/:id", handler.GetUserByID)
	r.PUT("user/:id", handler.UpdateUser)
	r.DELETE("user/:id", handler.DeleteUser)

}
