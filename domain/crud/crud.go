package crud

type CRUDesque[ID comparable, T any] interface {
	GetAll() ([]T, error)
	GetByID(ID) (T, bool, error)
	Create(ID, T) error
	Delete(ID) error
}

type CRUD[ID comparable, T any] struct {
	dict map[ID]T
}

func NewCRUD[ID comparable, T any]() *CRUD[ID, T] {
	items := make(map[ID]T)
	return &CRUD[ID, T]{dict: items}
}

func (r *CRUD[ID, T]) GetAll() ([]T, error) {
	var items []T
	for _, t := range r.dict {
		items = append(items, t)
	}
	return items, nil
}

func (r *CRUD[ID, T]) GetByID(id ID) (T, bool, error) {
	item, exists := r.dict[id]
	return item, exists, nil
}

func (r *CRUD[ID, T]) Create(id ID, t T) error {
	r.dict[id] = t
	return nil
}

func (r *CRUD[ID, T]) Delete(id ID) error {
	delete(r.dict, id)
	return nil
}
