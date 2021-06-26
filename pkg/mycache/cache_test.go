package mycache

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

	cache.Set("a", "aaa")
	cache.SetWithExpiration("ex", "expire", time.Millisecond*1)

	item1, found := cache.Get("a")
	if !found || !reflect.DeepEqual("aaa", item1) {
		t.Error("Did not find elem even though it was set to never expire")
	}
	item2, found := cache.Get("ex")
	if !found || !reflect.DeepEqual("expire", item2) {
		t.Error("Did not find elem even though it was set to never expire")
	}

	<-time.After(time.Millisecond * 20)

	_, found = cache.Get("ex")
	if found {
		t.Error("Found 'ex' when it should have been automatically deleted")
	}

}
