package backoffice

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	backofficeMock "gitlab.odds.team/worklog/api.odds-worklog/api/backoffice/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestHandleGetUserIncome(t *testing.T) {
	t.Run("get all userIncomes success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := backofficeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Get().Return(backofficeMock.MockUserIncomeList, nil)
		mockUsecase.EXPECT().GetKey().Return(&backofficeMock.MockBackOfficeKey, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/backoffice", strings.NewReader(backofficeMock.BackOfficeKeyReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetAllUserIncome(c)

		assert.Equal(t, http.StatusOK, rec.Code)

	})

	t.Run("get all userIncomes fail when invalid key", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := backofficeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetKey().Return(&backofficeMock.MockBackOfficeKey, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/backoffice", strings.NewReader(backofficeMock.InvalidBackOfficeKeyReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetAllUserIncome(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

	})

	t.Run("get all userIncomes fail when fail error when get all userIncome", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := backofficeMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Get().Return(nil, utils.ErrNotFound)
		mockUsecase.EXPECT().GetKey().Return(&backofficeMock.MockBackOfficeKey, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/backoffice", strings.NewReader(backofficeMock.BackOfficeKeyReqJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetAllUserIncome(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

	})

}
