package customer

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	cusMock "gitlab.odds.team/worklog/api.odds-worklog/api/customer/mock"
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
	t.Run("when create customer success, then return json models.Customer with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(cusMock.CustomerJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		cMock := cusMock.Customer
		usecaseMock := cusMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Create(&cMock).Return(&cMock, nil)

		h := &HttpHandler{usecaseMock}
		h.Create(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, cusMock.CustomerJson, rec.Body.String())
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

		usecaseMock := cusMock.NewMockUsecase(ctrl)
		h := &HttpHandler{usecaseMock}
		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when create customer error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cMock := cusMock.Customer
		usecaseMock := cusMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Create(&cMock).Return(&cMock, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(cusMock.CustomerJson))
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

		usecaseMock := cusMock.NewMockUsecase(ctrl)

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
	t.Run("when get customer list success, then return array json models.Customer with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		usecaseMock := cusMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Get().Return(cusMock.Customers, nil)

		h := &HttpHandler{usecaseMock}
		h.Get(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, cusMock.CustomersJson, rec.Body.String())
	})

	t.Run("when get customer list error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := cusMock.NewMockUsecase(ctrl)
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

		usecaseMock := cusMock.NewMockUsecase(ctrl)

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

func TestGetByID(t *testing.T) {
	t.Run("when get customer by id success, then return json models.Customer with status code 200", func(t *testing.T) {
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

		usecaseMock := cusMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().GetByID("1234").Return(&cusMock.Customer, nil)

		h := &HttpHandler{usecaseMock}
		h.GetByID(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, cusMock.CustomerJson, rec.Body.String())
	})

	t.Run("when get customer by id error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := cusMock.NewMockUsecase(ctrl)
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

		usecaseMock := cusMock.NewMockUsecase(ctrl)

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

func TestUpdate(t *testing.T) {
	t.Run("when update customer success, then return json models.Customer with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(cusMock.CustomerJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc3c001")

		cMock := cusMock.Customer
		usecaseMock := cusMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Update(&cMock).Return(&cMock, nil)

		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, cusMock.CustomerJson, rec.Body.String())
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
		c.SetParamValues("5bbcf2f90fd2df527bc3c001")

		usecaseMock := cusMock.NewMockUsecase(ctrl)
		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when update customer error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cMock := cusMock.Customer
		usecaseMock := cusMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Update(&cMock).Return(&cMock, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(cusMock.CustomerJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc3c001")

		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := cusMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc3c001")

		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestDelete(t *testing.T) {
	t.Run("when delete customer success, then return json models.Response with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		usecaseMock := cusMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Delete("1234").Return(nil)

		h := &HttpHandler{usecaseMock}
		h.Delete(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"message":"Delete customer success."}`, rec.Body.String())
	})

	t.Run("when delete customer error, then return json models.HTTPError status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := cusMock.NewMockUsecase(ctrl)
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

		usecaseMock := cusMock.NewMockUsecase(ctrl)

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
