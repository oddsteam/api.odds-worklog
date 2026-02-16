package login

import (
	"fmt"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/auth"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	UserUsecase user.Usecase
}

func NewUsecase(uu user.Usecase) Usecase {
	return &usecase{uu}
}

func (u *usecase) ValidateAndExtractToken(accessToken string) (models.Identity, error) {
	validator := auth.NewKeycloakValidator(
		os.Getenv("KEYCLOAK_SERVER_URL"),
		os.Getenv("KEYCLOAK_REALM"),
		os.Getenv("KEYCLOAK_CLIENT_ID"),
	)

	claims, err := validator.ValidateToken(accessToken)

	if err != nil {
		return models.Identity{}, err
	}

	// Check if journeyman role exists
	if !contains(claims.ResourceAccess["worklog"].Roles, "journeyman") {
		return models.Identity{}, fmt.Errorf("user does not have journeyman role")
	}

	return models.Identity{Email: claims.Email}, nil
}

// contains checks if a string exists in a slice of strings
func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func (u *usecase) CreateUserAndValidateEmail(email string) (*models.User, error) {
	if !isOddsTeam(email) {
		return nil, utils.ErrEmailIsNotOddsTeam
	}
	return u.CreateUser(email)
}

func (u *usecase) CreateUser(email string) (*models.User, error) {
	user := &models.User{}
	user.Email = email
	if user.IsAdmin() {
		user.Role = "admin"
	} else {
		user.Role = "individual"
	}
	return u.UserUsecase.Create(user)
}

func isOddsTeam(email string) bool {
	if len(email) < 10 {
		return false
	}

	host := email[len(email)-10:]
	return host == "@odds.team"
}

func handleToken(user *models.User) (*models.Token, error) {
	tok, err := genToken(
		&models.UserClaims{
			ID:         user.ID.Hex(),
			Role:       user.Role,
			StatusTavi: user.StatusTavi,
		},
	)
	if err != nil {
		return nil, err
	}

	firstLogin := "Y"
	if !user.IsFullnameEmpty() {
		firstLogin = "N"
	}

	token := &models.Token{
		Token:      tok,
		FirstLogin: firstLogin,
		User:       user,
	}

	return token, nil
}

func genToken(user *models.UserClaims) (string, error) {
	claims := &models.JwtCustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte("sMJuczqQPYzocl1s6SLj"))
	if err != nil {
		return "", fmt.Errorf("Generate token error: %s", err.Error())
	}
	return tok, nil
}
