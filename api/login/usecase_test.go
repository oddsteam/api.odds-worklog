package login

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestVerifyAudience(t *testing.T) {
	t.Run("956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com is currect Audience", func(t *testing.T) {
		assert.Equal(t, clientID, "956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com")
	})

	t.Run("956316396976-cnrmemp4r4coc62oqmn9uin7iq3o3eev.apps.googleusercontent.com2 is wrong Audience", func(t *testing.T) {
		assert.NotEqual(t, clientID, "956316396976-cnrmemp4r4coc62oqmn9uin7iq3o3eev.apps.googleusercontent.com")
	})
}

func TestHandleToken(t *testing.T) {
	t.Run("when user is first login got FirstLogin = 'Y'", func(t *testing.T) {
		user := new(models.User)
		token, err := handleToken(user)

		assert.NoError(t, err)
		assert.Equal(t, "Y", token.FirstLogin)
	})

	t.Run("when user is't first login got FirstLogin = 'N'", func(t *testing.T) {
		token, err := handleToken(&userMock.MockUser)

		assert.NoError(t, err)
		assert.Equal(t, "N", token.FirstLogin)
	})
}

func TestGenToken(t *testing.T) {
	t.Run("generate token success", func(t *testing.T) {
		token, err := genToken(&userMock.MockUser)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("create user success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		email := "abc@mail.com"
		user := new(models.User)
		user.Email = email
		user.CorporateFlag = "F"

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().CreateUser(gomock.Any()).Return(user, nil)

		usecase := NewUsecase(mockUsecase)
		userRes, err := usecase.CreateUser(email)

		assert.NoError(t, err)
		assert.Equal(t, email, userRes.Email)
		assert.Equal(t, "F", userRes.CorporateFlag)
	})
}
