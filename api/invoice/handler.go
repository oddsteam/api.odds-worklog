package invoice

import (
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/api/po"

	"gitlab.odds.team/worklog/api.odds-worklog/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type HttpHandler struct {
	usecase Usecase
}

func NewHttpHandler(g *echo.Group, s *mongo.Session) {
	invoiceRepo := NewRepository(s)
	poRepo := po.NewRepository(s)
	use := NewUsecase(invoiceRepo, poRepo)
	h := &HttpHandler{use}

	g = g.Group("/invoices")
	g.POST("", h.Create)
	g.GET("", h.Get)
	g.GET("/po/next-no", h.NextNo)
	g.GET("/po/:id", h.GetByPO)
}

func getUserFromToken(c echo.Context) *models.User {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}

// Create godoc
// @Summary Create New Invoice
// @Description Create New Invoice
// @Tags invoices
// @Accept json
// @Produce json
// @Param invoice body models.Invoice true  "id can be empty"
// @Success 200 {object} models.Invoice
// @Failure 403 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /invoices [post]
func (h *HttpHandler) Create(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	var i models.Invoice
	if err := c.Bind(&i); err != nil {
		return utils.NewError(c, http.StatusBadRequest, err)
	}

	invoice, err := h.usecase.Create(&i)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, invoice)
}

// Get godoc
// @Summary Get Invoice List
// @Description Get Invoice List
// @Tags invoices
// @Produce json
// @Success 200 {array} models.Invoice
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /invoices [get]
func (h *HttpHandler) Get(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	invoices, err := h.usecase.Get()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, invoices)
}

// GetByPO godoc
// @Summary Get Invoice List by PO
// @Description Get Invoice List by PO
// @Tags invoices
// @Produce json
// @Success 200 {array} models.Invoice
// @Failure 403 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /invoices/po/{id} [get]
func (h *HttpHandler) GetByPO(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	id := c.Param("id")
	invoices, err := h.usecase.GetByPO(id)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, invoices)
}

// NextNo godoc
// @Summary Get next invoice no
// @Description Get next invoice no
// @Tags invoices
// @Accept json
// @Produce json
// @Param invoiceNoReq body models.InvoiceNoReq true  "poId don't empty"
// @Success 200 {object} models.InvoiceNoRes
// @Failure 403 {object} utils.HTTPError
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /invoices/po/next-no [get]
func (h *HttpHandler) NextNo(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	var req models.InvoiceNoReq
	err := c.Bind(&req)
	if req.PoID == "" {
		return utils.NewError(c, http.StatusBadRequest, utils.ErrInvalidPath)
	}

	invoiceNo, err := h.usecase.NextNo(req.PoID)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, models.InvoiceNoRes{InvoiceNo: invoiceNo})
}
