package file

import (
	"os"

	"gitlab.odds.team/worklog/api.odds-worklog/models"

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
	user, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}

	user.Transcript = filename
	user, err = u.repo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *usercasse) UpdateImageProfileUser(id, filename string) error {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	user.ImageProfile = filename
	user, err = u.repo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *usercasse) UpdateDegreeCertificate(id, filename string) error {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}
	user.DegreeCertificate = filename
	user, err = u.repo.Update(user)
	if err != nil {
		return err
	}

	return nil
}

func (u *usercasse) GetPathTranscript(id string) (string, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	if user.Transcript == "" {
		return "", utils.ErrNoTranscriptFile
	}

	_, err = os.Open(user.Transcript)
	if err != nil {
		user.Transcript = ""
		u.repo.Update(user)
		return "", utils.ErrNoTranscriptFile
	}
	return user.Transcript, nil
}

func (u *usercasse) GetPathImageProfile(id string) (string, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	if user.ImageProfile == "" {
		return "", utils.ErrNoImageProfileFile
	}

	_, err = os.Open(user.ImageProfile)
	if err != nil {
		user.ImageProfile = ""
		u.repo.Update(user)
		return "", utils.ErrNoImageProfileFile
	}
	return user.ImageProfile, nil
}

func (u *usercasse) GetPathDegreeCertificate(id string) (string, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	if user.DegreeCertificate == "" {
		return "", utils.ErrNoDegreeCertificateFile
	}

	_, err = os.Open(user.DegreeCertificate)
	if err != nil {
		user.DegreeCertificate = ""
		u.repo.Update(user)
		return "", utils.ErrNoDegreeCertificateFile
	}
	return user.DegreeCertificate, nil
}

func (u *usercasse) RemoveTranscript(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func (u *usercasse) RemoveImage(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

func (u *usercasse) GetUserByID(id string) (*models.User, error) {
	return u.repo.GetByID(id)
}
