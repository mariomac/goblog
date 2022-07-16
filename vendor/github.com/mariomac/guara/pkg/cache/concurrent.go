package cache

import "sync"

// Concurrent wrapper to add thread-safety to any cache
type Concurrent[K comparable, V any] struct {
	m     sync.RWMutex
	inner Cache[K, V]
}

func NewConcurrent[K comparable, V any](inner Cache[K, V]) *Concurrent[K, V] {
	return &Concurrent[K, V]{inner: inner}
}

func (c *Concurrent[K, V]) Put(key K, value V) {
	c.m.Lock()
	defer c.m.Unlock()
	c.inner.Put(key, value)
}

func (c *Concurrent[K, V]) Get(key K) (value V, ok bool) {
	c.m.RLock()
	defer c.m.RUnlock()
	return c.inner.Get(key)
}

func (c *Concurrent[K, V]) Remove(key K) {
	c.m.Lock()
	defer c.m.Unlock()
	c.inner.Remove(key)
}
