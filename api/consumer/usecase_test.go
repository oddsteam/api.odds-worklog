package consumer

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	consumerMock "gitlab.odds.team/worklog/api.odds-worklog/api/consumer/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestGetByClientID(t *testing.T) {
	t.Run("when get client ID with old hardcode client id, then return consumer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		defaultClientID := "956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com"
		mockRepo := consumerMock.NewMockRepository(ctrl)

		usecase := NewUsecase(mockRepo)
		mockRepo.EXPECT().GetByClientID(defaultClientID).Return(nil, utils.ErrInvalidConsumer)
		mockRepo.EXPECT().GetByClientID(defaultClientID).Return(&models.Consumer{ClientID: defaultClientID}, nil).After(mockRepo.EXPECT().Create(&models.Consumer{ClientID: defaultClientID}).Return(&models.Consumer{ClientID: defaultClientID}, nil))

		consumer, err := usecase.GetByClientID(defaultClientID)

		assert.NoError(t, err)
		assert.Equal(t, defaultClientID, consumer.ClientID)
	})

	t.Run("when get client ID with stored client id, then return consumer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		defaultClientID := "956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com"
		clientID := "test.apps.googleusercontent.com"
		mockRepo := consumerMock.NewMockRepository(ctrl)

		usecase := NewUsecase(mockRepo)
		mockRepo.EXPECT().GetByClientID(defaultClientID).Return(&models.Consumer{ClientID: defaultClientID}, nil)
		mockRepo.EXPECT().GetByClientID(clientID).Return(&models.Consumer{ClientID: clientID}, nil)

		consumer, err := usecase.GetByClientID(clientID)

		assert.NoError(t, err)
		assert.Equal(t, clientID, consumer.ClientID)
	})

	t.Run("when get client ID with NOT stored client id, then return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		defaultClientID := "956316396976-mhb092ad69gn2olis0mtmc1fpe8blgn8.apps.googleusercontent.com"
		clientID := "test.apps.googleusercontent.com"
		mockRepo := consumerMock.NewMockRepository(ctrl)

		usecase := NewUsecase(mockRepo)
		mockRepo.EXPECT().GetByClientID(defaultClientID).Return(&models.Consumer{ClientID: defaultClientID}, nil)
		mockRepo.EXPECT().GetByClientID(gomock.Not(gomock.Eq(clientID))).Return(nil, utils.ErrInvalidConsumer)

		_, err := usecase.GetByClientID("invalid")

		assert.Equal(t, err, utils.ErrInvalidConsumer)
	})
}
