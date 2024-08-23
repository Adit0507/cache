package cache_test

import (
	"testing"
	"time"
	"caching-generics"
	"github.com/stretchr/testify/assert"
)

func TestCache_MaxSize(t*testing.T) {
	t.Parallel()

	c := cache.New[int, int](3, time.Minute)
	c.Upsert(1,1)
	c.Upsert(2,2)
	c.Upsert(3,3)
	
	got, found := c.Read(1)
	assert.True(t, found)
	assert.Equal(t, 1, got)

	c.Upsert(1, 10)
	c.Upsert(4,4)

	got, found = c.Read(2)
	assert.False(t, found)
	assert.Equal(t, 0, got)
}

func TestCache_TTL(t *testing.T) {
	t.Parallel()
	c := cache.New[string, string](5, time.Millisecond*100)
	c.Upsert("Norwegian", "Blue")

	got, found := c.Read("Norwegian")
	assert.True(t, found)
	assert.Equal(t, "Blue", got)

	time.Sleep(time.Millisecond * 200)
	got, found = c.Read("Norwegian")
	
	assert.False(t, found)
	assert.Equal(t, "", got)
}
