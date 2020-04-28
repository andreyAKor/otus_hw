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
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		// [aaa] => [9, 8, 7, 6, 5]
		for i := 0; i < 10; i++ {
			_ = c.Set(Key(strconv.Itoa(i)), i)
		}

		val, ok := c.Get("9")
		require.True(t, ok)
		require.Equal(t, 9, val)

		val, ok = c.Get("5")
		require.True(t, ok)
		require.Equal(t, 5, val)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)

		// [9, 8, 7, 6, 5] => [6, 5, 9, 8, 7]
		for i := 5; i < 7; i++ {
			_, _ = c.Get(Key(strconv.Itoa(i)))
		}

		// [6, 5, 9, 8, 7] => [0, 1, 2, 6, 5]
		for i := 0; i < 3; i++ {
			_ = c.Set(Key(strconv.Itoa(i)), i)
		}

		val, ok = c.Get("0")
		require.True(t, ok)
		require.Equal(t, 0, val)

		val, ok = c.Get("5")
		require.True(t, ok)
		require.Equal(t, 5, val)
	})

	t.Run("additional logic", func(t *testing.T) {
		c := NewCache(5)

		// [99, 98, 97, 96, 95]
		for i := 0; i < 100; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		val, ok := c.Get("95")
		require.True(t, ok)
		require.Equal(t, 95, val)

		// [99, 98, 97, 96, 95] => []
		c.Clear()

		val, ok = c.Get("95")
		require.False(t, ok)
		require.Nil(t, val)

		for i := 100; i < 200; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}

		val, ok = c.Get("195")
		require.True(t, ok)
		require.Equal(t, 195, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
	//t.Skip() // Remove if task with asterisk completed

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
