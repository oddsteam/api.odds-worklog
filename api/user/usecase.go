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
	user, err := u.repo.GetUserByEmail(m.Email)
	if err == nil {
		return user, utils.ErrConflict
	}

	return u.repo.CreateUser(m)
}

func (u *usecase) GetUser() ([]*models.User, error) {
	return u.repo.GetUser()
}

func (u *usecase) GetUserByRole(role string) ([]*models.User, error) {
	return u.repo.GetUserByRole(role)
}

func (u *usecase) GetUserByID(id string) (*models.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *usecase) GetUserBySiteID(id string) ([]*models.User, error) {
	return u.repo.GetUserBySiteID(id)
}

func (u *usecase) UpdateUser(m *models.User) (*models.User, error) {
	if err := m.ValidateRole(); err != nil {
		return nil, err
	}
	if m.Role == "admin" && !m.IsAdmin() {
		return nil, utils.ErrInvalidUserRole
	}
	if m.IsAdmin() && m.Role != "admin" {
		m.Role = "admin"
	}
	user, err := u.repo.UpdateUser(m)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}
