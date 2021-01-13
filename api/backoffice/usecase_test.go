package backoffice

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	backofficeMock "gitlab.odds.team/worklog/api.odds-worklog/api/backoffice/mock"
	siteMock "gitlab.odds.team/worklog/api.odds-worklog/api/site/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestUsecaseGetUserIncome(t *testing.T) {
	t.Run("get all userIncomes success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBackoffice := backofficeMock.NewMockRepository(ctrl)
		mockRepoBackoffice.EXPECT().Get().Return(backofficeMock.MockUserIncomeList, nil)

		mockRepoSite := siteMock.NewMockRepository(ctrl)
		mockRepoSite.EXPECT().GetSiteGroup().Return(siteMock.MockSites, nil)

		usecase := NewUsecase(mockRepoBackoffice, mockRepoSite)
		userIncomeRes, err := usecase.Get()

		assert.NotNil(t, userIncomeRes)
		assert.Nil(t, err)

	})

	t.Run("get all userIncomes fail error when get all userIncome", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBackoffice := backofficeMock.NewMockRepository(ctrl)
		mockRepoBackoffice.EXPECT().Get().Return(nil,utils.ErrNotFound)

		mockRepoSite := siteMock.NewMockRepository(ctrl)

		usecase := NewUsecase(mockRepoBackoffice, mockRepoSite)
		userIncomeRes, err := usecase.Get()

		assert.Nil(t, userIncomeRes)
		assert.NotNil(t, err)

	})

	t.Run("get all userIncomes fail error when get all site", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepoBackoffice := backofficeMock.NewMockRepository(ctrl)
		mockRepoBackoffice.EXPECT().Get().Return(backofficeMock.MockUserIncomeList, nil)

		mockRepoSite := siteMock.NewMockRepository(ctrl)
		mockRepoSite.EXPECT().GetSiteGroup().Return(nil,utils.ErrNotFound)

		usecase := NewUsecase(mockRepoBackoffice, mockRepoSite)
		userIncomeRes, err := usecase.Get()

		assert.Nil(t, userIncomeRes)
		assert.NotNil(t, err)

	})

}
