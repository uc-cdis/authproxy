package authproxy

// AuthzService is a common interface for different Authz providers (arborist, fence, ...)
type AuthzService interface {
	// Check whether the token grants the caller permission on resource.
	// true access authorizes the user, false denies access
	CheckAccess(token, action, resource string) (access bool, err error)
}

var fenceSingleton *FenceAuthzService
var arboristSingleton AuthzService

//
// AuthzFromRuntime builds an implementation
// of the AuthzService that conforms to the given config -
// currently just have FenceAuthzService and ArboristAuthzService
//
func AuthzFromRuntime(config *RuntimeConfig, cache PermCache) AuthzService {
	if nil == fenceSingleton {
		fenceSingleton = FenceConfigFromRuntime(config).Build(cache)
	}
	return fenceSingleton
}
