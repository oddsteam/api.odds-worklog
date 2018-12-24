package po

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	cusMock "gitlab.odds.team/worklog/api.odds-worklog/api/customer/mock"
	poMock "gitlab.odds.team/worklog/api.odds-worklog/api/po/mock"
)

func TestUsecase_Create(t *testing.T) {
	t.Run("when create po success, then return (*models.Po, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		po := poMock.Po
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		cusRepo.EXPECT().GetByID(po.CustomerId).Return(&cusMock.Customer, nil)
		poRepo.EXPECT().Create(&po).Return(&po, nil)

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.Create(&po)

		b, _ := json.Marshal(p)
		actaul := string(b)
		assert.NoError(t, err)
		assert.Equal(t, poMock.PoJson, actaul)
	})

	t.Run("when create po error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		po := poMock.Po
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		cusRepo.EXPECT().GetByID(po.CustomerId).Return(&cusMock.Customer, nil)
		poRepo.EXPECT().Create(&po).Return(nil, errors.New(""))

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.Create(&po)

		assert.Error(t, err)
		assert.Nil(t, p)
	})

	t.Run("when get customer error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		po := poMock.Po
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		cusRepo.EXPECT().GetByID(po.CustomerId).Return(nil, errors.New(""))

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.Create(&po)

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestUsecase_Get(t *testing.T) {
	t.Run("when get po's list success, then return ([]*models.Po, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		poes := poMock.Poes
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().Get().Return(poes, nil)

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.Get()

		b, _ := json.Marshal(p)
		actaul := string(b)
		assert.NoError(t, err)
		assert.Equal(t, poMock.PoesJson, actaul)
	})

	t.Run("when get po's list error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().Get().Return(nil, errors.New(""))

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.Get()

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestUsecase_GetByID(t *testing.T) {
	t.Run("when get po success, then return (*models.Po, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		po := poMock.Po
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().GetByID(po.ID.Hex()).Return(&po, nil)

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.GetByID(po.ID.Hex())

		b, _ := json.Marshal(p)
		actaul := string(b)
		assert.NoError(t, err)
		assert.Equal(t, poMock.PoJson, actaul)
	})

	t.Run("when get po error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		po := poMock.Po
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().GetByID(po.ID.Hex()).Return(nil, errors.New(""))

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.GetByID(po.ID.Hex())

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestUsecase_GetByCusID(t *testing.T) {
	t.Run("when get po's list by customer id success, then return ([]*models.Po, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		poes := poMock.Poes
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().GetByCusID("1234").Return(poes, nil)

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.GetByCusID("1234")

		b, _ := json.Marshal(p)
		actaul := string(b)
		assert.NoError(t, err)
		assert.Equal(t, poMock.PoesJson, actaul)
	})

	t.Run("when get po's list by customer id error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().GetByCusID("1234").Return(nil, errors.New(""))

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.GetByCusID("1234")

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestUsecase_Update(t *testing.T) {
	t.Run("when update po success, then return (*models.Po, nil)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		po := poMock.Po
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().GetByID(po.ID.Hex()).Return(&po, nil)
		poRepo.EXPECT().Update(&po).Return(&po, nil)

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.Update(&po)

		b, _ := json.Marshal(p)
		actaul := string(b)
		assert.NoError(t, err)
		assert.Equal(t, poMock.PoJson, actaul)
	})

	t.Run("when update po error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		po := poMock.Po
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().GetByID(po.ID.Hex()).Return(&po, nil)
		poRepo.EXPECT().Update(&po).Return(nil, errors.New(""))

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.Update(&po)

		assert.Error(t, err)
		assert.Nil(t, p)
	})

	t.Run("when get po error, then return (nil, error)", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		po := poMock.Po
		cusRepo := cusMock.NewMockRepository(ctrl)
		poRepo := poMock.NewMockRepository(ctrl)
		poRepo.EXPECT().GetByID(po.ID.Hex()).Return(nil, errors.New(""))

		u := NewUsecase(poRepo, cusRepo)
		p, err := u.Update(&po)

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}
