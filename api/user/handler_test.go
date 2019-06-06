package user

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"

	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
)

func TestCreate(t *testing.T) {
	t.Run("when create user success, return json models.User with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Create(&userMock.User).Return(&userMock.User, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(userMock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.Create(c)

		assert.Equal(t, http.StatusCreated, rec.Code)
	})

	t.Run("when content type is not valid it should return StatusBadRequest", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader("string"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.Create(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when usecase createUser is have error it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Create(&userMock.User).Return(&userMock.User, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(userMock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.Create(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

	t.Run("when request isn't admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(userMock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)

		handler := &HttpHandler{mockUsecase}
		handler.Create(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})
}

func TestGet(t *testing.T) {
	t.Run("when get user success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &userMock.User)
		mockUsecase.EXPECT().Get().Return(mockListUser, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.Get(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when current user is not admin it should return StatusForbidden", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &userMock.User)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)

		handler := &HttpHandler{mockUsecase}
		handler.Get(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("when getUser in usecase is have error  it should return StatusInternalServerError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockListUser := make([]*models.User, 0)
		mockListUser = append(mockListUser, &userMock.User)
		mockUsecase.EXPECT().Get().Return(mockListUser, errors.New(""))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)

		handler := &HttpHandler{mockUsecase}
		handler.Get(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})

}

func TestGetByEmail(t *testing.T) {
	t.Run("when get user by email success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetByEmail(userMock.User.Email).Return(&userMock.User, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("email")
		c.SetParamValues("test@abc.com")

		handler := &HttpHandler{mockUsecase}
		handler.GetByEmail(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when request is not admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("email")
		c.SetParamValues("test@abc.com")
		handler := &HttpHandler{mockUsecase}
		handler.GetByEmail(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("when get user by email error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetByEmail(userMock.User.Email).Return(nil, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("email")
		c.SetParamValues("test@abc.com")

		handler := &HttpHandler{mockUsecase}
		handler.GetByEmail(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestGetByID(t *testing.T) {
	t.Run("when get user by id success it should return status OK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetByID(userMock.User.ID.Hex()).Return(&userMock.User, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")
		handler := &HttpHandler{mockUsecase}
		handler.GetByID(c)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("when getUser in usecase is have error  it should return StatusNoContent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().GetByID(userMock.User.ID.Hex()).Return(&userMock.User, errors.New(""))
		e := echo.New()
		req := httptest.NewRequest(echo.GET, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.GetByID(c)

		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

}

func TestGetBySiteId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUsecase := userMock.NewMockUsecase(ctrl)
	mockListUser := make([]*models.User, 0)
	mockListUser = append(mockListUser, &userMock.User)
	mockUsecase.EXPECT().GetBySiteID("12345").Return(mockListUser, nil)

	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", userMock.TokenAdmin)
	c.SetParamNames("id")
	c.SetParamValues("12345")

	handler := &HttpHandler{mockUsecase}
	handler.GetBySiteID(c)

	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestUpdate(t *testing.T) {
	t.Run("when update user success, then return json models.User with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Update(&userMock.User, gomock.Any()).Return(&userMock.User, nil)

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(userMock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.Update(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, userMock.UserJson, rec.Body.String())
	})

	t.Run("when request is invalid, then return json models.HTTPError with status code 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bc89e26f37e2f0df54e6fef")

		handler := &HttpHandler{mockUsecase}
		handler.Update(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("when update user error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Update(&userMock.User, gomock.Any()).Return(&userMock.User, errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(userMock.UserJson))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.Update(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestUpdateStatusTavi(t *testing.T) {
	// t.Run("when update status tavi success, then return json models.Users with status code 200", func(t *testing.T) {
	// 	ctrl := gomock.NewController(t)
	// 	defer ctrl.Finish()

	// 	mockUsecase := userMock.NewMockUsecase(ctrl)
	// 	mockUsecase.EXPECT().UpdateStatusTavi(userMock.ListUser, gomock.Any()).Return(nil, nil)

	// 	e := echo.New()
	// 	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(userMock.UsersJson))
	// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 	rec := httptest.NewRecorder()
	// 	c := e.NewContext(req, rec)
	// 	c.Set("user", userMock.TokenUser)

	// 	handler := &HttpHandler{mockUsecase}
	// 	handler.UpdateStatusTavi(c)

	// 	assert.Equal(t, http.StatusOK, rec.Code)
	// 	assert.Equal(t, userMock.UsersJson, rec.Body.String())
	// })

	t.Run("when request is invalid, then return json models.HTTPError with status code 400", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.POST, "/", strings.NewReader(""))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("5bc89e26f37e2f0df54e6fef")

		handler := &HttpHandler{mockUsecase}
		handler.UpdateStatusTavi(c)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	// t.Run("when update status tavi error, then return json models.HTTPError with status code 500", func(t *testing.T) {
	// 	ctrl := gomock.NewController(t)
	// 	defer ctrl.Finish()

	// 	mockUsecase := userMock.NewMockUsecase(ctrl)
	// 	mockUsecase.EXPECT().UpdateStatusTavi(userMock.ListUser, gomock.Any()).Return(userMock.Users, errors.New(""))

	// 	e := echo.New()
	// 	req := httptest.NewRequest(echo.PUT, "/", strings.NewReader(userMock.UsersJson))
	// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	// 	rec := httptest.NewRecorder()
	// 	c := e.NewContext(req, rec)
	// 	c.Set("user", userMock.TokenUser)
	// 	c.SetParamNames("id")
	// 	c.SetParamValues("5bbcf2f90fd2df527bc39539")

	// 	handler := &HttpHandler{mockUsecase}
	// 	handler.UpdateStatusTavi(c)

	// 	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	// })
}

func TestDelete(t *testing.T) {
	t.Run("when delete user success, then return json models.Response with status code 200", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Delete(userMock.User.ID.Hex()).Return(nil)

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.Delete(c)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, `{"message":"Delete user success."}`, rec.Body.String())
	})

	t.Run("when request is not admin, then return json models.HTTPError with status code 403", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenUser)

		handler := &HttpHandler{mockUsecase}
		handler.Delete(c)

		assert.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("when delete error, then return json models.HTTPError with status code 500", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUsecase := userMock.NewMockUsecase(ctrl)
		mockUsecase.EXPECT().Delete(userMock.User.ID.Hex()).Return(errors.New(""))

		e := echo.New()
		req := httptest.NewRequest(echo.DELETE, "/", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.Set("user", userMock.TokenAdmin)
		c.SetParamNames("id")
		c.SetParamValues("5bbcf2f90fd2df527bc39539")

		handler := &HttpHandler{mockUsecase}
		handler.Delete(c)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}
