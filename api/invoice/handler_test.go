package invoice

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	mockInvoice "gitlab.odds.team/worklog/api.odds-worklog/api/invoice/mock"
)

func TestCreate(t *testing.T) {
	t.Run("when create invoice success, then return json models.Invoice with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mockInvoice.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		iMock := mockInvoice.Invoice
		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Create(&iMock).Return(&iMock, nil)

		h := &HttpHandler{uMock}
		h.Create(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, mockInvoice.InvoiceJson, rec.Body.String())
	})

	t.Run("when bind invoice error, then return json models.HTTPError with status code 422", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		uMock := mockInvoice.NewMockUsecase(ctrl)
		h := &HttpHandler{uMock}
		h.Create(c)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when create invoice error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		iMock := mockInvoice.Invoice
		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Create(&iMock).Return(&iMock, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mockInvoice.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := &HttpHandler{uMock}
		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
