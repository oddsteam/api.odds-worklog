package user

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"strings"

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

func (u *usecase) GetUserByType(corporateFlag string) ([]*models.User, error) {
	if corporateFlag != "Y" && corporateFlag != "N" {
		return nil, utils.ErrInvalidFlag
	}
	return u.repo.GetUserByType(corporateFlag)
}

func (u *usecase) GetUserByID(id string) (*models.User, error) {
	return u.repo.GetUserByID(id)
}

func (u *usecase) UpdateUser(m *models.User, file *multipart.FileHeader) (*models.User, error) {
	user, err := u.repo.UpdateUser(m)
	if err != nil {
		return nil, errors.New("Update user failed")
	}

	if file != nil {
		filename := getTranscriptFilename(user)
		err = saveTranscript(file, filename)
		if err != nil {
			return nil, errors.New("Save transcript failed")
		}
	}
	return user, nil
}

func (u *usecase) DeleteUser(id string) error {
	return u.repo.DeleteUser(id)
}

func saveTranscript(file *multipart.FileHeader, filename string) (err error) {
	if file == nil {
		return errors.New("File is nil!")
	}

	src, err := file.Open()
	if err != nil {
		return
	}
	defer src.Close()

	// Destination
	dst, err := os.Create(filename)
	if err != nil {
		return
	}
	defer dst.Close()

	// Copy
	_, err = io.Copy(dst, src)
	return
}

func getTranscriptFilename(u *models.User) (filename string) {
	path := "files/transcripts"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}

	r := utils.RandStringBytes(8)
	filename = fmt.Sprintf("%s/transcript_%s_%s_%s.pdf", path, strings.ToUpper(u.FirstName), strings.ToUpper(u.LastName), r)
	return
}
