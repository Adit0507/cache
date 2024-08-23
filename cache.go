package cache

import (
	"sync"
	"time"
)

// key-value storage
type Cache[K comparable, V any] struct {
	mu   sync.Mutex
	data map[K]entryWithTimeout[V]	
	ttl time.Duration
}

// adding expiration date
type entryWithTimeout[V any] struct {
	value   V
	expires time.Time // after this time, the value is useless
}

// creating a cache
func New[K comparable, V any](ttl time.Duration) Cache[K, V] {
	return Cache[K, V]{
		data: make(map[K]entryWithTimeout[V]),
		ttl:  ttl,
	}
}

// readingfrom cache
func (c *Cache[K, V]) Read(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var zeroV V
	e, ok := c.data[key]

	switch {
	case !ok:
		return zeroV, false
				
	case e.expires.Before(time.Now()):
		// since the Read() method is now altering the content
		// we cant use RWMutex  
		delete(c.data, key)
		return zeroV, false

	default:
		return e.value, true
	}
}

// overrides the value for current key
func (c *Cache[K, V]) Upsert(key K, value V) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = entryWithTimeout[V]{
		value:   value,
		expires: time.Now().Add(c.ttl),
	}

	return nil
}

// deleting key
func (c *Cache[K, V]) Delete(key K) {
	delete(c.data, key)
}
