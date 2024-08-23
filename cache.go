package cache

import (
	"slices"
	"sync"
	"time"
)

// key-value storage
type Cache[K comparable, V any] struct {
	mu   sync.Mutex
	data map[K]entryWithTimeout[V]
	ttl  time.Duration

	maxSize           int
	chronologicalKeys []K
}

// adding expiration date
type entryWithTimeout[V any] struct {
	value   V
	expires time.Time // after this time, the value is useless
}

// creating a cache
func New[K comparable, V any](maxSize int, ttl time.Duration) Cache[K, V] {
	return Cache[K, V]{
		data:              make(map[K]entryWithTimeout[V]),
		ttl:               ttl,
		maxSize:           maxSize,
		chronologicalKeys: make([]K, 0, maxSize),
	}
}

// inserting a key and value in the cache
func (c *Cache[K, V]) addKeyValue(key K, value V) {
	c.data[key] = entryWithTimeout[V]{
		value:   value,
		expires: time.Now().Add(c.ttl),
	}

	c.chronologicalKeys = append(c.chronologicalKeys, key)
}

// removes a key and its associated value from the cache.
func (c *Cache[K, V]) deleteKeYValue(key K) {
	c.chronologicalKeys = slices.DeleteFunc(c.chronologicalKeys, func(k K) bool { return k == key })
	delete(c.data, key)
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
		c.deleteKeYValue(key)
		return zeroV, false

	default:
		return e.value, true
	}
}

// overrides the value for current key
func (c *Cache[K, V]) Upsert(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, alreadyPresent := c.data[key]
	switch {
	case alreadyPresent:
		c.deleteKeYValue(key)
	case len(c.data) == c.maxSize:
		c.deleteKeYValue(c.chronologicalKeys[0])
	}

	c.addKeyValue(key, value)
}

// removes the entry for the specified key
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.deleteKeYValue(key)
}
