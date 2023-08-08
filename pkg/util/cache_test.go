package util

import (
	"testing"
	"time"
)

func TestRead(t *testing.T) {
	duration := time.Duration(1343432423)
	syncCache := NewSyncCache(duration)
	syncRequest := SyncRequest{RequestName: "test", Attempts: 2}
	expiry := int64(3432432426)
	syncCache.Update(syncRequest, expiry)
	sr, err := syncCache.Read("test")

	if err != nil {
		t.Errorf("Error looking up syncRequest")
	}
	if sr.RequestName != "test" {
		t.Errorf("Update failed")
	}
}

func TestUpdate(t *testing.T) {
	duration := time.Duration(1343432423)
	syncCache := NewSyncCache(duration)
	syncRequest := SyncRequest{RequestName: "test", Attempts: 2}
	expiry := int64(3432432426)
	syncCache.Update(syncRequest, expiry)
	cachedSyncRequest := syncCache.syncRequest["test"]

	if cachedSyncRequest.expiryTimestamp != expiry {
		t.Errorf("Expected timestamp %d, Actula timestatmp %d", expiry, cachedSyncRequest.expiryTimestamp)
	}
}

func TestDelete(t *testing.T) {
	duration := time.Duration(1343432423)
	syncCache := NewSyncCache(duration)
	syncRequest := SyncRequest{RequestName: "test", Attempts: 2}
	expiry := int64(3432432426)
	syncCache.Update(syncRequest, expiry)
	syncCache.Delete("test")
	_, ok := syncCache.syncRequest["test"]
	if ok {
		t.Errorf("test entry should not be present after delete operation")
	}
}
