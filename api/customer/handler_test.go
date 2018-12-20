package customer

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mock "gitlab.odds.team/worklog/api.odds-worklog/api/customer/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

// func TestCreateCustomer(t *testing.T) {
// 	t.Run("when create Customer success it should return status OK", func(t *testing.T) {
// 		ctrl := gomock.NewController(t)
// 		defer ctrl.Finish()

// 		mockUsecase := mock.NewMockUsecase(ctrl)
// 		mockUsecase.EXPECT().CreateCustomer(&mock.MockCustomer).Return(&mock.MockCustomer, nil)

// 		e := echo.New()
// 		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(mock.UserJson))
// 		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 		rec := httptest.NewRecorder()
// 		c := e.NewContext(req, rec)

// 		claims := &models.JwtCustomClaims{
// 			&mock.MockAdmin,
// 			jwt.StandardClaims{
// 				ExpiresAt: time.Now().Add(time.Hour * 1).Unix(),
// 			},
// 		}
// 		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 		c.Set("user", token)
// 		handler := &HttpHandler{mockUsecase}
// 		handler.CreateCustomer(c)

// 		assert.Equal(t, http.StatusCreated, rec.Code)
// 	})
// }

func TestGetCustomers(t *testing.T) {
	t.Run("when get user success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := mock.NewMockUsecase(ctrl)
		mockListCustomer := make([]*models.Customer, 0)
		mockListCustomer = append(mockListCustomer, &mock.MockCustomer)
		mockUsecase.EXPECT().GetCustomers().Return(mockListCustomer, nil)

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
		handler.GetCustomers(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
