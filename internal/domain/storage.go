package domain

type Storage interface {
	Save(filename string, data interface{}) error
}
