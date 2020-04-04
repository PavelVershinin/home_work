package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

// Cache Интерфейс кеша
type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	mu       sync.Mutex
	capacity int
	items    map[Key]*listItem
	queue    List
}

type cacheItem struct {
	key   Key
	value interface{}
}

// Set Добавит значение в кеш, если оно там было обновит и вернёт true
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if !ok {
		item = c.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
	} else {
		item.Value = cacheItem{
			key:   key,
			value: value,
		}
		c.queue.MoveToFront(item)
	}
	c.items[key] = item

	if !ok {
		// Грохнем элемент вышедший за пределы capacity, если такой есть
		for c.queue.Len() > c.capacity {
			i := c.queue.Back()
			c.queue.Remove(i)
			delete(c.items, i.Value.(cacheItem).key)
		}
	}

	return ok
}

// Get Вернёт значение из кеша по ключу
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item)
	} else {
		return nil, ok
	}

	return item.Value.(cacheItem).value, ok
}

// Clear Полностью очистит кеш
func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[Key]*listItem)
	c.queue = NewList()
}

// NewCache Создаст новый кеш
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[Key]*listItem),
		queue:    NewList(),
	}
}
