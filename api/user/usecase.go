package user

import (
	"errors"

	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/models"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usecase struct {
	repo     Repository
	siteRepo site.Repository
}

func NewUsecase(r Repository, sr site.Repository) Usecase {
	return &usecase{r, sr}
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
	users, err := u.repo.GetUser()
	if err != nil {
		return nil, err
	}

	sites, err := u.siteRepo.GetSiteGroup()
	if err != nil {
		return nil, err
	}

	for _, s := range sites {
		for i, us := range users {
			if s.ID.Hex() == us.SiteID {
				users[i].Site = s
				users[i].SiteID = ""
				break
			}
		}
	}
	return users, nil
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

func (u *usecase) UpdateUser(m *models.User, isAdmin bool) (*models.User, error) {
	if err := m.ValidateRole(); err != nil {
		return nil, err
	}
	if err := m.ValidateVat(); err != nil {
		return nil, err
	}
	if m.Role == "admin" && !isAdmin {
		return nil, utils.ErrInvalidUserRole
	}
	if m.IsAdmin() && m.Role != "admin" {
		m.Role = "admin"
	}

	user, err := u.repo.GetUserByID(m.ID.Hex())
	if err != nil {
		return nil, err
	}

	if m.FirstName != "" {
		user.FirstName = utils.ToFirstUpper(m.FirstName)
	}
	if m.LastName != "" {
		user.LastName = utils.ToFirstUpper(m.LastName)
	}
	if m.BankAccountName != "" {
		user.BankAccountName = m.BankAccountName
	}
	if m.BankAccountNumber != "" {
		user.BankAccountNumber = m.BankAccountNumber
	}
	if m.ThaiCitizenID != "" {
		user.ThaiCitizenID = m.ThaiCitizenID
	}
	if m.SlackAccount != "" {
		if err := utils.ValidateEmail(m.SlackAccount); err != nil {
			return nil, errors.New("Invalid slack acount.")
		}
		user.SlackAccount = m.SlackAccount
	}
	user.Role = m.Role
	user.Vat = m.Vat

	user, err = u.repo.UpdateUser(m)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}
