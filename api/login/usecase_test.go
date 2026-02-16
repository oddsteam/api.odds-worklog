package login

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestHandleToken(t *testing.T) {
	t.Run("when user is first login got FirstLogin = 'Y'", func(t *testing.T) {
		user := new(models.User)
		token, err := handleToken(user)

		assert.NoError(t, err)
		assert.Equal(t, "Y", token.FirstLogin)
	})

	t.Run("when user is't first login got FirstLogin = 'N'", func(t *testing.T) {
		token, err := handleToken(&userMock.User)

		assert.NoError(t, err)
		assert.Equal(t, "N", token.FirstLogin)
	})
}

func TestGenToken(t *testing.T) {
	t.Run("generate token success", func(t *testing.T) {
		token, err := genToken(
			&models.UserClaims{
				ID:         userMock.User.ID.Hex(),
				Role:       userMock.User.Role,
				StatusTavi: userMock.User.StatusTavi,
			},
		)

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("create user success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		email := "abc@odds.team"
		user := new(models.User)
		user.Email = email

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Create(gomock.Any()).Return(user, nil)

		usecase := NewUsecase(mockUsecase)
		userRes, err := usecase.CreateUserAndValidateEmail(email)

		assert.NoError(t, err)
		assert.Equal(t, email, userRes.Email)
	})

	t.Run("when email is not odds.team, then return error ErrEmailIsNotOddsTeam", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		email := "abc@mail.com"
		user := new(models.User)
		user.Email = email

		mockUsecase := userMock.NewMockUsecase(ctrl)

		usecase := NewUsecase(mockUsecase)
		userRes, err := usecase.CreateUserAndValidateEmail(email)

		assert.Nil(t, userRes)
		assert.EqualError(t, err, utils.ErrEmailIsNotOddsTeam.Error())
	})
}

func TestUsecase_isOddsTeam(t *testing.T) {
	assert.True(t, isOddsTeam("a@odds.team"))
	assert.False(t, isOddsTeam("a@xyzodds.team"))
	assert.False(t, isOddsTeam(""))
	assert.False(t, isOddsTeam("a@gmail.com"))
}

