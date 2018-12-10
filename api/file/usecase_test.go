package file

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
	userMock "gitlab.odds.team/worklog/api.odds-worklog/api/user/mock"
)

func TestUsecase_UpdateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	filename := "test.pdf"
	user := userMock.MockUser
	mockUserRepo := userMock.NewMockRepository(ctrl)
	mockUserRepo.EXPECT().GetUserByID(user.ID.Hex()).Return(&user, nil)

	user.Transcript = filename
	mockUserRepo.EXPECT().UpdateUser(&user).Return(&user, nil)

	usecase := NewUsecase(mockUserRepo)
	err := usecase.UpdateUser(user.ID.Hex(), filename)

	assert.NoError(t, err)

}
