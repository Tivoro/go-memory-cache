package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type Cache struct {
	store map[string]CacheValue
	expQueue ExpirationQueue
	mutex sync.Mutex
}

type CacheValue struct {
	val any
	ttl int
	exp time.Time
}

type ExpirationQueue []*CacheItem
type CacheItem struct {
	index int
	key string
	exp time.Time
}

func NewCache() *Cache {
	cache := &Cache{
		store: make(map[string]CacheValue),
		expQueue: make(ExpirationQueue, 0),
	}
	heap.Init(&cache.expQueue)
	cache.watchHeapExpiration()

	return cache
}

func (cache *Cache) watchHeapExpiration() {
	go func() {
		for {
			if cache.expQueue.Len() == 0 {
				time.Sleep(100 * time.Millisecond)
				continue
			}

			wait := time.Until(cache.expQueue[0].exp)
			time.Sleep(wait)

			for {
				if cache.expQueue.Len() > 0 && cache.expQueue[0].exp.Before(time.Now()) {
					item := heap.Pop(&cache.expQueue).(*CacheItem)
					fmt.Println("Deleting item:", item.key, item.exp)
					cache.Delete(item.key)
				} else {
					break
				}
			}
		}
	}()
}

func (cache *Cache) Set(key string, val any, ttl int) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	cache.store[key] = CacheValue{
		val: val,
		ttl: ttl,
		exp: time.Now().Add(time.Duration(int(time.Millisecond) * ttl)),
	}
	if ttl > 0 {
		heap.Push(&cache.expQueue, &CacheItem{ key: key, exp: cache.store[key].exp })
	}
}

func (cache *Cache) Get(key string) (any, bool) {
	item, ok := cache.store[key]
	if !ok {
		return nil, false
	}
	if item.ttl > 0 && time.Now().After(item.exp) {
		return nil, false
	}
	return item.val, true
}

func (cache *Cache) Delete(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	delete(cache.store, key)
}	

func (cache *Cache) Clear() {
	for key := range cache.store {
		cache.Delete(key)
	}
}

func (eq *ExpirationQueue) Push(x any) {
	n := len(*eq)
	item := x.(*CacheItem)
	item.index = n
	*eq = append(*eq, item)
}

func (eq *ExpirationQueue) Pop() any {
	old := *eq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*eq = old[0:n-1]
	return item
}

func (eq ExpirationQueue) Len() int {
	return len(eq)
}

func (eq ExpirationQueue) Less(i, j int) bool {
	return eq[i].exp.Before(eq[j].exp)
}

func (eq ExpirationQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
	eq[i].index = j
	eq[j].index = i
}
