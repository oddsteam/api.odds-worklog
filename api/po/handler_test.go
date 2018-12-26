package po

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
	poMock "gitlab.odds.team/worklog/api.odds-worklog/api/po/mock"
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
	t.Run("when create po success, then return json models.Po with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(poMock.PoJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		pMock := poMock.Po
		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Create(&pMock).Return(&pMock, nil)

		h := &HttpHandler{usecaseMock}
		h.Create(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, poMock.PoJson, rec.Body.String())
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

		usecaseMock := poMock.NewMockUsecase(ctrl)
		h := &HttpHandler{usecaseMock}
		h.Create(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when create po error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pMock := poMock.Po
		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Create(&pMock).Return(&pMock, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(poMock.PoJson))
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

		usecaseMock := poMock.NewMockUsecase(ctrl)

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
	t.Run("when get po list success, then return array json models.Po with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Get().Return(poMock.Poes, nil)

		h := &HttpHandler{usecaseMock}
		h.Get(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, poMock.PoesJson, rec.Body.String())
	})

	t.Run("when get po list error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := poMock.NewMockUsecase(ctrl)
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

		usecaseMock := poMock.NewMockUsecase(ctrl)

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

func TestGetByCusID(t *testing.T) {
	t.Run("when get po list by customer id success, then return array json models.Po with status code 200", func(t *testing.T) {
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

		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().GetByCusID("1234").Return(poMock.Poes, nil)

		h := &HttpHandler{usecaseMock}
		h.GetByCusID(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, poMock.PoesJson, rec.Body.String())
	})

	t.Run("when get po list by customer id error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().GetByCusID("1234").Return(nil, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
		h.GetByCusID(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := poMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		h := &HttpHandler{usecaseMock}
		h.GetByCusID(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("when get po by id success, then return json models.Po with status code 200", func(t *testing.T) {
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

		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().GetByID("1234").Return(&poMock.Po, nil)

		h := &HttpHandler{usecaseMock}
		h.GetByID(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, poMock.PoJson, rec.Body.String())
	})

	t.Run("when get po by id error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := poMock.NewMockUsecase(ctrl)
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

		usecaseMock := poMock.NewMockUsecase(ctrl)

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
	t.Run("when update po success, then return json models.Po with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(poMock.PoJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5c1f855d59fc7d06988c6e01")

		pMock := poMock.Po
		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Update(&pMock).Return(&pMock, nil)

		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, poMock.PoJson, rec.Body.String())
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
		c.SetParamValues("5c1f855d59fc7d06988c6e01")

		usecaseMock := poMock.NewMockUsecase(ctrl)
		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when update po error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		pMock := poMock.Po
		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Update(&pMock).Return(&pMock, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(poMock.PoJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5c1f855d59fc7d06988c6e01")

		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := poMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5c1f855d59fc7d06988c6e01")

		h := &HttpHandler{usecaseMock}
		h.Update(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestDelete(t *testing.T) {
	t.Run("when delete po success, then return json models.Response with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5c1f855d59fc7d06988c6e01")

		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Delete("5c1f855d59fc7d06988c6e01").Return(nil)

		h := &HttpHandler{usecaseMock}
		h.Delete(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"message":"Delete product owner success."}`, rec.Body.String())
	})

	t.Run("when delete po error, then return json models.HTTPError status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := poMock.NewMockUsecase(ctrl)
		usecaseMock.EXPECT().Delete("5c1f855d59fc7d06988c6e01").Return(errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5c1f855d59fc7d06988c6e01")

		h := &HttpHandler{usecaseMock}
		h.Delete(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		usecaseMock := poMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5c1f855d59fc7d06988c6e01")

		h := &HttpHandler{usecaseMock}
		h.Delete(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}
