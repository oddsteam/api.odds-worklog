package file

import (
	"os"

	"gitlab.odds.team/worklog/api.odds-worklog/api/user"
	"gitlab.odds.team/worklog/api.odds-worklog/pkg/utils"
)

type usercasse struct {
	repo user.Repository
}

func NewUsecase(repo user.Repository) Usecase {
	return &usercasse{repo}
}

func (u *usercasse) UpdateUser(id, filename string) error {
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

func (u *usercasse) GetPathTranscript(id string) (string, error) {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return "", err
	}
	if user.Transcript == "" {
		return "", utils.ErrNoTranscriptFile
	}

	_, err = os.Open(user.Transcript)
	if err != nil {
		user.Transcript = ""
		u.repo.UpdateUser(user)
		return "", utils.ErrNoTranscriptFile
	}
	return user.Transcript, nil
}

func (u *usercasse) UpdateImageProfileUser(id, filename string) error {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	user.ImageProfile = filename
	user, err = u.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
