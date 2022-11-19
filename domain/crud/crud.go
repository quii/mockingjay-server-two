package crud

type CRUD[ID comparable, T any] interface {
	GetAll() ([]T, error)
	GetByID(ID) (T, bool, error)
	Create(T) error
	Delete(ID) error
}
