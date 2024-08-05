package util

import (
	"errors"
	"sync"
	"time"
)

type SyncRequest struct {
	RequestName string `json:"RequestName"`
	CacheData   string `json:"CacheData,omitempty"`
	Attempts    int    `json:"Attempts,omitempty"`
}

type CachedSyncRequest struct {
	SyncRequest
	expiryTimestamp int64
}

type SyncCache struct {
	stop        chan struct{}
	wg          sync.WaitGroup
	mu          sync.RWMutex
	syncRequest map[string]CachedSyncRequest
}

func NewSyncCache(cleanupInterval time.Duration) *SyncCache {
	lc := &SyncCache{
		syncRequest: make(map[string]CachedSyncRequest),
		stop:        make(chan struct{}),
	}

	lc.wg.Add(1)
	go func(cleanupInterval time.Duration) {
		defer lc.wg.Done()
		lc.cleanupLoop(cleanupInterval)
	}(cleanupInterval)

	return lc
}

func (sc *SyncCache) cleanupLoop(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for {
		select {
		case <-sc.stop:
			return
		case <-t.C:
			sc.mu.Lock()
			for requestName, cu := range sc.syncRequest {
				if cu.expiryTimestamp <= time.Now().Unix() {
					delete(sc.syncRequest, requestName)
				}
			}
			sc.mu.Unlock()
		}
	}
}

// func (lc *SyncCache) stopCleanup() {
// 	close(lc.stop)
// 	lc.wg.Wait()
// }

func (sc *SyncCache) Update(sr SyncRequest, expiryTimestamp int64) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.syncRequest[sr.RequestName] = CachedSyncRequest{
		SyncRequest:     sr,
		expiryTimestamp: expiryTimestamp,
	}
}

var (
	errRequestNotInCache = errors.New("the request isn't in cache")
)

func (sc *SyncCache) Read(requestName string) (SyncRequest, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	sr, ok := sc.syncRequest[requestName]
	if !ok {
		return SyncRequest{}, errRequestNotInCache
	}

	return sr.SyncRequest, nil
}

func (sc *SyncCache) Delete(requestName string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	delete(sc.syncRequest, requestName)
}
