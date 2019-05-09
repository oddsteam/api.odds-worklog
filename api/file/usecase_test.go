package file

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

func TestUsecase_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filename := "test.pdf"
	user := userMock.User
	mockUserRepo := userMock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

	user.Transcript = filename
	mockUserRepo.EXPECT().Update(&user).Return(&user, nil)

	usecase := NewUsecase(mockUserRepo)
	err := usecase.UpdateUser(user.ID.Hex(), filename)

	assert.NoError(t, err)
}

func TestUsecase_UpdateImageProfileUser(t *testing.T) {
	t.Run("happy case update image profile user is pass", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		filename := "test.jpg"
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

		user.ImageProfile = filename
		mockUserRepo.EXPECT().Update(&user).Return(&user, nil)

		usecase := NewUsecase(mockUserRepo)
		err := usecase.UpdateImageProfileUser(user.ID.Hex(), filename)

		assert.NoError(t, err)
	})

	t.Run("when searching for an ID not found in GetById should return err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		filename := "test.jpg"
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, errors.New("error"))

		usecase := NewUsecase(mockUserRepo)
		err := usecase.UpdateImageProfileUser(user.ID.Hex(), filename)

		assert.Error(t, err)

	})

	t.Run("when update data fail should return err", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		filename := "test.jpg"
		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

		user.ImageProfile = filename
		mockUserRepo.EXPECT().Update(&user).Return(&user, errors.New("error"))
		usecase := NewUsecase(mockUserRepo)
		err := usecase.UpdateImageProfileUser(user.ID.Hex(), filename)

		assert.Error(t, err)

	})
}

func TestUsecase_UpdateDegreeCertificate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filename := "test.pdf"
	user := userMock.User
	mockUserRepo := userMock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

	user.DegreeCertificate = filename
	mockUserRepo.EXPECT().Update(&user).Return(&user, nil)

	usecase := NewUsecase(mockUserRepo)
	err := usecase.UpdateDegreeCertificate(user.ID.Hex(), filename)

	assert.NoError(t, err)
}

func TestUsecase_UpdateIDCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filename := "test.pdf"
	user := userMock.User
	mockUserRepo := userMock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

	user.IDCard = filename
	mockUserRepo.EXPECT().Update(&user).Return(&user, nil)

	usecase := NewUsecase(mockUserRepo)
	err := usecase.UpdateIDCard(user.ID.Hex(), filename)

	assert.NoError(t, err)
}

func TestUsecase_GetPathTranscript(t *testing.T) {
	t.Run("get path transcript success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User
		user.Transcript = "test.pdf"
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

		usecase := NewUsecase(mockUserRepo)
		filename, err := usecase.GetPathTranscript(user.ID.Hex())

		assert.NoError(t, err)
		assert.NotEmpty(t, filename)
	})

	t.Run("when transcript empty then return '', ErrNoTranscriptFile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

		usecase := NewUsecase(mockUserRepo)
		filename, err := usecase.GetPathTranscript(user.ID.Hex())

		assert.EqualError(t, err, utils.ErrNoTranscriptFile.Error())
		assert.Empty(t, filename)
	})

	t.Run(`when can't open transcript file 
			then update user with empty transcript 
			and return '', ErrNoTranscriptFile`,
		func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user := userMock.User
			user.Transcript = "no.pdf"
			mockUserRepo := userMock.NewMockRepository(ctrl)
			mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
			mockUserRepo.EXPECT().Update(gomock.Any())

			usecase := NewUsecase(mockUserRepo)
			filename, err := usecase.GetPathTranscript(user.ID.Hex())

			assert.EqualError(t, err, utils.ErrNoTranscriptFile.Error())
			assert.Empty(t, filename)
		})
}

func TestUsecase_GetPathImageProfile(t *testing.T) {
	t.Run("when ImageProfile empty then return '', ErrNoImageProfileFile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

		usecase := NewUsecase(mockUserRepo)
		filename, err := usecase.GetPathImageProfile(user.ID.Hex())

		assert.EqualError(t, err, utils.ErrNoImageProfileFile.Error())
		assert.Empty(t, filename)
	})

	t.Run(`when can't open ImageProfile file  then update user with empty ImageProfile and return '', ErrNoImageProfileFile`,
		func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user := userMock.User
			user.ImageProfile = "noimg.png"
			mockUserRepo := userMock.NewMockRepository(ctrl)
			mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
			mockUserRepo.EXPECT().Update(gomock.Any())

			usecase := NewUsecase(mockUserRepo)
			filename, err := usecase.GetPathImageProfile(user.ID.Hex())

			assert.EqualError(t, err, utils.ErrNoImageProfileFile.Error())
			assert.Empty(t, filename)
		})
}

func TestUsecase_GetPathDegreeCertificate(t *testing.T) {
	t.Run("get path degree certificate success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User
		user.DegreeCertificate = "test.pdf"
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

		usecase := NewUsecase(mockUserRepo)
		filename, err := usecase.GetPathDegreeCertificate(user.ID.Hex())

		assert.NoError(t, err)
		assert.NotEmpty(t, filename)
	})

	t.Run("when degree certificate empty then return '', ErrNoDegreeCertificateFile", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		user := userMock.User
		mockUserRepo := userMock.NewMockRepository(ctrl)
		mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)

		usecase := NewUsecase(mockUserRepo)
		filename, err := usecase.GetPathDegreeCertificate(user.ID.Hex())

		assert.EqualError(t, err, utils.ErrNoDegreeCertificateFile.Error())
		assert.Empty(t, filename)
	})

	t.Run(`when can't open degree certificate file
			then update user with empty degree certificate
			and return '', ErrNoDegreeCertificateFile`,
		func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			user := userMock.User
			user.DegreeCertificate = "no.pdf"
			mockUserRepo := userMock.NewMockRepository(ctrl)
			mockUserRepo.EXPECT().GetByID(user.ID.Hex()).Return(&user, nil)
			mockUserRepo.EXPECT().Update(gomock.Any())

			usecase := NewUsecase(mockUserRepo)
			filename, err := usecase.GetPathDegreeCertificate(user.ID.Hex())

			assert.EqualError(t, err, utils.ErrNoDegreeCertificateFile.Error())
			assert.Empty(t, filename)
		})
}
