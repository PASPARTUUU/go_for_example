package cachei

import (
	"reflect"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	type keyStruct struct {
		id  int
		val int
	}

	cache := New(0, time.Millisecond)

	cache.Set(1, 111)
	cache.Set("a", "aaa")
	cache.Set(keyStruct{id: 1, val: 1}, "it a struct")
	cache.SetWithExpiration("ex", "expire", time.Millisecond*1)

	item1, found := cache.Get(1)
	if !found || !reflect.DeepEqual(111, item1) {
		t.Error("Did not find elem even though it was set to never expire")
	}
	item2, found := cache.Get("a")
	if !found || !reflect.DeepEqual("aaa", item2) {
		t.Error("Did not find elem even though it was set to never expire")
	}
	item3, found := cache.Get(keyStruct{id: 1, val: 1})
	if !found || !reflect.DeepEqual("it a struct", item3) {
		t.Error("Did not find elem even though it was set to never expire")
	}
	item4, found := cache.Get("ex")
	if !found || !reflect.DeepEqual("expire", item4) {
		t.Error("Did not find elem even though it was set to never expire")
	}

	<-time.After(time.Millisecond * 20)

	_, found = cache.Get("ex")
	if found {
		t.Error("Found 'ex' when it should have been automatically deleted")
	}

}
