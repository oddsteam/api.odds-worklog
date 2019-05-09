package file

import "gitlab.odds.team/worklog/api.odds-worklog/models"

type Usecase interface {
	UpdateUser(id, filename string) error
	UpdateImageProfileUser(id, filename string) error
	UpdateDegreeCertificate(id, filename string) error
	UpdateIDCard(id, filename string) error
	GetPathTranscript(id string) (string, error)
	GetPathImageProfile(id string) (string, error)
	GetPathDegreeCertificate(id string) (string, error)
	RemoveTranscript(filename string) error
	RemoveDegreeCertificate(filename string) error
	RemoveImage(filename string) error
	GetUserByID(id string) (*models.User, error)
}
