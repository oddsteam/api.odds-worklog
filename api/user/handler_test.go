package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"

	mock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestCreateUser(t *testing.T) {
	t.Run("when create user success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().CreateUser(&mock.MockUser).Return(&mock.MockUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.CreateUser(c)

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
		handler.CreateUser(c)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when usecase createUser is have error it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().CreateUser(&mock.MockUser).Return(&mock.MockUser, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.CreateUser(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetUser(t *testing.T) {
	t.Run("when get user success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
		mockUsecase.EXPECT().GetUser().Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &models.JwtCustomClaims{
			&mock.MockAdmin,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		c.Set("user", token)
		handler := &HttpHandler{mockUsecase}
		handler.GetUser(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when current user is not admin it should return StatusUnauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &models.JwtCustomClaims{
			&mock.MockUser,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		c.Set("user", token)
		handler := &HttpHandler{mockUsecase}
		handler.GetUser(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when getUSer in usecase is have error  it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
		mockUsecase.EXPECT().GetUser().Return(mockListUser, errors.New(""))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		claims := &models.JwtCustomClaims{
			&mock.MockAdmin,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		c.Set("user", token)
		handler := &HttpHandler{mockUsecase}
		handler.GetUser(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}

func TestGetUserByID(t *testing.T) {
	t.Run("when get user by success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetUserByID(mock.MockUser.ID.Hex()).Return(&mock.MockUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		handler := &HttpHandler{mockUsecase}
		handler.GetUserByID(c)

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
		handler.GetUserByID(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when getUSer in usecase is have error  it should return StatusNoContent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetUserByID(mock.MockUser.ID.Hex()).Return(&mock.MockUser, errors.New(""))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.GetUserByID(c)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

}

func TestGetUserBySiteId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := mock.NewMockUsecase(ctrl)
	mockListUser := make([]*models.User, 0)
	mockListUser = append(mockListUser, &mock.MockUser)
	mockUsecase.EXPECT().GetUserBySiteID("12345").Return(mockListUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	claims := &models.JwtCustomClaims{
		&mock.MockAdmin,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	c.Set("user", token)
	c.SetParamNames("id")
	c.SetParamValues("12345")

	handler := &HttpHandler{mockUsecase}
	handler.GetUserBySiteID(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdateUser(t *testing.T) {
	t.Run("when update user success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
		mockUsecase.EXPECT().UpdateUser(&mock.MockUser, nil).Return(&mock.MockUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bc89e26f37e2f0df54e6fef")
		handler := &HttpHandler{mockUsecase}
		handler.UpdateUser(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when request is no have id it should return StatusBadRequest", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.UpdateUser(c)

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
		c.SetParamValues("5bc89e26f37e2f0df54e6fef")
		handler := &HttpHandler{mockUsecase}
		handler.UpdateUser(c)

		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	})

	t.Run("when UpdateUser in usecase have error it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
		mockUsecase.EXPECT().UpdateUser(&mock.MockUser, nil).Return(&mock.MockUser, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(mock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bc89e26f37e2f0df54e6fef")
		handler := &HttpHandler{mockUsecase}
		handler.UpdateUser(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("when delete user by success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().DeleteUser(mock.MockUser.ID.Hex()).Return(nil)

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		claims := &models.JwtCustomClaims{
			&mock.MockAdmin,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		handler := &HttpHandler{mockUsecase}
		handler.DeleteUser(c)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("when current user is not admin it should return StatusUnauthorized", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		claims := &models.JwtCustomClaims{
			&mock.MockUser,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		c.Set("user", token)
		handler := &HttpHandler{mockUsecase}
		handler.DeleteUser(c)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("when DeleteUser in usecase is have error it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().DeleteUser(mock.MockUser.ID.Hex()).Return(errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		claims := &models.JwtCustomClaims{
			&mock.MockAdmin,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		c.Set("user", token)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		handler := &HttpHandler{mockUsecase}
		handler.DeleteUser(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdatePartialUser(t *testing.T) {
	t.Run("when update user by method patch success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := new(models.User)
		mockUser.ID = bson.ObjectIdHex("5bbcf2f90fd2df527bc39539")
		mockUser.FirstName = "ODDS"
		mockUser.LastName = "junk"
		mockUser.Email = "xx@c.com"
		mockUser.BankAccountName = "ทดสอบชอบลงทุน"
		mockUser.BankAccountNumber = "123123123123"
		mockUser.ThaiCitizenID = "1234567890123"

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
		mockUsecase.EXPECT().GetUserByID(mock.MockUser.ID.Hex()).Return(mockUser, nil)
		mockUsecase.EXPECT().UpdateUser(mockUser, nil).Return(mockUser, nil)
		mockIoReader := `{"firstName" : "ODDS","lastName" : "junk","email" : "xx@c.com"}`

		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(mockIoReader))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.UpdatePartialUser(c)

		userByte, _ := json.Marshal(mockUser)
		UserJson := string(userByte)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, UserJson, rec.Body.String())
	})

	t.Run("should return InternalError if no have requestBody", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
		mockUsecase.EXPECT().GetUserByID(mock.MockUser.ID.Hex()).Return(&mock.MockUser, nil)

		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.UpdatePartialUser(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestIsUserAdmin(t *testing.T) {
	t.Run("it should return true if user is admin", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		claims := &models.JwtCustomClaims{
			&mock.MockAdmin,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		c.Set("user", token)
		isAdmin, _ := IsUserAdmin(c)

		assert.Equal(t, true, isAdmin)
	})

	t.Run("it should return false with message if user is not admin", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodPatch, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		claims := &models.JwtCustomClaims{
			&mock.MockUser,
			jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		c.Set("user", token)
		isAdmin, mess := IsUserAdmin(c)

		assert.Equal(t, false, isAdmin)
		assert.Equal(t, "ไม่มีสิทธิในการใช้งาน", mess)
	})
}
