package usecases

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
	mock_usecases "gitlab.odds.team/worklog/api.odds-worklog/business/usecases/mock"
)

func TestNormalizeSapExportFailureLimit(t *testing.T) {
	assert.Equal(t, 100, normalizeSapExportFailureLimit(0))
	assert.Equal(t, 100, normalizeSapExportFailureLimit(-1))
	assert.Equal(t, 50, normalizeSapExportFailureLimit(50))
	assert.Equal(t, 200, normalizeSapExportFailureLimit(200))
	assert.Equal(t, 200, normalizeSapExportFailureLimit(9999))
}

func TestViewSAPExportFailuresUsecase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockList := mock_usecases.NewMockForListingSAPExportFailures(ctrl)
	mockList.EXPECT().List(100).Return([]*models.SAPExportFailureLog{{IncomeID: "x"}}, nil)

	uc := NewViewSAPExportFailuresUsecase(mockList)
	logs, err := uc.List(0)
	assert.NoError(t, err)
	assert.Len(t, logs, 1)
	assert.Equal(t, "x", logs[0].IncomeID)
}
