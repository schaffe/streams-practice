package main

import (
	"container/list"
	"sync"
	"testing"

	"github.com/go4org/hashtriemap"
)

// CacheSize determines how big the cache can grow
const CacheSize = 100

// KeyStoreCacheLoader is an interface for the KeyStoreCache
type KeyStoreCacheLoader interface {
	// Load implements a function where the cache should gets it's content from
	Load(string) string
}

type page struct {
	Key   string
	Value string
}

// KeyStoreCache is a LRU cache for string key-value pairs
type KeyStoreCache struct {
	mu    sync.Mutex
	cache hashtriemap.HashTrieMap[string, *list.Element]
	pages list.List
	size  uint
	load  func(string) string
}

// New creates a new KeyStoreCache
func New(load KeyStoreCacheLoader) *KeyStoreCache {
	return &KeyStoreCache{
		load:  load.Load,
		cache: hashtriemap.HashTrieMap[string, *list.Element]{},
	}
}

// Get gets the key from cache, loads it from the source if needed
func (k *KeyStoreCache) Get(key string) string {

	if e, ok := k.cache.Load(key); ok {
		k.mu.Lock()
		k.pages.MoveToFront(e)
		k.mu.Unlock()
		return e.Value.(page).Value
	}
	// Miss - load from database and save it in cache
	// if cache is full remove the least used item
	k.mu.Lock()
	p := page{key, k.load(key)}

	if k.size >= CacheSize {
		end := k.pages.Back()
		// remove from map
		k.cache.Delete(end.Value.(page).Key)
		// remove from list
		k.pages.Remove(end)
	} else {
		k.size++
	}
	k.pages.PushFront(p)
	k.cache.Swap(key, k.pages.Front())
	k.mu.Unlock()
	return p.Value
}

// Loader implements KeyStoreLoader
type Loader struct {
	DB *MockDB
}

// Load gets the data from the database
func (l *Loader) Load(key string) string {
	val, err := l.DB.Get(key)
	if err != nil {
		panic(err)
	}

	return val
}

func run(t *testing.T) (*KeyStoreCache, *MockDB) {
	loader := Loader{
		DB: GetMockDB(),
	}
	cache := New(&loader)

	RunMockServer(cache, t)

	return cache, loader.DB
}

func main() {
	run(nil)
}
