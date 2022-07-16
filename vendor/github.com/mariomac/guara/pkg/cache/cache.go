package cache

type Cache[K comparable, V any] interface {
	// Put a value into the cache.
	Put(key K, value V)
	// Get looks up a key's value from the cache.
	Get(key K) (value V, ok bool)
	// Remove the provided key from the cache.
	Remove(key K)
}
