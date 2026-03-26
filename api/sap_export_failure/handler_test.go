package sap_export_failure

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	sap_export_failure_mock "gitlab.odds.team/worklog/api.odds-worklog/api/sap_export_failure/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
)

func TestHttpHandler_List(t *testing.T) {
	t.Run("when user is not admin it should return StatusForbidden and not call usecase", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := sap_export_failure_mock.NewMockForViewingSAPExportFailures(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)

		h := &HttpHandler{uc: mockUC}
		err := h.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("when user is admin it should return OK and list from usecase", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := sap_export_failure_mock.NewMockForViewingSAPExportFailures(ctrl)
		mockUC.EXPECT().List(0).Return([]*models.SAPExportFailureLog{}, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		h := &HttpHandler{uc: mockUC}
		err := h.List(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
