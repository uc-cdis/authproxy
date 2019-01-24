package authproxy

// Valid results from AuthService.CheckAccess
const (
	AUTHZ_OK    = iota
	AUTHZ_NOTOK = iota
)

// AuthzService is a common interface for different Authz providers (arborist, fence, ...)
type AuthzService interface {
	// Check whether the token grants the caller permission on resource.
	CheckAccess(token, action, resource string) (access int, err error)
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
