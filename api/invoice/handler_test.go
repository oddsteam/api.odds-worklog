package invoice

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	mockInvoice "gitlab.odds.team/worklog/api.odds-worklog/api/invoice/mock"
	mockUser "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
)

func TestGetUserFromToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", mockUser.TokenAdmin)

	user := getUserFromToken(c)
	b, _ := json.Marshal(user)
	actual := string(b)
	assert.Equal(t, mockUser.MockAdminJson, actual)
}

func TestCreate(t *testing.T) {
	t.Run("when create invoice success, then return json models.Invoice with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mockInvoice.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)

		iMock := mockInvoice.Invoice
		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Create(&iMock).Return(&iMock, nil)

		h := &HttpHandler{uMock}
		h.Create(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, mockInvoice.InvoiceJson, rec.Body.String())
	})

	t.Run("when request is invalid, then return json models.HTTPError with status code 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)

		uMock := mockInvoice.NewMockUsecase(ctrl)
		h := &HttpHandler{uMock}
		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
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
		c.Set("user", mockUser.TokenAdmin)

		h := &HttpHandler{uMock}
		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.Create(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestGet(t *testing.T) {
	t.Run("when get invoice list success, then return array json models.Invoice with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Get().Return(mockInvoice.Invoices, nil)

		h := &HttpHandler{uMock}
		h.Get(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, mockInvoice.InvoicesJson, rec.Body.String())
	})

	t.Run("when get invoice list error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Get().Return(nil, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)

		h := &HttpHandler{uMock}
		h.Get(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenUser)

		h := &HttpHandler{uMock}
		h.Get(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestGetByPO(t *testing.T) {
	t.Run("when get invoice list by PO success, then return array json models.Invoice with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().GetByPO("1234").Return(mockInvoice.Invoices, nil)

		h := &HttpHandler{uMock}
		h.GetByPO(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, mockInvoice.InvoicesJson, rec.Body.String())
	})

	t.Run("when get invoice list by PO error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().GetByPO("1234").Return(nil, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.GetByPO(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.GetByPO(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("when get invoice by id success, then return json models.Invoice with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().GetByID("1234").Return(&mockInvoice.Invoice, nil)

		h := &HttpHandler{uMock}
		h.GetByID(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, mockInvoice.InvoiceJson, rec.Body.String())
	})

	t.Run("when get invoice by id error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().GetByID("1234").Return(nil, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.GetByID(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.GetByID(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestNextNo(t *testing.T) {
	t.Run("when request is valid, then reuturn json models.InvoiceNoRes with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		no := "2018_001"
		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().NextNo("1234").Return(no, nil)

		h := &HttpHandler{uMock}
		h.NextNo(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"invoiceNo":"2018_001"}`, rec.Body.String())
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.NextNo(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("when nextNo is error, then reuturn json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().NextNo("1234").Return("", errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.NextNo(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("when update invoice success, then return json models.Invoice with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mockInvoice.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		iMock := mockInvoice.Invoice
		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Update(&iMock).Return(&iMock, nil)

		h := &HttpHandler{uMock}
		h.Update(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, mockInvoice.InvoiceJson, rec.Body.String())
	})

	t.Run("when request is invalid, then return json models.HTTPError with status code 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		uMock := mockInvoice.NewMockUsecase(ctrl)
		h := &HttpHandler{uMock}
		h.Update(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when update invoice error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		iMock := mockInvoice.Invoice
		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Update(&iMock).Return(&iMock, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mockInvoice.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.Update(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.Update(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestDelete(t *testing.T) {
	t.Run("when delete invoice success, then return json models.Response with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Delete("1234").Return(nil)

		h := &HttpHandler{uMock}
		h.Delete(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when delete invoice error, then return json models.HTTPError status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)
		uMock.EXPECT().Delete("1234").Return(errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.Delete(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		uMock := mockInvoice.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", mockUser.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{uMock}
		h.Delete(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
