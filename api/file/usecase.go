package file

import (
	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
)

type usercasse struct {
	repo user.Repository
}

func NewUsecase(repo user.Repository) Usecase {
	return &usercasse{repo}
}

func (u usercasse) UpdateUser(id, filename string) error {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	user.Transcript = filename
	user, err = u.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
