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
	*listItem
	value interface{}
}

// Set Добавит значение в кеш, если оно там было обновит и вернёт true
func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()

	item, ok := c.items[key]

	if !ok {
		if c.queue.Len() == c.capacity {
			last := c.queue.Back()            // Находим крайний элемент списка
			c.queue.Remove(last)              // Удаляем его из очереди
			delete(c.items, last.Value.(Key)) // Удаляем его из мапки
		}
		item.listItem = c.queue.PushFront(key)
	} else {
		c.queue.MoveToFront(item.listItem)
	}

	item.value = value
	c.items[key] = item

	c.mu.Unlock()

	return ok
}

// Get Вернёт значение из кеша по ключу
func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()

	item, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(item.listItem)
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
