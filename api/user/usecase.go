package user

import (
	"unicode"

	"gitlab.odds.team/worklog/api.odds-worklog/api/site"
	"gitlab.odds.team/worklog/api.odds-worklog/business/models"
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

	siteByID := make(map[string]*models.Site, len(sites))
	for _, s := range sites {
		siteByID[s.ID.Hex()] = s
	}
	for i, us := range users {
		if s, ok := siteByID[us.SiteID]; ok {
			users[i].Site = s
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

func (u *usecase) Update(userFromRequest *models.User, actor *models.UserClaims) (*models.User, error) {
	currentUser, err := u.repo.GetByID(userFromRequest.ID.Hex())
	if err != nil {
		return nil, err
	}

	if actor == nil {
		return nil, utils.ErrPermissionDenied
	}
	isSelf := actor.ID == currentUser.ID.Hex()
	isAdmin := actor.IsAdmin()
	isUserAdmin := actor.Role == "user-admin"
	if !isSelf && !isAdmin && !isUserAdmin {
		return nil, utils.ErrPermissionDenied
	}

	if isUserAdmin && !isSelf {
		updated := *currentUser
		if userFromRequest.SiteID != "" {
			updated.SiteID = userFromRequest.SiteID
		}
		return u.repo.Update(&updated)
	}

	// Allow partial updates (e.g. site-only) to omit role/vat; keep current values.
	if userFromRequest.Role == "" {
		userFromRequest.Role = currentUser.Role
	}
	if userFromRequest.Vat == "" {
		userFromRequest.Vat = currentUser.Vat
	}

	if err := userFromRequest.ValidateRole(); err != nil {
		return nil, err
	}
	if err := userFromRequest.ValidateVat(); err != nil {
		return nil, err
	}
	// Only full admin may promote a user TO admin; user-admin may still update existing admins (e.g. site).
	if userFromRequest.Role == "admin" && !isAdmin && currentUser.Role != "admin" {
		return nil, utils.ErrInvalidUserRole
	}

	user := NewUser(*currentUser)
	err = user.prepareDataForUpdateFrom(*userFromRequest)
	if err != nil {
		return nil, err
	}

	persistedUser, err := u.repo.Update(user.data)
	if err != nil {
		return nil, err
	}

	return persistedUser, nil
}

func extractNumbers(input string) string {
	var result []rune

	for _, char := range input {
		if unicode.IsDigit(char) {
			result = append(result, char)
		}
	}

	return string(result)
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
