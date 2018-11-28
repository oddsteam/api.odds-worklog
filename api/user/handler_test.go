package user

import (
	"encoding/json"
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
	// "gitlab.odds.team/worklog/api.odds-worklog/user"
)

func TestCreateUser(t *testing.T) {
	t.Run("when create user success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := new(models.User)
		mockUser.ID = bson.ObjectIdHex("5bbcf2f90fd2df527bc39539")
		mockUser.FullNameEn = "นายทดสอบชอบลงทุน"
		mockUser.Email = "test@abc.com"
		mockUser.BankAccountName = "ทดสอบชอบลงทุน"
		mockUser.BankAccountNumber = "123123123123"
		mockUser.ThaiCitizenID = "1234567890123"
		mockUser.CorporateFlag = "Y"

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().CreateUser(mockUser).Return(&mock.MockUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		handler := &HttpHandler{mockUsecase}
		handler.CreateUser(c)

		assert.Equal(t, http.StatusCreated, rec.Code)
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

}

func TestGetUserByID(t *testing.T) {
	t.Run("when get user by success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
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
}

func TestUpdateUser(t *testing.T) {
	t.Run("when update user success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := new(models.User)
		mockUser.ID = bson.ObjectIdHex("5bbcf2f90fd2df527bc39539")
		mockUser.FullNameEn = "นายทดสอบชอบลงทุน"
		mockUser.Email = "test@abc.com"
		mockUser.BankAccountName = "ทดสอบชอบลงทุน"
		mockUser.BankAccountNumber = "123123123123"
		mockUser.ThaiCitizenID = "1234567890123"
		mockUser.CorporateFlag = "Y"

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
		mockUsecase.EXPECT().UpdateUser(mockUser).Return(&mock.MockUser, nil)

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
}

func TestDeleteUser(t *testing.T) {
	t.Run("when delete user by success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
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
}

func TestUpdatePartialUser(t *testing.T) {
	t.Run("when update user by method patch success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUser := new(models.User)
		mockUser.ID = bson.ObjectIdHex("5bbcf2f90fd2df527bc39539")
		mockUser.FullNameEn = "ODDS junk"
		mockUser.Email = "xx@c.com"
		mockUser.BankAccountName = "ทดสอบชอบลงทุน"
		mockUser.BankAccountNumber = "123123123123"
		mockUser.ThaiCitizenID = "1234567890123"
		mockUser.CorporateFlag = "Y"

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &mock.MockUser)
		mockUsecase.EXPECT().GetUserByID(mock.MockUser.ID.Hex()).Return(mockUser, nil)
		mockUsecase.EXPECT().UpdateUser(mockUser).Return(mockUser, nil)
		mockIoReader := `{"fullnameEh" : "ODDS junk","email" : "xx@c.com"}`

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

		mockUser := new(models.User)
		mockUser.ID = bson.ObjectIdHex("5bbcf2f90fd2df527bc39539")
		mockUser.FullNameEn = "นายทดสอบชอบลงทุน"
		mockUser.Email = "test@abc.com"
		mockUser.BankAccountName = "ทดสอบชอบลงทุน"
		mockUser.BankAccountNumber = "123123123123"
		mockUser.ThaiCitizenID = "1234567890123"
		mockUser.CorporateFlag = "Y"

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
