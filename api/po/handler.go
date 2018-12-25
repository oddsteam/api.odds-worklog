package po

import (
	"net/http"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/customer"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"gopkg.in/mgo.v2/bson"
)

type HttpHandler struct {
	Usecase Usecase
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	ur := NewRepository(session)
	custRepo := customer.NewRepository(session)
	uc := NewUsecase(ur, custRepo)
	handler := &HttpHandler{uc}

	r = r.Group("/poes")
	r.POST("", handler.Create)
	r.PUT("/:id", handler.Update)
	r.GET("", handler.Get)
	r.GET("/:id", handler.GetByID)
	r.GET("/customer/:id", handler.GetByCusID)
	r.DELETE("/:id", handler.Delete)
}

func getUserFromToken(c echo.Context) *models.User {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}

// Create godoc
// @Summary Create PO
// @Description Create PO
// @Tags po
// @Accept json
// @Produce json
// @Param po body models.Po true "customer id is require"
// @Success 200 {object} models.Po
// @Failure 203 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /poes [post]
func (h *HttpHandler) Create(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	var po models.Po
	if err := c.Bind(&po); err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}
	resPo, err := h.Usecase.Create(&po)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, resPo)
}

// Update godoc
// @Summary Update PO
// @Description Update PO
// @Tags po
// @Accept json
// @Produce json
// @Param id path string true "id is Po ID"
// @Param po body models.Po true "customer id is require"
// @Success 200 {object} models.Po
// @Failure 203 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /poes/{id} [put]
func (h *HttpHandler) Update(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	var po models.Po
	if err := c.Bind(&po); err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	id := c.Param("id")
	po.ID = bson.ObjectIdHex(id)
	res, err := h.Usecase.Update(&po)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// Get godoc
// @Summary Get Po List
// @Description Get Po List
// @Tags po
// @Accept json
// @Produce json
// @Success 200 {array} models.Po
// @Failure 203 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /poes [get]
func (h *HttpHandler) Get(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	res, err := h.Usecase.Get()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// GetByID godoc
// @Summary Get Po
// @Description Get Po By Id
// @Tags po
// @Accept json
// @Produce json
// @Param id path string true "Po ID"
// @Success 200 {object} models.Po
// @Failure 203 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /poes/{id} [get]
func (h *HttpHandler) GetByID(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	id := c.Param("id")
	po, err := h.Usecase.GetByID(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, po)
}

// GetByCusID godoc
// @Summary Get Poes By Customer Id
// @Description Get Get Poes By Customer Id
// @Tags po
// @Accept json
// @Produce json
// @Param id path string true "Poes ID"
// @Success 200 {object} models.Po
// @Failure 203 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /poes/customer/{id} [get]
func (h *HttpHandler) GetByCusID(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}
	id := c.Param("id")
	po, err := h.Usecase.GetByCusID(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, po)
}

// Delete godoc
// @Summary Delete Po
// @Description Delete Po
// @Tags po
// @Accept json
// @Produce json
// @Param id path string true "Po ID"
// @Success 200 {object} models.Response
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /poes/{id} [delete]
func (h *HttpHandler) Delete(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}
	id := c.Param("id")
	err := h.Usecase.Delete(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, models.Response{Message: "Delete product owner success."})
}
