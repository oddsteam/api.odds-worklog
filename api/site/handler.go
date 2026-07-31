package site

import (
	"errors"
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/bsonutil"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	validator "gopkg.in/go-playground/validator.v9"
)

type HttpHandler struct {
	Usecase Usecase
}

func getUserFromToken(c echo.Context) *models.UserClaims {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}

func isRequestValid(m *models.Site) (bool, error) {
	if err := validator.New().Struct(m); err != nil {
		return false, err
	}
	return true, nil
}

// CreateSiteGroup godoc
// @Summary Create Site Group
// @Description Create Site Group
// @Tags sites
// @Accept  json
// @Produce  json
// @Param site body models.Site true  "id can empty"
// @Success 200 {array} models.Site
// @Failure 400 {object} utils.HTTPError
// @Failure 403 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /sites [post]
func (h *HttpHandler) CreateSiteGroup(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsUserManager() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	var site models.Site
	if err := c.Bind(&site); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}

	if ok, err := isRequestValid(&site); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	resSite, err := h.Usecase.CreateSiteGroup(&site)
	if err != nil {
		if err == utils.ErrConflict {
			return utils.NewError(c, http.StatusConflict, err)
		}
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, resSite)
}

// UpdateSiteGroup godoc
// @Summary Update Site By Id
// @Description Update Site By Id
// @Tags sites
// @Accept  json
// @Produce  json
// @Param  id path string true "Site ID"
// @Param sites body models.Site true  "id can empty"
// @Success 200 {object} models.Site
// @Failure 400 {object} utils.HTTPError
// @Failure 403 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /sites/{id} [put]
func (h *HttpHandler) UpdateSiteGroup(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsUserManager() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	site := models.Site{
		ID: bsonutil.MustObjectIDFromHex(id),
	}
	if err := c.Bind(&site); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}

	if ok, err := isRequestValid(&site); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	res, err := h.Usecase.UpdateSiteGroup(&site)
	if err != nil {
		if err == utils.ErrConflict {
			return utils.NewError(c, http.StatusConflict, err)
		}
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// GetSiteGroup godoc
// @Summary List Site
// @Description get site list
// @Tags sites
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Site
// @Failure 500 {object} utils.HTTPError
// @Router /sites [get]
func (h *HttpHandler) GetSiteGroup(c echo.Context) error {
	site, err := h.Usecase.GetSiteGroup()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, site)
}

// GetSiteGroupByID godoc
// @Summary Get Site By Id
// @Description Get Site By Id
// @Tags sites
// @Accept json
// @Produce json
// @Param id path string true "Site ID"
// @Success 200 {object} models.Site
// @Failure 204 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Router /sites/{id} [get]
func (h *HttpHandler) GetSiteGroupByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	site, err := h.Usecase.GetSiteGroupByID(id)
	if err != nil {
		return utils.NewError(c, http.StatusNoContent, err)
	}
	return c.JSON(http.StatusOK, site)
}

// DeleteSiteGroup godoc
// @Summary Delete Site
// @Description Delete Site By Id
// @Tags sites
// @Accept json
// @Produce json
// @Param id path string true "Site ID"
// @Success 204 {object} models.Site
// @Failure 400 {object} utils.HTTPError
// @Failure 403 {object} utils.HTTPError
// @Failure 409 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /sites/{id} [delete]
func (h *HttpHandler) DeleteSiteGroup(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsUserManager() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	err := h.Usecase.DeleteSiteGroup(id)
	if err != nil {
		if err == utils.ErrSiteInUse {
			return utils.NewError(c, http.StatusConflict, err)
		}
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	ur := NewRepository(session)
	uc := NewUsecase(ur, newUserBySiteReader(session))
	handler := &HttpHandler{uc}
	r = r.Group("/sites")
	r.POST("", handler.CreateSiteGroup)
	r.PUT("/:id", handler.UpdateSiteGroup)
	r.GET("", handler.GetSiteGroup)
	r.GET("/:id", handler.GetSiteGroupByID)
	r.DELETE("/:id", handler.DeleteSiteGroup)
}
