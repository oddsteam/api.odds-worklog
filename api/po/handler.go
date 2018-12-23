package po

import (
	"errors"
	"net/http"

	"gopkg.in/mgo.v2/bson"

	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	validator "gopkg.in/go-playground/validator.v9"
)

type HttpHandler struct {
	Usecase Usecase
}

func isRequestValid(m *models.Po) (bool, error) {
	if err := validator.New().Struct(m); err != nil {
		return false, err
	}
	return true, nil
}

// Create godoc
// @Summary Create PO
// @Description Create PO
// @Tags po
// @Accept  json
// @Produce  json
// @Param poes body models.Po true "require customer id"
// @Success 200 {object} models.Po
// @Failure 500 {object} utils.HTTPError
// @Router /poes [post]
func (h *HttpHandler) Create(c echo.Context) error {
	var po models.Po
	if err := c.Bind(&po); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}
	if ok, err := isRequestValid(&po); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}
	resPo, err := h.Usecase.Create(&po)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, resPo)

}

// Update godoc
// @Summary Update PO
// @Description Update PO
// @Tags po
// @Accept  json
// @Produce  json
// @Param  id path string true "Po ID"
// @Param poes body models.Po
// @Success 200 {object} models.Po
// @Failure 500 {object} utils.HTTPError
// @Router /poes/{id} [put]
func (h *HttpHandler) Update(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return utils.NewError(c, http.StatusBadRequest, errors.New("invalid path"))
	}
	po := models.Po{
		ID: bson.ObjectIdHex(id),
	}
	if err := c.Bind(&po); err != nil {
		return utils.NewError(c, http.StatusUnprocessableEntity, err)
	}
	if ok, err := isRequestValid(&po); !ok {
		return utils.NewError(c, http.StatusBadRequest, err)
	}
	res, err := h.Usecase.Update(&po)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, res)
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	ur := NewRepository(session)
	uc := NewUsecase(ur)
	handler := &HttpHandler{uc}
	r = r.Group("/poes")
	r.POST("", handler.Create)
	r.PUT("/:id", handler.Update)
}
