package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type cacheItem struct {
	Key   Key
	Value interface{}
}

type lruCache struct {
	capacity int
	lst      List
	indexes  map[Key]*Item
	locker   sync.RWMutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.locker.Lock()
	defer c.locker.Unlock()

	cacheValue := &cacheItem{key, value}
	if itm, ok := c.indexes[key]; ok {
		itm.Value = cacheValue
		c.lst.MoveToFront(itm)
		return true
	}
	c.indexes[key] = c.lst.PushFront(cacheValue)
	if c.lst.Len() > c.capacity {
		back := c.lst.Back()
		c.lst.Remove(back)
		delete(c.indexes, (back.Value.(*cacheItem)).Key)
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.locker.RLock()
	defer c.locker.RUnlock()

	if itm, ok := c.indexes[key]; ok {
		c.lst.MoveToFront(itm)
		return (itm.Value.(*cacheItem)).Value, ok
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.locker.Lock()
	defer c.locker.Unlock()

	c.lst.Clear()
	c.indexes = map[Key]*Item{}
}

func NewCache(capacity int) Cache {
	return &lruCache{capacity: capacity, lst: NewList(), indexes: map[Key]*Item{}}
}
