package site

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.odds.team/worklog/api.odds-worklog/business/models"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	mock "gitlab.odds.team/worklog/api.odds-worklog/api/site/mock"
)

func TestCreateSiteGroup(t *testing.T) {
	t.Run("when create site success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		site := models.Site{Name: "ktb"}
		siteJson := `{"name": "ktb"}`
		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().CreateSiteGroup(&site).Return(&mock.MockSite, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(siteJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.CreateSiteGroup(c)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("when content type is not valid it should return StatusUnprocessableEntity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader("string"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.CreateSiteGroup(c)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when usecase createSiteGroup is have error it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		site := models.Site{Name: "ktb"}
		siteJson := `{"name": "ktb"}`
		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().CreateSiteGroup(&site).Return(&mock.MockSite, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(siteJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.CreateSiteGroup(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("when update site success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().UpdateSiteGroup(&mock.MockSite).Return(&mock.MockSite, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mock.SiteJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.UpdateSiteGroup(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when request is no have id it should return StatusBadRequest", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mock.SiteJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.UpdateSiteGroup(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when content type is not valid it should return StatusUnprocessableEntity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader("string"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		handler := &HttpHandler{mockUsecase}
		handler.UpdateSiteGroup(c)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when UpdateSiteGroup in usecase have error it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().UpdateSiteGroup(&mock.MockSite).Return(&mock.MockSite, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mock.SiteJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.UpdateSiteGroup(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetSiteGroup(t *testing.T) {
	t.Run("when get site success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetSiteGroup().Return(mock.MockSites, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetSiteGroup(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when get site list in usecase is have error  it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetSiteGroup().Return(mock.MockSites, errors.New(""))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetSiteGroup(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}

func TestGetUserByID(t *testing.T) {
	t.Run("when get site by id success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetSiteGroupByID(mock.MockSite.ID.Hex()).Return(&mock.MockSite, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		handler := &HttpHandler{mockUsecase}
		handler.GetSiteGroupByID(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when request is no have id it should return StatusBadRequest", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.GetSiteGroupByID(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when get site by id in usecase is have error, it should return StatusNoContent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetSiteGroupByID(mock.MockSite.ID.Hex()).Return(&mock.MockSite, errors.New(""))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.GetSiteGroupByID(c)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

}

func TestDeleteUser(t *testing.T) {
	t.Run("when delete site by id success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().DeleteSiteGroup(mock.MockSite.ID.Hex()).Return(nil)

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		handler := &HttpHandler{mockUsecase}
		handler.DeleteSiteGroup(c)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("when DeleteSite in usecase is have error it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().DeleteSiteGroup(mock.MockSite.ID.Hex()).Return(errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		handler := &HttpHandler{mockUsecase}
		handler.DeleteSiteGroup(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
