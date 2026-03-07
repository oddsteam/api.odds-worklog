package income

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	incomeMock "gitlab.odds.team/worklog/api.odds-worklog/business/models/mock"
)

func TestUsecaseGetByRole(t *testing.T) {
	t.Run("when get by role success it should return users", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoIncome := incomeMock.NewMockRepository(ctrl)
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByRole("corporate").Return(userMock.Users, nil)

		uc := NewUsecase(mockRepoIncome, mockUserRepo)
		res, err := uc.GetByRole("corporate")
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, userMock.Users, res)
	})
}
