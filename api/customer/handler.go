package customer

import (
	"net/http"

	"github.com/globalsign/mgo/bson"
	"gitlab.odds.team/worklog/api.odds-worklog/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type HttpHandler struct {
	Usecase Usecase
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	repo := NewRepository(session)
	uc := NewUsecase(repo)
	handler := &HttpHandler{uc}

	r = r.Group("/customers")
	r.POST("", handler.Create)
	r.PUT("/:id", handler.Update)
	r.GET("", handler.Get)
	r.GET("/:id", handler.GetByID)
	r.DELETE("/:id", handler.Delete)
}

func getUserFromToken(c echo.Context) *models.User {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}

// Create godoc
// @Summary Create Customer
// @Description Create Customer
// @Tags customers
// @Accept json
// @Produce json
// @Param customer body models.Customer true  " "
// @Success 200 {object} models.Customer
// @Failure 403 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /customers [post]
func (h *HttpHandler) Create(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	var customer models.Customer
	if err := c.Bind(&customer); err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	res, err := h.Usecase.Create(&customer)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// Get godoc
// @Summary List customer
// @Description get customer list
// @Tags customers
// @Accept json
// @Produce json
// @Success 200 {array} models.Customer
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /customers [get]
func (h *HttpHandler) Get(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}
	customers, err := h.Usecase.Get()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, customers)
}

// Update godoc
// @Summary Update Customer
// @Description Update Customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body models.Customer true  "id can empty"
// @Success 200 {object} models.Customer
// @Failure 403 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /customers/{id} [put]
func (h *HttpHandler) Update(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	var ct models.Customer
	if err := c.Bind(&ct); err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	id := c.Param("id")
	ct.ID = bson.ObjectIdHex(id)
	customer, err := h.Usecase.Update(&ct)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, customer)
}

// GetById godoc
// @Summary Get Customer By Id
// @Description Get Customer By Id
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Customer
// @Failure 204 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Router /customers/{id} [get]
func (h *HttpHandler) GetByID(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	id := c.Param("id")
	customer, err := h.Usecase.GetByID(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, customer)
}

// Delete godoc
// @Summary Delete Customer
// @Description Delete Customer
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Response
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /customers/{id} [delete]
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
	return c.JSON(http.StatusOK, models.Response{Message: "Delete customer success."})
}
