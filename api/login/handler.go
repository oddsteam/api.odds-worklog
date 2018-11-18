package login

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	userRepo := user.NewRepository(session)
	uc := newUsecase()

	r.POST("/login", func(c echo.Context) error {
		return login(c, userRepo)
	})

	r.POST("/login-google", func(c echo.Context) error {
		return loginGoogle(uc, c, userRepo)
	})
}

// loginGoogle godoc
// @Summary Login with google account
// @Description Login with google account only odds.team
// @Tags login
// @Accept  json
// @Produce  json
// @Param login body models.Login true  "id is token from google login in font-end"
// @Success 200 {object} models.Token
// @Failure 401 {object} utils.HTTPError
// @Router /login-google [post]
func loginGoogle(log Usecase, c echo.Context, userRepo user.Repository) error {
	var login models.Login
	if err := c.Bind(&login); err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	profile, err := log.GetTokenInfo(login.Token)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	if !verifyAudience(profile.Audience) {
		return utils.NewError(c, http.StatusUnauthorized, errors.New("Token is not account @odds.team"))
	}

	user := &models.User{}
	user.Email = profile.Email
	user.CorporateFlag = "F"

	user, err = userRepo.CreateUser(user)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	token, err := handleToken(user)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusOK, token)
}

// Login godoc
// @Summary Login
// @Description Login get token
// @Tags login
// @Accept  json
// @Produce  json
// @Param login body models.Login true  "id is userId"
// @Success 200 {object} models.Token
// @Failure 401 {object} utils.HTTPError
// @Router /login [post]
func login(c echo.Context, userRepo user.Repository) error {
	var login models.Login
	if err := c.Bind(&login); err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	user, err := userRepo.GetUserByID(login.Token)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	token, err := handleToken(user)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusOK, token)
}
