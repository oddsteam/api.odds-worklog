package file

type Usecase interface {
	UpdateUser(id, filename string) error
}
