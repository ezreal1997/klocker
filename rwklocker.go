package klocker

import (
	"fmt"
	"sync"
)

type RWKMutex struct {
	mu      sync.Mutex
	locks   map[string]*sync.RWMutex
	rCounts map[string]int
	wCounts map[string]int
}

func (rw *RWKMutex) lazyInit() {
	if rw.locks == nil {
		rw.locks = make(map[string]*sync.RWMutex)
		rw.rCounts = make(map[string]int)
		rw.wCounts = make(map[string]int)
	}
}

func (rw *RWKMutex) tryClean(key string) {
	if rw.rCounts[key] == 0 && rw.wCounts[key] == 0 {
		delete(rw.locks, key)
		delete(rw.rCounts, key)
		delete(rw.wCounts, key)
	}
}

func (rw *RWKMutex) Lock(key string) {
	rw.mu.Lock()
	rw.lazyInit()
	lock, ok := rw.locks[key]
	if !ok {
		lock = &sync.RWMutex{}
		rw.locks[key] = lock
	}
	rw.wCounts[key]++
	rw.mu.Unlock()
	lock.Lock()
}

func (rw *RWKMutex) Unlock(key string) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	lock, ok := rw.locks[key]
	if !ok || rw.wCounts[key] == 0 {
		panic(fmt.Sprintf("klocker: unlock unlocked rwkmutex of %#v", key))
	}
	lock.Unlock()
	rw.wCounts[key]--
	rw.tryClean(key)
}

func (rw *RWKMutex) RLock(key string) {
	rw.mu.Lock()
	rw.lazyInit()
	lock, ok := rw.locks[key]
	if !ok {
		lock = &sync.RWMutex{}
		rw.locks[key] = lock
	}
	rw.rCounts[key]++
	rw.mu.Unlock()
	lock.RLock()
}

func (rw *RWKMutex) RUnlock(key string) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	lock, ok := rw.locks[key]
	if !ok || rw.rCounts[key] == 0 {
		panic(fmt.Sprintf("klocker: runlock unlocked rwkmutex of %#v", key))
	}
	lock.RUnlock()
	rw.rCounts[key]--
	rw.tryClean(key)
}

type rKLocker RWKMutex

func (r *rKLocker) Lock(key string)   { (*RWKMutex)(r).RLock(key) }
func (r *rKLocker) Unlock(key string) { (*RWKMutex)(r).RUnlock(key) }

func (rw *RWKMutex) RKLocker() KLocker {
	return (*rKLocker)(rw)
}

func (rw *RWKMutex) Locker(key string) sync.Locker {
	return &locker{kl: rw, key: key}
}

func (rw *RWKMutex) RLocker(key string) sync.Locker {
	return &locker{kl: rw.RKLocker(), key: key}
}
