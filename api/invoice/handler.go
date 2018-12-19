package invoice

import (
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/models"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type HttpHandler struct {
	usecase Usecase
}

func NewHttpHandler(g *echo.Group, s *mongo.Session) {
	repo := NewRepository(s)
	use := NewUsecase(repo)
	h := &HttpHandler{use}

	g = g.Group("/invoices")
	g.POST("", h.Create)
}

// Create godoc
// @Summary Create New Invoice
// @Description Create New Invoice
// @Tags invoices
// @Accept json
// @Produce json
// @Param invoice body models.Invoice true  "id can be empty"
// @Success 200 {array} models.Invoice
// @Failure 422 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /invoices [post]
func (h *HttpHandler) Create(c echo.Context) error {
	var i models.Invoice
	if err := c.Bind(&i); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}

	invoice, err := h.usecase.Create(&i)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, invoice)
}
