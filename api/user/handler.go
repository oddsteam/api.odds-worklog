package user

import (
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"

	"github.com/globalsign/mgo/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

type HttpHandler struct {
	Usecase Usecase
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	sr := site.NewRepository(session)
	ur := NewRepository(session)
	uc := NewUsecase(ur, sr)
	handler := &HttpHandler{uc}

	r = r.Group("/users")
	r.GET("", handler.Get)
	r.POST("", handler.Create)
	r.GET("/:id", handler.GetByID)
	r.GET("/email/:email", handler.GetByEmail)
	r.GET("/site/:id", handler.GetBySiteID)
	r.PUT("/:id", handler.Update)
	r.DELETE("/:id", handler.Delete)
}

func getUserFromToken(c echo.Context) *models.User {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}

// Create godoc
// @Summary Create User
// @Description Create User
// @Tags users
// @Accept json
// @Produce json
// @Param user body models.User true  "id can empty"
// @Success 200 {object} models.User
// @Failure 400 {object} utils.HTTPError
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users [post]
func (h *HttpHandler) Create(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	var u models.User
	if err := c.Bind(&u); err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	user, err := h.Usecase.Create(&u)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, user)
}

// Get godoc
// @Summary List user
// @Description get user list
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users [get]
func (h *HttpHandler) Get(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}
	users, err := h.Usecase.Get()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}

// GetById godoc
// @Summary Get User By Id
// @Description Get User By Id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.User
// @Failure 204 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Router /users/{id} [get]
func (h *HttpHandler) GetByID(c echo.Context) error {
	id := c.Param("id")
	user, err := h.Usecase.GetByID(id)
	if err != nil {
		return utils.NewError(c, http.StatusNoContent, err)
	}
	return c.JSON(http.StatusOK, user)
}

// GetByEmail godoc
// @Summary Get User By Email
// @Description Get User By Email
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Email"
// @Success 200 {array} models.User
// @Failure 204 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Router /users/{email} [get]
func (h *HttpHandler) GetByEmail(c echo.Context) error {
	email := c.Param("email")
	user, err := h.Usecase.GetByEmail(email)
	if err != nil {
		return utils.NewError(c, http.StatusNoContent, err)
	}
	return c.JSON(http.StatusOK, user)
}

// GetBySiteId godoc
// @Summary Get User By Site Id
// @Description Get User By Site Id
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "Site id"
// @Success 200 {array} models.User
// @Failure 204 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Router /users/site/{id} [get]
func (h *HttpHandler) GetBySiteID(c echo.Context) error {
	u := getUserFromToken(c)
	if !u.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	id := c.Param("id")
	user, err := h.Usecase.GetBySiteID(id)
	if err != nil {
		return utils.NewError(c, http.StatusNoContent, err)
	}
	return c.JSON(http.StatusOK, user)
}

// Update godoc
// @Summary Update User
// @Description Update User
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body models.User true  "id can empty"
// @Success 200 {object} models.User
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users/{id} [put]
func (h *HttpHandler) Update(c echo.Context) error {
	var u models.User
	if err := c.Bind(&u); err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}
	id := c.Param("id")
	u.ID = bson.ObjectIdHex(id)
	ut := getUserFromToken(c)
	user, err := h.Usecase.Update(&u, ut.IsAdmin())
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary Delete User
// @Description Delete User
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /users/{id} [delete]
func (h *HttpHandler) Delete(c echo.Context) error {
	u := getUserFromToken(c)
	if !u.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	id := c.Param("id")
	err := h.Usecase.Delete(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, models.Response{Message: "Delete user success."})
}
