package file

type Usecase interface {
	UpdateUser(id, filename string) error
	GetPathTranscript(id string) (string, error)
}
