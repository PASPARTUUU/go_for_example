package myorganizer

import (
	"reflect"
	"testing"
)

type customCache struct {
	items map[string]interface{}
}

func (cc customCache) New() *customCache {
	return &customCache{
		items: make(map[string]interface{}),
	}
}
func (cc customCache) Set(key string, item interface{}) {
	cc.items[key] = item
}
func (cc customCache) Get(key string) interface{} {
	val, _ := cc.items[key]
	return val
}

func TestT(t *testing.T) {
	ccache := customCache{}.New()

	type area struct {
		id string

		Fold *Fold
	}

	var rootArea = area{
		id:   "111",
		Fold: NewRoot(ccache, "111"),
	}

	var area2 = area{
		id:   "222",
		Fold: rootArea.Fold.New("222"),
	}
	var area3 = area{
		id:   "333",
		Fold: rootArea.Fold.New("333"),
	}
	var area4 = area{
		id:   "444",
		Fold: area3.Fold.New("444"),
	}

	if area2.Fold.Parent != rootArea.Fold.ID {
		t.Error("the child has the wrong parent")
	}
	if area3.Fold.Parent != rootArea.Fold.ID {
		t.Error("the child has the wrong parent")
	}

	if eq := reflect.DeepEqual(rootArea.Fold.GetItemByID(area2.Fold.ID), *area2.Fold); eq == false {
		t.Error("fold does not match")
	}

	childs := rootArea.Fold.GetChildrens()

	for _, el := range childs {
		if el.ID != area2.Fold.ID &&
			el.ID != area3.Fold.ID &&
			el.ID != area4.Fold.ID {
			t.Error("the child has the wrong parent")
		}
	}

	// TODO: проверить
	// (it *Fold) SetElement(elem interface{})
	// (it *Fold) GetElements() []interface{}
}
