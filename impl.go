package localcache

import (
	"sync"
	"time"
)

const defaultExpiredTime = 30 * time.Second

type localCache struct {
	data map[string]cacheData
	m    sync.RWMutex
}

type cacheData struct {
	value interface{}
	wiper *time.Timer
}

// New will return a localcache instance
func New() Cache {
	inst := &localCache{
		data: make(map[string]cacheData),
		m:    sync.RWMutex{},
	}

	return inst
}

// Get returns exist content from local cache
func (lc *localCache) Get(k string) (ret interface{}, e error) {
	lc.m.RLock()
	defer lc.m.RUnlock()

	if d, ok := lc.data[k]; ok {
		ret = d.value
		return
	}

	e = ErrKeyNonExist
	return
}

// Set stores value in local cache with key k
func (lc *localCache) Set(k string, v interface{}) error {
	lc.m.Lock()
	defer lc.m.Unlock()

	// Avoid older timer to evict
	if d, ok := lc.data[k]; ok {
		d.wiper.Stop()
	}

	lc.data[k] = cacheData{
		value: v,
		wiper: time.AfterFunc(defaultExpiredTime, func() {
			lc.evict(k)
		}),
	}

	return nil
}

func (lc *localCache) evict(k string) {
	lc.m.Lock()
	defer lc.m.Unlock()
	delete(lc.data, k)
}
