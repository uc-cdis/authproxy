package authproxy

import (
	"sync"
	"time"
)

// RPCacheEntry in a PermissioCache
type RPCacheEntry struct {
	Token      string
	Action     string
	Resource   string
	HasAccess  bool
	CreateTime int64
}

// RingPCache with up to 1000 entries,
// multiple readers, but a single writer
// on top of a ring buffer (?)
type RingPCache struct {
	Tail        int
	Lock        sync.RWMutex
	Entries     []RPCacheEntry
	TimeoutSecs int
}

// NewRingPCache implements PermissionCaches
func NewRingPCache(cacheSize, timeoutSecs int) *RingPCache {
	if cacheSize < 1000 { // 1000 is min
		cacheSize = 1000
	}
	if timeoutSecs < 10 {
		timeoutSecs = 10
	}
	return &RingPCache{
		Tail:        0,
		Entries:     make([]RPCacheEntry, cacheSize, cacheSize),
		TimeoutSecs: timeoutSecs,
	}
}

// _lookup helper returns pointer to matching entry or nil
// Assumes caller acquires lock
func (cache *RingPCache) _lookup(token, action, resource string) *RPCacheEntry {
	// All lookup keys are required
	if token == "" || action == "" || resource == "" {
		return nil
	}
	var tTooOld = time.Now().Unix() - int64(cache.TimeoutSecs)
	for loopCount := 0; loopCount < len(cache.Entries); loopCount = loopCount + 1 {
		i := (cache.Tail - 1 - loopCount + len(cache.Entries)) % len(cache.Entries)
		if entry := &cache.Entries[i]; entry.CreateTime > tTooOld && len(entry.Token) == len(token) && entry.Token == token && entry.Action == action && entry.Resource == resource {
			return entry
		}
	}
	return nil
}

// Lookup the cache entry if any that matches the given token, action, and resource,
// and was created within the last 60 seconds
// Is hash map faster for 1000 entries?  Probably ...
func (cache *RingPCache) Lookup(token, action, resource string) (result bool, ok bool) {
	cache.Lock.RLock()
	defer cache.Lock.RUnlock()
	entryPtr := cache._lookup(token, action, resource)
	if nil != entryPtr {
		return entryPtr.HasAccess, true
	}
	return false, false
}

// Add a new entry to the cache.
func (cache *RingPCache) Add(token, action, resource string, hasAccess bool) {
	cache.Lock.Lock()
	defer cache.Lock.Unlock()
	// avoid duplicate entries
	entryPtr := cache._lookup(token, action, resource)
	if nil == entryPtr { // no duplicate entry, so overwrite the Tail
		entryPtr = &cache.Entries[cache.Tail]
		cache.Tail = (cache.Tail + 1) % len(cache.Entries)
		entryPtr.Token = token
		entryPtr.Action = action
		entryPtr.Resource = resource
	}
	entryPtr.HasAccess = hasAccess
	entryPtr.CreateTime = time.Now().Unix()
}
