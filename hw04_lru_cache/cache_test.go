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

		beginKey := 100_000_000
		endKey := beginKey + 6_000
		count := (endKey - beginKey) / 2
		middleKey := beginKey + count

		c := NewCache(count)

		for i := beginKey; i < endKey; i++ {
			require.False(t, c.Set(Key(strconv.Itoa(i)), i))
		}
		for i := beginKey; i < middleKey; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.False(t, ok)
			require.Nil(t, val)
		}
		for i := middleKey; i < endKey; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.True(t, ok)
			require.Equal(t, val, i)
		}
		c.Clear()
		for i := middleKey; i < endKey; i++ {
			val, ok := c.Get(Key(strconv.Itoa(i)))
			require.False(t, ok)
			require.Nil(t, val, i)
		}
	})

	t.Run("duplicates", func(t *testing.T) {
		c := NewCache(1)

		var key Key = "duplicate"
		val1 := "duplicate_VALUE"
		val2 := "NEW_duplicate_VALUE"

		require.False(t, c.Set(key, val1))
		require.True(t, c.Set(key, val2))

		val, ok := c.Get(key)
		require.True(t, ok)
		require.NotNil(t, val)
		require.NotEqual(t, val, val1)
		require.Equal(t, val, val2)
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
