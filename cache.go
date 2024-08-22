package main

import (
	"sync"
	"time"
)

// key-value storage
type Cache[K comparable, V any] struct {
	mu sync.Mutex
	data map[K]V
}

// ttl
type entryWithTimeout[V any] struct {
	value V
	expires time.Time	// after this time, the value is useless
}

// creating a cache
func New[K comparable, V any]() Cache[K, V] {
	return Cache[K, V]{
		data: make(map[K]V),
	}
}

// readingfrom cache
func (c *Cache[K, V])  Read(key K) (V, bool) {
	v, found := c.data[key]

	return v, found
}

// overrides the value for current key
func (c *Cache[K, V]) Upsert(key K, value V) error {
	c.data[key] = value

	return nil
}

// deleting key
func (c *Cache[K, V]) Delete(key K) {
	delete(c.data, key)
}