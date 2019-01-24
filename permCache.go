package authproxy

// PermCache standard interface for cacheing permission data
type PermCache interface {
	Lookup(token, action, resource string) (result bool, ok bool)
	Add(token, action, resource string, hasAccess bool)
}

// MockCache NOOP implementation of PermCache
type MockCache struct {
}

// Add NOOP for MockCache
func (cache *MockCache) Add(token, action, resource string, hasAccess bool) {
}

// Lookup for MockCache always returns false, false
func (cache *MockCache) Lookup(token, action, resource string) (result bool, ok bool) {
	return false, false
}

// NullCache global mock cache
var NullCache = &MockCache{}

var cacheSingleton PermCache

//
// AuthzFromRuntime builds an implementation
// of the AuthzService that conforms to the given config -
// currently just have FenceAuthzService and ArboristAuthzService
//
func PermCacheFromRuntime(config *RuntimeConfig) PermCache {
	if nil == cacheSingleton {
		cacheSingleton = NewRingPCache(config.CacheSize, config.CacheTimeoutSecs)
	}
	return cacheSingleton
}
