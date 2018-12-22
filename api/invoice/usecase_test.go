package invoice

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	invoiceMock "gitlab.odds.team/worklog/api.odds-worklog/api/invoice/mock"
	poMock "gitlab.odds.team/worklog/api.odds-worklog/api/po/mock"
)

func TestUsecase_Create(t *testing.T) {
	t.Run("when create invoice success, then return (*models.Invoice, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		inv := invoiceMock.Invoice
		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Create(&inv).Return(&inv, nil)
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
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
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
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
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
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
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		invoices, err := u.Get()

		assert.Error(t, err)
		assert.Nil(t, invoices)
	})
}

func TestUsecase_GetByPO(t *testing.T) {
	t.Run("when get invoice list by PO success, then return ([]*models.Invoice, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByPO("1234").Return(invoiceMock.Invoices, nil)
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		invoices, err := u.GetByPO("1234")

		b, _ := json.Marshal(invoices)
		j := string(b)
		assert.NoError(t, err)
		assert.Equal(t, invoiceMock.InvoicesJson, j)
	})

	t.Run("when get invoice list by PO error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByPO("1234").Return(nil, errors.New(""))
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		invoices, err := u.GetByPO("1234")

		assert.Error(t, err)
		assert.Nil(t, invoices)
	})
}

func TestUsecase_GetByID(t *testing.T) {
	t.Run("when get invoice by id success, then return (*models.Invoice, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID("1234").Return(&invoiceMock.Invoice, nil)
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		invoice, err := u.GetByID("1234")

		b, _ := json.Marshal(invoice)
		j := string(b)
		assert.NoError(t, err)
		assert.Equal(t, invoiceMock.InvoiceJson, j)
	})

	t.Run("when get invoice by id error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().GetByID("1234").Return(nil, errors.New(""))
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		invoice, err := u.GetByID("1234")

		assert.Error(t, err)
		assert.Nil(t, invoice)
	})
}

func TestUsecase_NextNo(t *testing.T) {
	t.Run("when get last invoice success, then return (string, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mIvoice := invoiceMock.Invoice
		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Last(mIvoice.PoID).Return(&mIvoice, nil)
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		invoiceNo, err := u.NextNo(invoiceMock.Invoice.PoID)

		assert.NoError(t, err)
		assert.NotEmpty(t, invoiceNo)
	})

	t.Run("when get last invoice error, then return (string, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mIvoice := invoiceMock.Invoice
		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Last(mIvoice.PoID).Return(nil, errors.New(""))
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		invoiceNo, err := u.NextNo(mIvoice.PoID)

		assert.NoError(t, err)
		assert.NotEmpty(t, invoiceNo)
	})
}

func TestUsecase_NewNo(t *testing.T) {
	var lastNo, expected, actual string
	var err error

	t.Run("when last invoice no in current year, then return invoice no increase by 1", func(t *testing.T) {
		lastNo = "2018_001"
		expected = "2018_002"
		actual, err = newNo(lastNo)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)

		lastNo = "2018_099"
		expected = "2018_100"
		actual, err = newNo(lastNo)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)

		lastNo = "2018_101"
		expected = "2018_102"
		actual, err = newNo(lastNo)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})

	t.Run("when create new invoice no, then return invoice no is yyyy_001 (yyyy is year pattern)", func(t *testing.T) {
		lastNo = ""
		expected = fmt.Sprintf("%04d_001", time.Now().Year())
		actual, err = newNo(lastNo)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})

	t.Run("when last invoice no isn't in cerrent year, then return invoice no is yyyy_001 in current year (yyyy is year pattern)", func(t *testing.T) {
		lastNo = "2017_900"
		expected = fmt.Sprintf("%04d_001", time.Now().Year())
		actual, err = newNo(lastNo)
		assert.Equal(t, expected, actual)
		assert.NoError(t, err)
	})

	t.Run(`when new invoice no over limit 999, then return error "Over limit 999 invoices."`, func(t *testing.T) {
		lastNo = fmt.Sprintf("%04d_999", time.Now().Year())
		expected = ""
		actual, err = newNo(lastNo)
		assert.Equal(t, expected, actual)
		assert.EqualError(t, errors.New("Over limit 999 invoices."), err.Error())
	})
}

func TestUsecase_Delete(t *testing.T) {
	t.Run("when delete invoice success, then return nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Delete("1234").Return(nil)
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		err := u.Delete("1234")
		assert.NoError(t, err)
	})

	t.Run("when delete invoice error, then return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := invoiceMock.NewMockRepository(ctrl)
		mockRepo.EXPECT().Delete("1234").Return(errors.New(""))
		mockPoRepo := poMock.NewMockRepository(ctrl)

		u := NewUsecase(mockRepo, mockPoRepo)
		err := u.Delete("1234")

		assert.Error(t, err)
	})
}
