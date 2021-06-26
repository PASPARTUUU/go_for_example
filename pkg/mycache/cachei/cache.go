package cachei

// sources:
// https://habr.com/ru/post/359078/
// https://github.com/patrickmn/go-cache

import (
	"encoding/json"
	"errors"
	"sync"
	"time"
)

const (
	// NoExpiration - For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// DefaultExpiration - For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or (e.g. 5 minutes.)
	DefaultExpiration time.Duration = time.Minute * 5

	// NoCG -
	NoCG time.Duration = -1
	// DefaultCleanupInterval -
	DefaultCleanupInterval time.Duration = time.Minute * 5
)

// Cache -
type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[interface{}]item
}

// item -
type item struct {
	value      interface{}
	takenAt    time.Time
	expiration int64
}

// New -
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	// инициализируем карту(map) в паре ключ(string)/значение(Item)
	items := make(map[interface{}]item)

	if defaultExpiration == 0 {
		defaultExpiration = DefaultExpiration
	}
	if defaultExpiration < 0 {
		defaultExpiration = DefaultExpiration
	}
	if cleanupInterval == 0 {
		cleanupInterval = DefaultCleanupInterval
	}

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.StartGC()
	}

	return &cache
}

// Set -
func (c *Cache) Set(key interface{}, value interface{}) {
	c.set(key, value, 0)
}

// SetWithExpiration -
func (c *Cache) SetWithExpiration(key interface{}, value interface{}, duration time.Duration) {
	c.set(key, value, duration)
}

func (c *Cache) set(key interface{}, value interface{}, duration time.Duration) {

	var expiration int64

	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = c.defaultExpiration
	}

	if c.defaultExpiration < 0 {
		expiration = -1
	} else {
		// Устанавливаем время истечения кеша
		if duration > 0 {
			expiration = time.Now().Add(duration).UnixNano()
		}
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = item{
		value:      value,
		expiration: expiration,
		takenAt:    time.Now(),
	}

}

// Get -
func (c *Cache) Get(key interface{}) (interface{}, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return nil, false
	}

	item.takenAt = time.Now()

	return item.value, true
}

// Unmarshal - работает только для экспотрируемых полей
func (c *Cache) Unmarshal(key interface{}, res interface{}) (bool, error) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]
	if !found {
		return false, nil
	}
	item.takenAt = time.Now()

	b, err := json.Marshal(item.value)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(b, res)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Delete -
func (c *Cache) Delete(key interface{}) error {
	c.Lock()
	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}

// StartGC -
func (c *Cache) StartGC() {
	go c.runGC()
}

func (c *Cache) runGC() {

	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}
	}
}

// expiredKeys - возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []interface{}) {
	if c.defaultExpiration == NoExpiration {
		return keys
	}

	c.RLock()
	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.expiration && i.expiration != -1 {
			keys = append(keys, k)
		}
		if i.expiration == 0 && i.takenAt.UnixNano() > time.Now().Add(-time.Minute*5).UnixNano() {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems - удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []interface{}) {
	c.Lock()
	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
