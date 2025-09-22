package memcache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpiryQueue(t *testing.T) {
	cache := NewCache()

	t.Run("Handles queue expiry", func(t *testing.T) {
		cache.Set("1", 1, 25)
		cache.Set("2", 2, 0)
		cache.Set("3", 3, 150)

		_, ok1 := cache.Get("1")
		_, ok2 := cache.Get("2")
		_, ok3 := cache.Get("3")

		assert.Equal(t, ok1, true)
		assert.Equal(t, ok2, true)
		assert.Equal(t, ok3, true)

		time.Sleep(100 * time.Millisecond)

		_, ok1 = cache.Get("1")
		_, ok2 = cache.Get("2")
		_, ok3 = cache.Get("3")
		assert.Equal(t, ok1, false)
		assert.Equal(t, ok2, true)
		assert.Equal(t, ok3, true)
	
		time.Sleep(100 * time.Millisecond)
		_, ok1 = cache.Get("1")
		_, ok2 = cache.Get("2")
		_, ok3 = cache.Get("3")
		assert.Equal(t, ok1, false)
		assert.Equal(t, ok2, true)
		assert.Equal(t, ok3, false)
	})
}
