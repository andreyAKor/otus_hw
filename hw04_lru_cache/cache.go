package hw04_lru_cache //nolint:golint,stylecheck

import (
	"sync"
)

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key Key) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                              // Очистить кэш
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*Item
	mux      sync.Mutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

// Добавляет значение в кэш по ключу
func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mux.Lock()
	defer l.mux.Unlock()

	if i, ok := l.items[key]; ok {
		i.Value.(*cacheItem).value = value
		l.queue.MoveToFront(i)

		return true
	}

	i := l.queue.PushFront(&cacheItem{key, value})
	l.queue.MoveToFront(i)
	l.items[key] = i

	if l.queue.Len() > l.capacity {
		i := l.queue.Back()
		l.queue.Remove(i)

		delete(l.items, i.Value.(*cacheItem).key)
	}

	return false
}

// Возвращает значение из кэша по ключу
func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mux.Lock()
	defer l.mux.Unlock()

	if i, ok := l.items[key]; ok {
		l.queue.MoveToFront(i)

		return i.Value.(*cacheItem).value, true
	}

	return nil, false
}

// Очиститка кэша
func (l *lruCache) Clear() {
	l.mux.Lock()
	defer l.mux.Unlock()

	for key, i := range l.items {
		l.queue.Remove(i)
		delete(l.items, key)
	}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*Item),
	}
}
