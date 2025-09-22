package memcache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache := NewCache()

	t.Run("Can set/get items in cache", func (t *testing.T) {
		cache.Set("pelle", "svanslös", 0)

		val, ok := cache.Get("pelle")
		assert.Equal(t, ok, true)
		assert.Equal(t, val, "svanslös")
	})

	t.Run("Can delete item", func(t *testing.T) {
		cache.Delete("pelle")

		val, ok := cache.Get("pelle")
		assert.Equal(t, ok, false)
		assert.Equal(t, val, nil)
	})

	t.Run("Can clear cache", func(t *testing.T) {
		cache.Set("1", 1, 0)
		cache.Clear()

		val, ok := cache.Get("1")
		assert.Equal(t, ok, false)
		assert.Equal(t, val, nil)
	})
}
