package login

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/httputil"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
)

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	userRepo := user.NewRepository(session)
	r.POST("/login", func(c echo.Context) error {
		return login(c, userRepo)
	})
}

// Login godoc
// @Summary Login
// @Description Login get token
// @Tags login
// @Accept  json
// @Produce  json
// @Param login body models.Login true  "id is userId"
// @Success 200 {object} models.Token
// @Failure 401 {object} httputil.HTTPError
// @Router /api/v1/login [post]
func login(c echo.Context, userRepo user.Repository) error {
	var u models.Login
	if err := c.Bind(&u); err != nil {
		return httputil.NewError(c, http.StatusUnauthorized, err)
	}

	user, err := userRepo.GetUserByID(u.ID)
	if err != nil {
		return httputil.NewError(c, http.StatusUnauthorized, err)
	}

	user.BankAccountName = ""
	user.BankAccountNumber = ""
	user.ThaiCitizenID = ""

	claims := &models.JwtCustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("GmkZGF3CmpZNs88dLvbV"))
	if err != nil {
		return httputil.NewError(c, http.StatusUnauthorized, err)
	}
	tk := &models.Token{
		Token: t,
	}
	return c.JSON(http.StatusOK, tk)
}
