package user

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"

	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
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

// CreateUser godoc
// @Summary Create User
// @Description Create User
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body models.User true  "id can empty"
// @Success 200 {array} models.User
// @Failure 400 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users [post]
func (h *HttpHandler) CreateUser(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}

	if ok, err := isRequestValid(&u); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	user, err := h.Usecase.CreateUser(&u)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary List user
// @Description get user list
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {array} models.User
// @Failure 500 {object} utils.HTTPError
// @Router /users [get]
func (h *HttpHandler) GetUser(c echo.Context) error {
	checkUser, message := IsUserAdmin(c)
	if !checkUser {
		return c.JSON(http.StatusUnauthorized, message)
	}
	users, err := h.Usecase.GetUser()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

// GetUserById godoc
// @Summary Get User By Id
// @Description Get User By Id
// @Tags users
// @Accept  multipart/form-data
// @Produce  json
// @Param  id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 204 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Router /users/{id} [get]
func (h *HttpHandler) GetUserByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	user, err := h.Usecase.GetUserByID(id)
	if err != nil {
		return utils.NewError(c, http.StatusNoContent, err)
	}
	return c.JSON(http.StatusOK, user)
}

// GetUserBySiteId godoc
// @Summary Get User By Site Id
// @Description Get User By Site Id
// @Tags users
// @Accept  multipart/form-data
// @Produce  json
// @Param  id path string true "Site id"
// @Success 200 {object} models.User
// @Failure 204 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Router /users/site/{id} [get]
func (h *HttpHandler) GetUserBySiteID(c echo.Context) error {
	isAdmin, message := IsUserAdmin(c)
	if !isAdmin {
		return c.JSON(http.StatusUnauthorized, message)
	}

	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	user, err := h.Usecase.GetUserBySiteID(id)
	if err != nil {
		return utils.NewError(c, http.StatusNoContent, err)
	}
	return c.JSON(http.StatusOK, user)
}

// UpdateUserById godoc
// @Summary Update User By Id
// @Description Update User By Id
// @Tags users
// @Accept  multipart/form-data
// @Produce  json
// @Param  id path string true "User ID"
// @Param user body models.User true  "id can empty"
// @Success 200 {object} models.User
// @Failure 400 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users/{id} [put]
func (h *HttpHandler) UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	u := models.User{
		ID: bson.ObjectIdHex(id),
	}
	if err := c.Bind(&u); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}

	if ok, err := isRequestValid(&u); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	ut := getUserFromToken(c)
	u.Email = ut.Email

	file, _ := c.FormFile("file")
	user, err := h.Usecase.UpdateUser(&u, file)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete User
// @Description Delete User By Id
// @Tags users
// @Accept  multipart/form-data
// @Produce  json
// @Param  id path string true "User ID"
// @Success 204 {object} models.User
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users/{id} [delete]
func (h *HttpHandler) DeleteUser(c echo.Context) error {
	checkUser, message := IsUserAdmin(c)
	if !checkUser {
		return c.JSON(http.StatusUnauthorized, message)
	}
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	err := h.Usecase.DeleteUser(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// UpdatePartialUser godoc
// @Summary Update Partial User
// @Description Delete Update Partial User
// @Tags users
// @Accept  multipart/form-data
// @Produce  json
// @Param  id path string true "User ID"
// @Param user body models.User true  "id can empty"
// @Success 200 {object} models.User
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users/{id} [patch]
func (h *HttpHandler) UpdatePartialUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	user, err := h.Usecase.GetUserByID(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	b, _ := ioutil.ReadAll(c.Request().Body)
	err = json.Unmarshal(b, &user)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	file, _ := c.FormFile("file")
	newUser, err := h.Usecase.UpdateUser(user, file)
	return c.JSON(http.StatusOK, newUser)
}

func IsUserAdmin(c echo.Context) (bool, string) {
	u := getUserFromToken(c)
	if u.IsAdmin() {
		return true, ""
	}
	return false, "ไม่มีสิทธิในการใช้งาน"
}

func getUserFromToken(c echo.Context) *models.User {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	ur := NewRepository(session)
	uc := NewUsecase(ur)
	handler := &HttpHandler{uc}
	r = r.Group("/users")
	r.GET("", handler.GetUser)
	r.POST("", handler.CreateUser)
	r.GET("/:id", handler.GetUserByID)
	r.GET("/site/:id", handler.GetUserBySiteID)
	r.PUT("/:id", handler.UpdateUser)
	r.DELETE("/:id", handler.DeleteUser)
	r.PATCH("/:id", handler.UpdatePartialUser)
}
