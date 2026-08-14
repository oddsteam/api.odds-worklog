package usecases

import "gitlab.odds.team/worklog/api.odds-worklog/business/models"

type manageUsersUsecase struct {
	users ForManagingUsers
	sites ForListingSites
}

func NewManageUsersUsecase(users ForManagingUsers, sites ForListingSites) ForUsingUsers {
	return &manageUsersUsecase{users, sites}
}

func (u *manageUsersUsecase) Create(m *models.User) (*models.User, error) {
	if err := m.ValidateEmail(); err != nil {
		return nil, err
	}
	user, err := u.users.GetByEmail(m.Email)
	if err == nil {
		return user, models.ErrConflict
	}

	return u.users.Create(m)
}

func (u *manageUsersUsecase) Get() ([]*models.User, error) {
	users, err := u.users.Get()
	if err != nil {
		return nil, err
	}

	sites, err := u.sites.GetSiteGroup()
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

func (u *manageUsersUsecase) GetByRole(role string) ([]*models.User, error) {
	return u.users.GetByRole(role)
}

func (u *manageUsersUsecase) GetByID(id string) (*models.User, error) {
	return u.users.GetByID(id)
}

func (u *manageUsersUsecase) GetByEmail(email string) (*models.User, error) {
	return u.users.GetByEmail(email)
}

func (u *manageUsersUsecase) GetBySiteID(id string) ([]*models.User, error) {
	return u.users.GetBySiteID(id)
}

func (u *manageUsersUsecase) Update(userFromRequest *models.User, actor *models.UserClaims) (*models.User, error) {
	currentUser, err := u.users.GetByID(userFromRequest.ID.Hex())
	if err != nil {
		return nil, err
	}

	if actor == nil {
		return nil, models.ErrPermissionDenied
	}
	isSelf := actor.ID == currentUser.ID.Hex()
	isAdmin := actor.IsAdmin()
	isUserAdmin := actor.Role == "user-admin"
	if !isSelf && !isAdmin && !isUserAdmin {
		return nil, models.ErrPermissionDenied
	}

	if isUserAdmin && !isSelf {
		updated := *currentUser
		if userFromRequest.SiteID != "" {
			updated.SiteID = userFromRequest.SiteID
		}
		return u.users.Update(&updated)
	}

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
	if userFromRequest.Role == "admin" && !isAdmin && currentUser.Role != "admin" {
		return nil, models.ErrInvalidUserRole
	}

	updated := *currentUser
	updated.ApplyProfileUpdate(*userFromRequest)
	return u.users.Update(&updated)
}

func (u *manageUsersUsecase) UpdateStatusTavi(m []*models.StatusTavi, isAdmin bool) ([]*models.User, error) {
	var users []*models.User

	for i := 0; i < len(m); i++ {
		if err := m[i].User.ValidateRole(); err != nil {
			return nil, err
		}
		if m[i].User.Role == "admin" && !isAdmin {
			return nil, models.ErrInvalidUserRole
		}

		user, err := u.users.GetByID(m[i].User.ID.Hex())
		if err != nil {
			return nil, err
		}

		user.StatusTavi = m[i].User.StatusTavi
		user.Role = user.Role
		user.Vat = user.Vat

		user, err = u.users.Update(user)
		if err != nil {
			return nil, err
		}
		user.DailyIncome = ""
		user.ThaiCitizenID = ""

		users = append(users, user)
	}
	return users, nil
}

func (u *manageUsersUsecase) Delete(id string) error {
	user, err := u.users.GetByID(id)
	if err != nil {
		return err
	}
	_, err = u.users.CreateArchivedUser(*user)
	if err != nil {
		return err
	}
	return u.users.Delete(id)
}
