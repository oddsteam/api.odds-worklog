package backoffice

import (
	"net/http"

	"gitlab.odds.team/worklog/api.odds-worklog/api/site"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
	"fmt"
	"crypto/sha256"
    "encoding/hex"
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
	r.POST("", handler.GetAllUserIncome)
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
func (h *HttpHandler) GetAllUserIncome(c echo.Context) error {
	token := c.QueryParam("token");
	s := "backoffice"
    hasher := sha256.New()
    hasher.Write([]byte(s))
    sha1_hash := hex.EncodeToString(hasher.Sum(nil))
	fmt.Println(s, sha1_hash)
	
	if token == sha1_hash {
		users, err := h.Usecase.Get()
		if err != nil {
			return utils.NewError(c, http.StatusInternalServerError, err)
		}
		return c.JSON(http.StatusOK, users)
	}else{
		return c.JSON(http.StatusOK, "invalid token")
	}
}