package invoice

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	invoiceMock "gitlab.odds.team/worklog/api.odds-worklog/api/invoice/mock"
)

func TestUsecase_Create(t *testing.T) {
	t.Run("when create invoice success, then return (*models.Invoice, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		inv := invoiceMock.Invoice
		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Create(&inv).Return(&inv, nil)

		u := NewUsecase(mockRepo)
		invoice, err := u.Create(&inv)

		assert.NoError(t, err)
		assert.Equal(t, inv.PoID, invoice.PoID)
		assert.Equal(t, inv.InvoiceNo, invoice.InvoiceNo)
		assert.Equal(t, inv.Amount, invoice.Amount)
	})

	t.Run("when create invoice error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		inv := invoiceMock.Invoice
		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Create(&inv).Return(nil, errors.New("error"))

		u := NewUsecase(mockRepo)
		invoice, err := u.Create(&inv)

		assert.Error(t, err)
		assert.Nil(t, invoice)
	})
}

func TestUsecase_Get(t *testing.T) {
	t.Run("when get invoice list success, then return ([]*models.Invoice, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Get().Return(invoiceMock.Invoices, nil)

		u := NewUsecase(mockRepo)
		invoices, err := u.Get()

		b, _ := json.Marshal(invoices)
		j := string(b)
		assert.NoError(t, err)
		assert.Equal(t, invoiceMock.InvoicesJson, j)
	})

	t.Run("when get invoice list error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Get().Return(nil, errors.New(""))

		u := NewUsecase(mockRepo)
		invoices, err := u.Get()

		assert.Error(t, err)
		assert.Nil(t, invoices)
	})
}
