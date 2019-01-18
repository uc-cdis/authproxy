package authproxy

import (
	"testing"
)

func TestFenceAuthzInvalidToken(t *testing.T) {
	var authz AuthzService = FenceConfigFromRuntime(&GlobalConfig).Build()
	result, _ := authz.CheckAccess("bogustoken", "prometheus-admin", "*")
	if result != AUTHZ_NOTOK {
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

func TestFenceAuthzValidToken(t *testing.T) {

}
