package authproxy

import (
	"fmt"
	"os"
	"testing"
)

func TestFenceAuthzInvalidToken(t *testing.T) {
	var authz AuthzService = FenceConfigFromRuntime(&GlobalConfig).Build(NullCache)
	result, _ := authz.CheckAccess("bogustoken", "prometheus-admin", "*")
	if result {
		t.Error("bogus access check should fail")
	}
}

func TestFenceResponse(t *testing.T) {
	test1 := []byte(`{ "frick": "Frack", "is_admin": false, "whatever": "bla" }`)
	info := FenceUserResponseFromJSON(test1)
	if nil == info || info.IsAdmin {
		t.Error("json1 should not be admin")
	}
	test2 := []byte(`{ "frick": "Frack", "is_admin": true, "whatever": "bla" }`)
	info = FenceUserResponseFromJSON(test2)
	if nil == info || !info.IsAdmin {
		t.Error("json2 should be admin")
	}
}

// TestFenceAuthzValidToken set environment variable AUTHZ_TEST_TOKEN
func TestFenceAuthzValidToken(t *testing.T) {
	GlobalConfig.LoadEnv()
	var cache = PermCacheFromRuntime(&GlobalConfig)
	var authz AuthzService = FenceConfigFromRuntime(&GlobalConfig).Build(cache)
	var token = os.Getenv("AUTHPROXY_TEST_TOKEN")
	if len(token) < 10 {
		t.Skip("AUTHPROXY_TEST_TOKEN not set")
	}
	for i := 0; i < 10; i = i + 1 { // calls 2-9 should use cache
		result, err := authz.CheckAccess(token, "prometheus-admin", "*")
		if nil != err {
			t.Error(fmt.Sprintf("unexpected error: %v", err))
		}
		if !result {
			t.Error("test token access check should succeed")
		}
	}

	var ringCache *RingPCache = cache.(*RingPCache)
	if ringCache.Tail != 1 {
		t.Error(fmt.Sprintf("Expected cache to have one entry 1 != %v", ringCache.Tail))
	}

	for i := 0; i < 10; i = i + 1 { // calls 2-9 should use cache
		result, err := authz.CheckAccess(token+"bogus", "prometheus-admin", "*")
		if nil != err {
			t.Error(fmt.Sprintf("unexpected error with bad token: %v", err))
		}
		if result {
			t.Error("test token access check should not succeed")
		}
	}

	if ringCache.Tail != 2 {
		t.Error(fmt.Sprintf("Expected cache to have 2 entries 2 != %v", ringCache.Tail))
	}

}
