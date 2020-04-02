package hw04_lru_cache //nolint:golint,stylecheck

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(t *testing.T) {
		cache := NewCache(5)
		for i := 0; i < 5; i++ {
			cache.Set(Key(strconv.Itoa(i)), i)
		}
		// [4, 3, 2, 1, 0]

		zero, ok := cache.Get(Key("0")) // [0, 4, 3, 2, 1]
		require.True(t, ok)
		require.Equal(t, 0, zero)

		four, ok := cache.Get(Key("4")) // [4, 0, 3, 2, 1]
		require.True(t, ok)
		require.Equal(t, 4, four)

		ok = cache.Set(Key("5"), 5) // [5, 4, 0, 3, 2]
		require.False(t, ok)

		_, ok = cache.Get(Key("4")) // [4, 5, 0, 3, 2]
		require.True(t, ok)

		_, ok = cache.Get(Key("1")) // [4, 5, 0, 3, 2]
		require.False(t, ok)

		cache.Clear() // []
		for i := 0; i < 6; i++ {
			_, ok = cache.Get(Key(strconv.Itoa(i)))
			require.False(t, ok)
		}
	})
}

func TestCacheMultithreading(t *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
