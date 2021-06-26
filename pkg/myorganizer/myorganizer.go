package myorganizer

import (
	"github.com/gofrs/uuid"
)

type IStorage interface {
	Set(key string, item interface{})
	Get(key string) interface{}
}

// Fold -
type Fold struct {
	ID string

	Parent    string
	Childrens map[string]struct{}

	Elems []interface{}

	storage IStorage
}

// NewRoot -
func NewRoot(storage IStorage, customID ...string) *Fold {
	var id string = uuid.Must(uuid.NewV4()).String()
	if len(customID) != 0 {
		if customID[0] != "" {
			id = customID[0]
		}
	}

	item := Fold{
		ID:        id,
		Parent:    "",
		Childrens: make(map[string]struct{}),
		storage:   storage,
	}

	item.storage.Set(id, item)

	return &item
}

// New -
func (it *Fold) New(customID ...string) *Fold {
	var id string = uuid.Must(uuid.NewV4()).String()
	if len(customID) != 0 {
		if customID[0] != "" {
			id = customID[0]
		}
	}

	it.Childrens[id] = struct{}{}

	item := Fold{
		ID:        id,
		Parent:    it.ID,
		Childrens: make(map[string]struct{}),
		storage:   it.storage,
	}

	it.storage.Set(id, item)

	return &item
}

// GetChildren -
func (it *Fold) GetChildrens() []Fold {
	var childrens []Fold

	it.getChildrens(&childrens)

	return childrens
}

func (it *Fold) getChildrens(ppl *[]Fold) {
	for t := range it.Childrens {
		child := it.GetItemByID(t)

		*ppl = append(*ppl, child)

		child.getChildrens(ppl)
	}
}

// SetElement -
func (it *Fold) SetElement(elem interface{}) {
	item := it.GetItemByID(it.ID)
	item.Elems = append(item.Elems, elem)

	*it = item

	it.storage.Set(it.ID, item)
}

// GetElements -
func (it *Fold) GetElements() []interface{} {
	var res []interface{}
	res = fill(it.Elems, res)
	items := it.GetChildrens()

	for _, el := range items {
		res = fill(el.Elems, res)
	}

	return res
}

// GetItemByID - получение из хранилища без ограничений видимости
func (it *Fold) GetItemByID(id string) Fold {
	return it.storage.Get(id).(Fold)
}

func fill(src, dst []interface{}) []interface{} {
	for _, el := range src {
		dst = append(dst, el)
	}
	return dst
}
