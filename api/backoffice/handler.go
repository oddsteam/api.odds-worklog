package backoffice

import (
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"

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

	r = r.Group("/backoffice")
	r.GET("", handler.Get)
}

func getUserFromToken(c echo.Context) *models.UserClaims {
	t := c.Get("user").(*jwt.Token)
	claims := t.Claims.(*models.JwtCustomClaims)
	return claims.User
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
// @Router /backoffice [get]
func (h *HttpHandler) Get(c echo.Context) error {
	// user := getUserFromToken(c)
	// if !user.IsAdmin() {
	// 	return utils.NewError(c, http.StatusForbidden, utils.ErrPermissionDenied)
	// }
	users, err := h.Usecase.Get()
	if err != nil {
		return utils.NewError(c, http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, users)
}
