package file

type Usecase interface {
	UpdateUser(id, filename string) error
	UpdateImageProfileUser(id, filename string) error
	GetPathTranscript(id string) (string, error)
	GetPathImageProfile(id string) (string, error)
}
