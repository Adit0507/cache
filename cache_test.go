package cache_test

import (
	"testing"
	"time"
	"caching-generics"
	"github.com/stretchr/testify/assert"
)

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
