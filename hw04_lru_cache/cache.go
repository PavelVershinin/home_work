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
	items    map[Key]cacheItem
	queue    List
}

type cacheItem struct {
	value interface{}
}

// Set Добавит значение в кеш, если оно там было обновит и вернёт true
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()

	_, ok := c.items[key]
	c.items[key] = cacheItem{value: value}

	// Добавим ключ в начало списка,
	// если элемент уже есть в списке, он будет перемещён в начало, потому что список теперь уникальный
	c.queue.PushFront(key)

	if !ok {
		// Грохнем элемент вышедший за пределы capacity, если такой есть
		for c.queue.Len() > c.capacity {
			i := c.queue.Back()
			c.queue.Remove(i)
			delete(c.items, i.Value.(Key))
		}
	}

	c.mu.Unlock()

	return ok
}

// Get Вернёт значение из кеша по ключу
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()

	item, ok := c.items[key]
	if ok {
		c.queue.PushFront(key)
	}

	c.mu.Unlock()

	return item.value, ok
}

// Clear Полностью очистит кеш
func (c *lruCache) Clear() {
	c.mu.Lock()
	c.items = make(map[Key]cacheItem)
	c.queue = NewList()
	c.mu.Unlock()
}

// NewCache Создаст новый кеш
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		items:    make(map[Key]cacheItem),
		queue:    NewList(),
	}
}
