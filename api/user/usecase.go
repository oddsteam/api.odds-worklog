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

func (u *usecase) Create(m *models.User) (*models.User, error) {
	err := utils.ValidateEmail(m.Email)
	if err != nil {
		return nil, err
	}
	user, err := u.repo.GetByEmail(m.Email)
	if err == nil {
		return user, utils.ErrConflict
	}

	return u.repo.Create(m)
}

func (u *usecase) Get() ([]*models.User, error) {
	users, err := u.repo.Get()
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

func (u *usecase) GetByRole(role string) ([]*models.User, error) {
	return u.repo.GetByRole(role)
}

func (u *usecase) GetByID(id string) (*models.User, error) {
	return u.repo.GetByID(id)
}

func (u *usecase) GetByEmail(email string) (*models.User, error) {
	return u.repo.GetByEmail(email)
}

func (u *usecase) GetBySiteID(id string) ([]*models.User, error) {
	return u.repo.GetBySiteID(id)
}

func (u *usecase) Update(m *models.User, isAdmin bool) (*models.User, error) {
	if err := m.ValidateRole(); err != nil {
		return nil, err
	}
	if err := m.ValidateVat(); err != nil {
		return nil, err
	}
	if m.Role == "admin" && !isAdmin {
		return nil, utils.ErrInvalidUserRole
	}

	user, err := u.repo.GetByID(m.ID.Hex())
	if err != nil {
		return nil, err
	}

	if m.FirstName != "" {
		user.FirstName = utils.ToFirstUpper(m.FirstName)
	}
	if m.LastName != "" {
		user.LastName = utils.ToFirstUpper(m.LastName)
	}
	if m.CorporateName != "" {
		user.CorporateName = m.CorporateName
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
	if m.SiteID != "" {
		user.SiteID = m.SiteID
	}
	if m.Project != "" {
		user.Project = m.Project
	}
	if m.DailyIncome != "" {
		user.DailyIncome = m.DailyIncome
	}
	if m.Address != "" {
		user.Address = m.Address
	}
	if m.StartDate != "" {
		user.StartDate = m.StartDate
	}

	user.StatusTavi = m.StatusTavi
	user.Role = m.Role
	user.Vat = m.Vat
	user.Phone = m.Phone
	user, err = u.repo.Update(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *usecase) UpdateStatusTavi(m []*models.StatusTavi, isAdmin bool) ([]*models.User, error) {
	var users []*models.User

	for i := 0; i < len(m); i++ {
		if err := m[i].User.ValidateRole(); err != nil {
			return nil, err
		}
		// if err := m[i].User.ValidateVat(); err != nil {
		// 	return nil, err
		// }
		if m[i].User.Role == "admin" && !isAdmin {
			return nil, utils.ErrInvalidUserRole
		}

		user, err := u.repo.GetByID(m[i].User.ID.Hex())
		if err != nil {
			return nil, err
		}

		if m[i].User.FirstName != "" {
			user.FirstName = user.FirstName
		}
		if m[i].User.LastName != "" {
			user.LastName = user.LastName
		}
		if m[i].User.CorporateName != "" {
			user.CorporateName = user.CorporateName
		}
		if m[i].User.BankAccountName != "" {
			user.BankAccountName = user.BankAccountName
		}
		if m[i].User.BankAccountNumber != "" {
			user.BankAccountNumber = user.BankAccountNumber
		}
		if m[i].User.ThaiCitizenID != "" {
			user.ThaiCitizenID = user.ThaiCitizenID
		}
		if m[i].User.SlackAccount != "" {
			if err := utils.ValidateEmail(m[i].User.SlackAccount); err != nil {
				return nil, errors.New("Invalid slack acount.")
			}
			user.SlackAccount = user.SlackAccount
		}
		if m[i].User.SiteID != "" {
			user.SiteID = user.SiteID
		}
		if m[i].User.Project != "" {
			user.Project = user.Project
		}
		if m[i].User.DailyIncome != "" {
			user.DailyIncome = user.DailyIncome
		}
		if m[i].User.Address != "" {
			user.Address = user.Address
		}

		user.StatusTavi = m[i].User.StatusTavi
		user.Role = user.Role
		user.Vat = user.Vat

		user, err = u.repo.Update(user)
		if err != nil {
			return nil, err
		}
		user.DailyIncome = ""
		user.ThaiCitizenID = ""

		users = append(users, user)

	}
	return users, nil
}

func (u *usecase) Delete(id string) error {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	_, err = u.repo.CreateArchivedUser(*user)
	if err != nil {
		return err
	}
	return u.repo.Delete(id)
}
