package customer

import (
	"errors"
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	validator "gopkg.in/go-playground/validator.v9"
	"gopkg.in/mgo.v2/bson"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type HttpHandler struct {
	Usecase Usecase
}

func isRequestValid(m *models.Customer) (bool, error) {
	if err := validator.New().Struct(m); err != nil {
		return false, err
	}
	return true, nil
}

// CreateCustomer godoc
// @Summary Create Customer
// @Description Create Customer
// @Tags customer
// @Accept json
// @Produce json
// @Param customer body models.IncomeReq true  " "
// @Success 200 {object} models.Customer
// @Failure 400 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /customer [post]
func (h *HttpHandler) CreateCustomer(c echo.Context) error {
	isAdmin, message := IsUserAdmin(c)
	if !isAdmin {
		return c.JSON(http.StatusUnauthorized, message)
	}

	var customer models.Customer
	if err := c.Bind(&customer); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}
	if ok, err := isRequestValid(&customer); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	res, err := h.Usecase.CreateCustomer(&customer)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

// GetCustomers godoc
// @Summary List customer
// @Description get customer list
// @Tags customers
// @Accept  json
// @Produce  json
// @Success 200 {array} models.Customer
// @Failure 500 {object} utils.HTTPError
// @Router /customer [get]
func (h *HttpHandler) GetCustomers(c echo.Context) error {
	checkUser, message := IsUserAdmin(c)
	if !checkUser {
		return c.JSON(http.StatusUnauthorized, message)
	}
	customers, err := h.Usecase.GetCustomers()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, customers)
}

// UpdateCustomerById godoc
// @Summary Update Customer By Id
// @Description Update Customer By Id
// @Tags customer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Param customer body models.Customer true  "id can empty"
// @Success 200 {object} models.Customer
// @Failure 400 {object} utils.HTTPError
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /customer/{id} [put]
func (h *HttpHandler) UpdateCustomer(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	u := models.Customer{
		ID: bson.ObjectIdHex(id),
	}
	if err := c.Bind(&u); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}

	if ok, err := isRequestValid(&u); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	ut := getUserFromToken(c)
	user, err := h.Usecase.UpdateCustomer(&u, ut.IsAdmin())
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, user)
}

// GetCustomerById godoc
// @Summary Get Customer By Id
// @Description Get Customer By Id
// @Tags customers
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 200 {object} models.Customer
// @Failure 204 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Router /customer/{id} [get]
func (h *HttpHandler) GetCustomerByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	user, err := h.Usecase.GetCustomerByID(id)
	if err != nil {
		return utils.NewError(c, http.StatusNoContent, err)
	}
	return c.JSON(http.StatusOK, user)
}

// DeleteCustomer godoc
// @Summary Delete Customer
// @Description Delete Customer By Id
// @Tags customer
// @Accept json
// @Produce json
// @Param id path string true "Customer ID"
// @Success 204 {object} models.Customer
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /customer/{id} [delete]
func (h *HttpHandler) DeleteCustomer(c echo.Context) error {
	checkUser, message := IsUserAdmin(c)
	if !checkUser {
		return c.JSON(http.StatusUnauthorized, message)
	}
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}

	err := h.Usecase.DeleteCustomer(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func getUserFromToken(c echo.Context) *models.User {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}
func IsUserAdmin(c echo.Context) (bool, string) {
	u := getUserFromToken(c)
	if u.IsAdmin() {
		return true, ""
	}
	return false, "ไม่มีสิทธิในการใช้งาน"
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	customerRepo := NewRepository(session)
	userRepo := user.NewRepository(session)
	uc := NewUsecase(customerRepo, userRepo)

	handler := &HttpHandler{uc}

	r = r.Group("/customer")
	r.POST("", handler.CreateCustomer)
	r.PUT("/:id", handler.UpdateCustomer)
	r.GET("", handler.GetCustomers)
	r.GET("/:id", handler.GetCustomerByID)
	r.DELETE("/:id", handler.DeleteCustomer)

}
