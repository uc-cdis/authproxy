package authproxy

import (
	"fmt"
	"testing"
)

func TestCleanPath(t *testing.T) {
	prefix := "/frickjack//bla/frickjack//"
	handler := NewAuthzHandler(nil, "/frickjack//bla/frickjack//", nil)
	if nil == handler.Log {
		t.Error("Log not initialized")
	}
	if "frickjack/bla/frickjack" != handler.Root {
		t.Errorf("bad root path %v", handler.Root)
	}
	if _, _, ok := handler.ParsePath(prefix); ok {
		t.Errorf("should have failed parse of path %v", prefix)
	}
	for i := 0; i < 10; i = i + 1 {
		action := fmt.Sprintf("action%v", i)
		resource := "resource" + action + "/*/bla/" + action
		path := action + "/" + resource
		if _, _, ok := handler.ParsePath(path); ok {
			t.Errorf("should have failed to parse %v", path)
		}
		path = handler.Root + "/" + path + "//"
		if resultAction, resultResource, ok := handler.ParsePath(path); !(ok && action == resultAction && resource == resultResource) {
			t.Errorf("should have parsed %v != %v/%v/%v", handler.Root, resultAction, resultResource, path)
		}
	}
}
