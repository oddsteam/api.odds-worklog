package sap_export_failure

import (
	"net/http"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/business/usecases"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"gitlab.odds.team/worklog/api.odds-worklog/repositories"
)

type HttpHandler struct {
	uc usecases.ForViewingSAPExportFailures
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	lister := repositories.NewSAPExportFailureLister(session)
	uc := usecases.NewViewSAPExportFailuresUsecase(lister)
	handler := &HttpHandler{uc: uc}

	g := r.Group("/sap-export-failures")
	g.GET("", handler.List)
}

func getUserFromToken(c echo.Context) *models.UserClaims {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
}

func (h *HttpHandler) List(c echo.Context) error {
	user := getUserFromToken(c)
	if !user.IsAdmin() {
		return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	}

	limit := 0
	if q := c.QueryParam("limit"); q != "" {
		if n, err := strconv.Atoi(q); err == nil {
			limit = n
		}
	}

	logs, err := h.uc.List(limit)
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, logs)
}
