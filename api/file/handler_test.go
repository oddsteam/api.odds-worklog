package file

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	fileMock "gitlab.odds.team/worklog/api.odds-worklog/api/file/mock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func Test_getTranscriptFilename(t *testing.T) {
	u := userMock.User

	filename := getTranscriptFilename(&u)
	assert.NotEmpty(t, filename)

	path := "files/transcripts"
	filenameExp := fmt.Sprintf("%s/%s_%s_", path, strings.ToLower(u.FirstName), strings.ToLower(u.LastName))
	assert.Contains(t, filename, filenameExp)
	assert.Contains(t, filename, ".pdf")
	assert.Equal(t, len(filenameExp)+16, len(filename))
}

func TestDownloadTranscript(t *testing.T) {
	t.Run("download transcript file success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.Admin
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetPathTranscript(u.ID.Hex()).Return("test.pdf", nil)
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(u.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.DownloadTranscript(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when param id is empty then return status code 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase := fileMock.NewMockUsecase(ctrl)
		handler := &HttpHandler{mockUsecase}
		handler.DownloadTranscript(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"code":400,"message":"Invalid path"}`, rec.Body.String())
	})

	t.Run("when not owner transcript file and not admin then return status code 401", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.User
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		handler := &HttpHandler{mockUsecase}
		handler.DownloadTranscript(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, `{"code":401,"message":"Permission denied."}`, rec.Body.String())
	})

	t.Run("when transcript file error then return status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.User
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetPathTranscript(u.ID.Hex()).Return("", errors.New(""))
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(u.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.DownloadTranscript(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDownloadImageProfile(t *testing.T) {
	t.Run("download ImageProfile file success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.Admin
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetPathImageProfile(u.ID.Hex()).Return("test.pdf", nil)
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(u.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.DownloadImageProfile(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when param id is empty then return status code 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUsecase := fileMock.NewMockUsecase(ctrl)
		handler := &HttpHandler{mockUsecase}
		handler.DownloadImageProfile(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `{"code":400,"message":"Invalid path"}`, rec.Body.String())
	})

	t.Run("when not owner ImageProfile file and not admin then return status code 401", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.User
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		handler := &HttpHandler{mockUsecase}
		handler.DownloadImageProfile(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, `{"code":401,"message":"Permission denied."}`, rec.Body.String())
	})

	t.Run("when ImageProfile file error then return status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.User
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetPathImageProfile(u.ID.Hex()).Return("", errors.New(""))
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(u.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.DownloadImageProfile(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestRemoveTranscript(t *testing.T) {
	t.Run("when remove transcript success status code should be ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.User
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetPathTranscript(u.ID.Hex()).Return("test.pdf", nil)
		mockUsecase.EXPECT().RemoveTranscript("test.pdf").Return(nil)
		mockUsecase.EXPECT().UpdateUser(u.ID.Hex(), "").Return(nil)
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(u.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.RemoveTranscript(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
	t.Run("when user no have transcript status code should be internal server error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.User
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetPathTranscript(u.ID.Hex()).Return("test.pdf", errors.New(""))
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(u.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.RemoveTranscript(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Equal(t, `{"code":500,"message":"No transcript file."}`, rec.Body.String())
	})
	t.Run("when remove transcript is not success status code should be internal server error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.User
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetPathTranscript(u.ID.Hex()).Return("test.pdf", nil)
		mockUsecase.EXPECT().RemoveTranscript("test.pdf").Return(errors.New(""))
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues(u.ID.Hex())

		handler := &HttpHandler{mockUsecase}
		handler.RemoveTranscript(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when user send request to method but no have token status code should be unauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		u := userMock.User
		mockUsecase := fileMock.NewMockUsecase(ctrl)
		claims := &models.JwtCustomClaims{
			&u,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues("1234")

		handler := &HttpHandler{mockUsecase}
		handler.RemoveTranscript(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Equal(t, `{"code":401,"message":"Permission denied."}`, rec.Body.String())
	})
}
