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
	invoiceMock "gitlab.odds.team/worklog/api.odds-worklog/api/invoice/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
)

func TestGetUserFromToken(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.POST, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", userMock.TokenAdmin)

	user := getUserFromToken(c)
	b, _ := json.Marshal(user)
	actual := string(b)
	assert.Equal(t, userMock.AdminJson, actual)
}

func TestCreate(t *testing.T) {
	t.Run("when create invoice success, then return json models.Invoice with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(invoiceMock.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		iMock := invoiceMock.Invoice
		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Create(&iMock).Return(&iMock, nil)

		h := &HttpHandler{usecaseMock}
		h.Create(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, invoiceMock.InvoiceJson, rec.Body.String())
	})

	t.Run("when request is invalid, then return json models.HTTPError with status code 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		h := &HttpHandler{usecaseMock}
		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when create invoice error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		iMock := invoiceMock.Invoice
		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Create(&iMock).Return(&iMock, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(invoiceMock.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		h := &HttpHandler{usecaseMock}
		h.Create(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
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
		c.Set("user", userMock.TokenAdmin)

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Get().Return(invoiceMock.Invoices, nil)

		h := &HttpHandler{usecaseMock}
		h.Get(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, invoiceMock.InvoicesJson, rec.Body.String())
	})

	t.Run("when get invoice list error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Get().Return(nil, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		h := &HttpHandler{usecaseMock}
		h.Get(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)

		h := &HttpHandler{usecaseMock}
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
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().GetByPO("1234").Return(invoiceMock.Invoices, nil)

		h := &HttpHandler{usecaseMock}
		h.GetByPO(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, invoiceMock.InvoicesJson, rec.Body.String())
	})

	t.Run("when get invoice list by PO error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().GetByPO("1234").Return(nil, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
		h.GetByPO(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
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
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().GetByID("1234").Return(&invoiceMock.Invoice, nil)

		h := &HttpHandler{usecaseMock}
		h.GetByID(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, invoiceMock.InvoiceJson, rec.Body.String())
	})

	t.Run("when get invoice by id error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().GetByID("1234").Return(nil, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
		h.GetByID(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
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
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		no := "2018_001"
		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().NextNo("1234").Return(no, nil)

		h := &HttpHandler{usecaseMock}
		h.NextNo(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"invoiceNo":"2018_001"}`, rec.Body.String())
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
		h.NextNo(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("when nextNo is error, then reuturn json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().NextNo("1234").Return("", errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
		h.NextNo(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("when update invoice success, then return json models.Invoice with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(invoiceMock.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc30000")

		iMock := invoiceMock.Invoice
		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Update(&iMock).Return(&iMock, nil)

		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, invoiceMock.InvoiceJson, rec.Body.String())
	})

	t.Run("when request is invalid, then return json models.HTTPError with status code 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when update invoice error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		iMock := invoiceMock.Invoice
		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Update(&iMock).Return(&iMock, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(invoiceMock.InvoiceJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc30000")

		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc30000")

		h := &HttpHandler{usecaseMock}
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
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Delete("1234").Return(nil)

		h := &HttpHandler{usecaseMock}
		h.Delete(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when delete invoice error, then return json models.HTTPError status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Delete("1234").Return(errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
		h.Delete(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := invoiceMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
		h.Delete(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
