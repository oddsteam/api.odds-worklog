package site

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	mock "gitlab.odds.team/worklog/api.odds-worklog/api/site/mock"
)

func TestUsecase_CreateSiteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetSiteGroupByName(mock.MockSite.Name).Return(nil, nil)
	mockRepo.EXPECT().CreateSiteGroup(&mock.MockSite).Return(&mock.MockSite, nil)

	uc := NewUsecase(mockRepo)
	site, err := uc.CreateSiteGroup(&mock.MockSite)

	assert.NoError(t, err)
	assert.Equal(t, mock.MockSite.ID.Hex(), site.ID.Hex())
	assert.Equal(t, mock.MockSite.Name, site.Name)
}

func TestUsecase_UpdateSiteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().UpdateSiteGroup(&mock.MockSite).Return(&mock.MockSite, nil)

	uc := NewUsecase(mockRepo)
	site, err := uc.UpdateSiteGroup(&mock.MockSite)

	assert.NoError(t, err)
	assert.Equal(t, mock.MockSite.ID.Hex(), site.ID.Hex())
	assert.Equal(t, mock.MockSite.Name, site.Name)
}

func TestUsecase_GetSiteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().GetSiteGroup().Return(mock.MockSites, nil)

	uc := NewUsecase(mockRepo)
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
	mockRepo.EXPECT().GetSiteGroupByID(mock.MockSite.ID.Hex()).Return(&mock.MockSite, nil)

	uc := NewUsecase(mockRepo)
	site, err := uc.GetSiteGroupByID(mock.MockSite.ID.Hex())

	assert.NoError(t, err)
	assert.Equal(t, mock.MockSite.ID.Hex(), site.ID.Hex())
	assert.Equal(t, mock.MockSite.Name, site.Name)
}

func TestUsecase_DeleteSiteGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	mockRepo.EXPECT().DeleteSiteGroup(mock.MockSite.ID.Hex()).Return(nil)

	uc := NewUsecase(mockRepo)
	err := uc.DeleteSiteGroup(mock.MockSite.ID.Hex())

	assert.NoError(t, err)
}
