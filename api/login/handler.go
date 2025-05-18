package login

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.odds.team/worklog/api.odds-worklog/api/consumer"
	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/mongo"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type HttpHandler struct {
	Usecase Usecase
}

func NewHttpHandler(r *echo.Group, session *mongo.Session) {
	siteRepo := site.NewRepository(session)
	userRepo := user.NewRepository(session)
	userUsecase := user.NewUsecase(userRepo, siteRepo)
	consumerUsecase := consumer.NewUsecase(consumer.NewRepository(session))
	loginUsecase := NewUsecase(userUsecase, consumerUsecase)
	handler := &HttpHandler{loginUsecase}

	r.POST("/login", func(c echo.Context) error {
		return login(c, userRepo)
	})
	r.POST("/login-google", handler.loginGoogle)
	r.POST("/login-keycloak", handler.loginKeycloak)
}

// loginKeycloak godoc
// @Summary Login with keycloak
// @Description Login with keycloak account only those with group worklog
// @Tags login
// @Accept  json
// @Produce  json
// @Param login body models.Login true  "id is token from Identity Provider from font-end"
// @Success 200 {object} models.Token
// @Failure 401 {object} utils.HTTPError
// @Router /login-keycloak [post]
func (h *HttpHandler) loginKeycloak(c echo.Context) error {
	var login models.Login
	if err := c.Bind(&login); err != nil {
		return utils.NewError(c, http.StatusUnauthorized, utils.ErrBadRequest)
	}

	if login.Token == "" {
		return utils.NewError(c, http.StatusUnauthorized, utils.ErrBadRequest)
	}

	tokenInfo, err := h.Usecase.ValidateAndExtractToken(login.Token)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	user, err := h.Usecase.CreateUser(tokenInfo.Email)
	if err != nil && err != utils.ErrConflict {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	token, err := handleToken(user)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	return c.JSON(http.StatusOK, token)
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
func (h *HttpHandler) loginGoogle(c echo.Context) error {
	var login models.Login
	if err := c.Bind(&login); err != nil {
		return utils.NewError(c, http.StatusUnauthorized, utils.ErrBadRequest)
	}

	if login.Token == "" {
		return utils.NewError(c, http.StatusUnauthorized, utils.ErrBadRequest)
	}

	tokenInfo, err := h.Usecase.GetTokenInfo(login.Token)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	user, err := h.Usecase.CreateUserAndValidateEmail(tokenInfo.Email)
	if err != nil && err != utils.ErrConflict {
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

	user, err := userRepo.GetByID(login.Token)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}

	token, err := handleToken(user)
	if err != nil {
		return utils.NewError(c, http.StatusUnauthorized, err)
	}
	return c.JSON(http.StatusOK, token)
}
