package login

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	consumerMock "gitlab.odds.team/worklog/api.odds-worklog/api/consumer/mock"
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

		mockConsumerUsecase := consumerMock.NewMockUsecase(ctrl)

		usecase := NewUsecase(mockUsecase, mockConsumerUsecase)
		userRes, err := usecase.CreateUser(email)

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
		mockConsumerUsecase := consumerMock.NewMockUsecase(ctrl)

		usecase := NewUsecase(mockUsecase, mockConsumerUsecase)
		userRes, err := usecase.CreateUser(email)

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

func TestValidConsumerClientID(t *testing.T) {
	t.Run("when consumer client id is stored, then return valid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockConsumerUsecase := consumerMock.NewMockUsecase(ctrl)
		mockConsumerUsecase.EXPECT().GetByClientID(gomock.Eq("ThisIsValidClientID")).Return(&models.Consumer{}, nil)

		usecase := NewUsecase(mockUsecase, mockConsumerUsecase)
		isValid := usecase.IsValidConsumerClientID("ThisIsValidClientID")

		assert.True(t, isValid)
	})

	t.Run("when consumer client id is NOT stored, then return INvalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockConsumerUsecase := consumerMock.NewMockUsecase(ctrl)
		mockConsumerUsecase.EXPECT().GetByClientID(gomock.Not(gomock.Eq("ThisIsValidClientID"))).Return(nil, utils.ErrInvalidConsumer)

		usecase := NewUsecase(mockUsecase, mockConsumerUsecase)
		isValid := usecase.IsValidConsumerClientID("ThisIs InValid ClientID")

		assert.False(t, isValid)
	})
}
