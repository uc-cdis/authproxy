package authproxy

import (
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

//
// AuthzHandler implement http.Handler for auth requests
// Handles GET requests of form /root/action/resource
// Calls out to provided authzService, and returns 200 on success, 401 on fail
//
type AuthzHandler struct {
	Authz         AuthzService
	Root          string
	RxSingleSlash *regexp.Regexp
	RxTrimSlash   *regexp.Regexp
	Log           *log.Logger
}

// NewAuthzHandler init helper.  If nil log supplied, then
// creates default
func NewAuthzHandler(authz AuthzService, root string, logIn *log.Logger) *AuthzHandler {
	rxTrim, err := regexp.Compile(`(^/+)|(/+$)`)
	if nil != err {
		return nil
	}
	rxSingle, err := regexp.Compile(`/+`)
	if nil != err {
		return nil
	}
	result := &AuthzHandler{
		Authz:         authz,
		RxSingleSlash: rxSingle,
		RxTrimSlash:   rxTrim,
		Log:           logIn,
	}
	if nil == result.Log {
		result.Log = log.New(os.Stdout, "", log.LstdFlags)
	}
	result.Root = result.CleanPath(root)
	return result
}

// CleanPath trim leading and trailing / and dedup /
func (handler *AuthzHandler) CleanPath(in string) string {
	return string(handler.RxSingleSlash.ReplaceAll(handler.RxTrimSlash.ReplaceAll([]byte(in), []byte(``)), []byte(`/`)))
}

// ParsePath extracts action and resource from path of form handle.Root / action / resource,
// where action may not contain /
func (handler *AuthzHandler) ParsePath(in string) (action, resource string, ok bool) {
	clean := handler.CleanPath(in)
	if len(clean) < len(handler.Root)+2 || clean[:len(handler.Root)+1] != handler.Root+"/" {
		return "", "", false
	}
	clean = clean[len(handler.Root)+1:]
	slashIndex := strings.Index(clean, "/")
	if slashIndex < 1 || slashIndex == len(clean)-1 {
		return "", "", false
	}
	action = clean[:slashIndex]
	resource = clean[slashIndex+1:]
	return action, resource, true
}

func (handler *AuthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache") // response depends on auth header, so do not cache
	w.Header().Set("Content-Type", "application/json")

	if r.Method != "GET" && r.Method != "" { // only support GET
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "invalid http method"}`))
		return
	}

	prefix := "bearer "
	token := r.Header.Get("Authorization")
	if len(token) < len(prefix)+10 || token[:len(prefix)] != prefix {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`{"error": "invalid token"}`))
		return
	}

	action, resource, ok := handler.ParsePath(r.URL.Path)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad path"}`))
		return
	}
	if result, err := handler.Authz.CheckAccess(token, action, resource); nil != err {
		handler.Log.Printf("error accessing backend %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "server side"}`))
		return
	} else if !result {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "not authorized"}`))
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "go go go!"}`))
		return
	}
}
