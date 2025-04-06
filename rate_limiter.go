package ratelimiter

import (
	"sync"
	"time"
)

type Request struct {
	Time int64
}

type Cache struct {
	TimeFrame int64
	Limit     int
	Requests  map[string][]Request
	mu        sync.Mutex
}

func (cache *Cache) rateLimit(origin string) bool {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	reqs := cache.Requests[origin]
	if len(reqs) >= cache.Limit {
		return false
	} else {
		cache.Requests[origin] = append(reqs, Request{time.Now().UnixMilli()})
		return true
	}
}

func MakeRateLimiter(expireTime time.Duration, limit int) func(string) bool {

	newCache := Cache{
		TimeFrame: expireTime.Milliseconds(),
		Requests:  map[string][]Request{},
		mu:        sync.Mutex{},
		Limit:     limit,
	}

	go newCache.newCleaner(expireTime.Seconds())

	return newCache.rateLimit
}

func (c *Cache) newCleaner(interval float64) {
	for range time.Tick(time.Second * time.Duration(interval)) {
		c.cleanCache()
	}
}

func (cache *Cache) cleanCache() {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	now := time.Now().UnixMilli()
	for k, v := range cache.Requests {
		unexpired := []Request{}
		for _, r := range v {
			if r.Time+cache.TimeFrame > now {
				unexpired = append(unexpired, r)
			}
		}

		cache.Requests[k] = unexpired
	}
}
