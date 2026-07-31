package site

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	mock "gitlab.odds.team/worklog/api.odds-worklog/api/site/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestUsecase_CreateSiteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockUsers := mock.NewMockForGettingUsersBySiteID(ctrl)
	mockRepo.EXPECT().GetSiteGroupByName(mock.MockSite.Name).Return(nil, nil)
	mockRepo.EXPECT().CreateSiteGroup(&mock.MockSite).Return(&mock.MockSite, nil)

	uc := NewUsecase(mockRepo, mockUsers)
	site, err := uc.CreateSiteGroup(&mock.MockSite)

	assert.NoError(t, err)
	assert.Equal(t, mock.MockSite.ID.Hex(), site.ID.Hex())
	assert.Equal(t, mock.MockSite.Name, site.Name)
}

func TestUsecase_UpdateSiteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockUsers := mock.NewMockForGettingUsersBySiteID(ctrl)
	mockRepo.EXPECT().GetSiteGroupByName(mock.MockSite.Name).Return(&mock.MockSite, nil)
	mockRepo.EXPECT().UpdateSiteGroup(&mock.MockSite).Return(&mock.MockSite, nil)

	uc := NewUsecase(mockRepo, mockUsers)
	site, err := uc.UpdateSiteGroup(&mock.MockSite)

	assert.NoError(t, err)
	assert.Equal(t, mock.MockSite.ID.Hex(), site.ID.Hex())
	assert.Equal(t, mock.MockSite.Name, site.Name)
}

func TestUsecase_UpdateSiteGroup_Conflict(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockUsers := mock.NewMockForGettingUsersBySiteID(ctrl)
	mockRepo.EXPECT().GetSiteGroupByName(mock.MockSite.Name).Return(&mock.MockSite2, nil)

	uc := NewUsecase(mockRepo, mockUsers)
	site, err := uc.UpdateSiteGroup(&mock.MockSite)

	assert.Nil(t, site)
	assert.EqualError(t, err, utils.ErrConflict.Error())
}

func TestUsecase_GetSiteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockUsers := mock.NewMockForGettingUsersBySiteID(ctrl)
	mockRepo.EXPECT().GetSiteGroup().Return(mock.MockSites, nil)

	uc := NewUsecase(mockRepo, mockUsers)
	site, err := uc.GetSiteGroup()

	assert.NoError(t, err)
	assert.Equal(t, 2, len(site))
	assert.Equal(t, mock.MockSites[0].ID.Hex(), site[0].ID.Hex())
	assert.Equal(t, mock.MockSites[0].Name, site[0].Name)
	assert.Equal(t, mock.MockSites[1].ID.Hex(), site[1].ID.Hex())
	assert.Equal(t, mock.MockSites[1].Name, site[1].Name)
}

func TestUsecase_GetSiteGroupById(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockUsers := mock.NewMockForGettingUsersBySiteID(ctrl)
	mockRepo.EXPECT().GetSiteGroupByID(mock.MockSite.ID.Hex()).Return(&mock.MockSite, nil)

	uc := NewUsecase(mockRepo, mockUsers)
	site, err := uc.GetSiteGroupByID(mock.MockSite.ID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, mock.MockSite.ID.Hex(), site.ID.Hex())
	assert.Equal(t, mock.MockSite.Name, site.Name)
}

func TestUsecase_DeleteSiteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockUsers := mock.NewMockForGettingUsersBySiteID(ctrl)
	mockUsers.EXPECT().GetBySiteID(mock.MockSite.ID.Hex()).Return([]*models.User{}, nil)
	mockRepo.EXPECT().DeleteSiteGroup(mock.MockSite.ID.Hex()).Return(nil)

	uc := NewUsecase(mockRepo, mockUsers)
	err := uc.DeleteSiteGroup(mock.MockSite.ID.Hex())

	assert.NoError(t, err)
}

func TestUsecase_DeleteSiteGroup_InUse(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockUsers := mock.NewMockForGettingUsersBySiteID(ctrl)
	mockUsers.EXPECT().GetBySiteID(mock.MockSite.ID.Hex()).Return([]*models.User{{}}, nil)

	uc := NewUsecase(mockRepo, mockUsers)
	err := uc.DeleteSiteGroup(mock.MockSite.ID.Hex())

	assert.EqualError(t, err, utils.ErrSiteInUse.Error())
}
