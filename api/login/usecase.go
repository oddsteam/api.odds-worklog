package login

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	oauth2 "google.golang.org/api/oauth2/v2"
)

const clientID = "956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com"

type usecase struct {
	UserUsecase user.Usecase
}

func NewUsecase(uu user.Usecase) Usecase {
	return &usecase{uu}
}

func (u *usecase) GetTokenInfo(idToken string) (*oauth2.Tokeninfo, error) {
	oauth2Service, err := oauth2.New(&http.Client{})
	if err != nil {
		return nil, err
	}

	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfo, err := tokenInfoCall.IdToken(idToken).Do()

	if err != nil {
		return nil, err
	}

	if !verifyAudience(tokenInfo.Audience) {
		return nil, errors.New("Token is not account @odds.team")
	}
	return tokenInfo, nil
}

func (u *usecase) CreateUser(email string) (*models.User, error) {
	user := &models.User{}
	user.Email = email
	user.CorporateFlag = "F"
	return u.UserUsecase.CreateUser(user)
}

func verifyAudience(aud string) bool {
	return aud == clientID
}

func handleToken(user *models.User) (*models.Token, error) {
	user.BankAccountName = ""
	user.BankAccountNumber = ""
	user.ThaiCitizenID = ""

	tok, err := genToken(user)
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
		IdUser:     user.ID.Hex(),
	}

	return token, nil
}

func genToken(user *models.User) (string, error) {
	claims := &models.JwtCustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tok, err := token.SignedString([]byte("GmkZGF3CmpZNs88dLvbV"))
	if err != nil {
		return "", fmt.Errorf("Generate token error: %s", err.Error())
	}
	return tok, nil
}
