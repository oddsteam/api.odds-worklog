package customer

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	cusMock "gitlab.odds.team/worklog/api.odds-worklog/api/customer/mock"
)

func TestUsecase_Create(t *testing.T) {
	t.Run("when create customer success, then return (*models.Customer, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Create(&cus).Return(&cus, nil)

		u := NewUsecase(mockRepo)
		customer, err := u.Create(&cus)

		b, _ := json.Marshal(customer)
		js := string(b)
		assert.NoError(t, err)
		assert.Equal(t, cusMock.CustomerJson, js)
	})

	t.Run("when create customer error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Create(&cus).Return(nil, errors.New("error"))

		u := NewUsecase(mockRepo)
		customer, err := u.Create(&cus)

		assert.Error(t, err)
		assert.Nil(t, customer)
	})
}

func TestUsecase_Get(t *testing.T) {
	t.Run("when get customer list success, then return ([]*models.Customer, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Get().Return(cusMock.Customers, nil)

		u := NewUsecase(mockRepo)
		customers, err := u.Get()

		b, _ := json.Marshal(customers)
		js := string(b)
		assert.NoError(t, err)
		assert.Equal(t, cusMock.CustomersJson, js)
	})

	t.Run("when get customer list error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Get().Return(nil, errors.New("error"))

		u := NewUsecase(mockRepo)
		customers, err := u.Get()

		assert.Error(t, err)
		assert.Nil(t, customers)
	})
}

func TestUsecase_GetByID(t *testing.T) {
	t.Run("when get customer success, then return ([]*models.Customer, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(cus.ID.Hex()).Return(&cus, nil)

		u := NewUsecase(mockRepo)
		customer, err := u.GetByID(cus.ID.Hex())

		b, _ := json.Marshal(customer)
		js := string(b)
		assert.NoError(t, err)
		assert.Equal(t, cusMock.CustomerJson, js)
	})

	t.Run("when get customer list error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(cus.ID.Hex()).Return(nil, errors.New("error"))

		u := NewUsecase(mockRepo)
		customer, err := u.GetByID(cus.ID.Hex())

		assert.Error(t, err)
		assert.Nil(t, customer)
	})
}

func TestUsecase_Update(t *testing.T) {
	t.Run("when update customer success, then return (*models.Customer, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(cus.ID.Hex()).Return(&cus, nil)
		mockRepo.EXPECT().Update(&cus).Return(&cus, nil)

		u := NewUsecase(mockRepo)
		customer, err := u.Update(&cus)

		b, _ := json.Marshal(customer)
		js := string(b)
		assert.NoError(t, err)
		assert.Equal(t, cusMock.CustomerJson, js)
	})

	t.Run("when update customer error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(cus.ID.Hex()).Return(&cus, nil)
		mockRepo.EXPECT().Update(&cus).Return(nil, errors.New("error"))

		u := NewUsecase(mockRepo)
		customer, err := u.Update(&cus)

		assert.Error(t, err)
		assert.Nil(t, customer)
	})

	t.Run("when get customer error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID(cus.ID.Hex()).Return(nil, errors.New("error"))

		u := NewUsecase(mockRepo)
		customer, err := u.Update(&cus)

		assert.Error(t, err)
		assert.Nil(t, customer)
	})
}

func TestUsecase_Delete(t *testing.T) {
	t.Run("when delete customer success, then return error nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Delete(cus.ID.Hex()).Return(nil)

		u := NewUsecase(mockRepo)
		err := u.Delete(cus.ID.Hex())
		assert.NoError(t, err)
	})

	t.Run("when delete customer error, then return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cus := cusMock.Customer
		mockRepo := cusMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Delete(cus.ID.Hex()).Return(errors.New("error"))

		u := NewUsecase(mockRepo)
		err := u.Delete(cus.ID.Hex())

		assert.Error(t, err)
	})
}
