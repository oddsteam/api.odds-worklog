package user

import (
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	repo Repository
}

func NewUsecase(r Repository) Usecase {
	return &usecase{r}
}

func (u *usecase) CreateUser(m *models.User) (*models.User, error) {
	err := utils.ValidateEmail(m.Email)
	if err != nil {
		return nil, err
	}
	_, err = u.repo.GetUserByEmail(m.Email)
	if err != nil {
		return nil, err
	}

	return u.repo.CreateUser(m)
}

func (u *usecase) GetUser() ([]*models.User, error) {
	return u.repo.GetUser()
}

func (u *usecase) GetUserByType(corporateFlag string) ([]*models.User, error) {
	return u.repo.GetUserByType(corporateFlag)
}

func (u *usecase) GetUserByID(id string) (*models.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *usecase) UpdateUser(m *models.User) (*models.User, error) {
	return u.repo.UpdateUser(m)
}

func (u *usecase) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}
