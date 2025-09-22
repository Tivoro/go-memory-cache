package memcache

import "time"

type ExpirationQueue []*CacheItem

type CacheItem struct {
	index int
	key string
	exp time.Time
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
