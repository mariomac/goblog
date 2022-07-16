package cache

import (
	"container/list"
)

// Sizable is any type whose size in bytes can be measured
type Sizable interface {
	SizeBytes() int
}

// LRU cache. It is not safe for concurrent access.
type LRU[K comparable, V Sizable] struct {
	sizeBytes    int
	maxSizeBytes int
	ll           *list.List
	cache        map[K]*list.Element
}

type entry[K comparable, V Sizable] struct {
	key   K
	value V
}

// NewLRU creates a new Cache that can store a maximum total bytes as the sum of the sizes
// of all the stored values.
func NewLRU[K comparable, V Sizable](maxSizeBytes int) *LRU[K, V] {
	return &LRU[K, V]{
		maxSizeBytes: maxSizeBytes,
		ll:           list.New(),
		cache:        map[K]*list.Element{},
	}
}

// Put a value into the cache.
func (c *LRU[K, V]) Put(key K, value V) {
	defer c.evictAll()
	c.sizeBytes += value.SizeBytes()
	if ee, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ee)
		c.sizeBytes -= ee.Value.(*entry[K, V]).value.SizeBytes()
		ee.Value.(*entry[K, V]).value = value
		return
	}
	ele := c.ll.PushFront(&entry[K, V]{key: key, value: value})
	c.cache[key] = ele
	// TODO: evict
}

// Get looks up a key's value from the cache.
func (c *LRU[K, V]) Get(key K) (value V, ok bool) {
	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		return ele.Value.(*entry[K, V]).value, true
	}
	return
}

// Remove the provided key from the cache.
func (c *LRU[K, V]) Remove(key K) {
	if ele, hit := c.cache[key]; hit {
		c.removeElement(ele)
	}
}

// RemoveOldest removes the oldest item from the cache.
func (c *LRU[K, V]) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *LRU[K, V]) removeElement(e *list.Element) {
	c.ll.Remove(e)
	kv := e.Value.(*entry[K, V])
	c.sizeBytes -= kv.value.SizeBytes()
	delete(c.cache, kv.key)
}

// evictAll removes the oldest entries until the cache reaches the given size.
func (c *LRU[K, V]) evictAll() {
	if c.sizeBytes <= c.maxSizeBytes {
		return
	}
	for c.sizeBytes > c.maxSizeBytes {
		c.removeOldest()
	}
}
