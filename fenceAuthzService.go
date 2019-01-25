package authproxy

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

// FenceConfig configuration for FenceAuthzService
// TODO - setup a ConfigManager to centralize configuration
type FenceConfig struct {
	EndPointURL string
}

// MakeFenceConfig with default values
func MakeFenceConfig() *FenceConfig {
	return &FenceConfig{"http://fence-service"}
}

// Build a FenceAuthzService configured with this configuration,
// and the given cache
func (config *FenceConfig) Build(cache PermCache) *FenceAuthzService {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	fence := &FenceAuthzService{
		Config: *config,
		Client: &http.Client{
			Transport: tr,
		},
		Cache: cache,
	}

	return fence
}

// FenceConfigFromRuntime extracts FenceConfig from RuntimeConfig
func FenceConfigFromRuntime(rtConfig *RuntimeConfig) *FenceConfig {
	return &FenceConfig{rtConfig.FenceEndpointURL}
}

// FenceAuthzService implements AuthzService interface
type FenceAuthzService struct {
	Config FenceConfig
	Client *http.Client
	Cache  PermCache
	// Cache tokens for up to 1 minute
}

// FenceUserResponse from fence/user endpoint
type FenceUserResponse struct {
	IsAdmin bool `json:"is_admin"`
}

// FenceUserResponseFromJSON returns null of jsonStr decode fails
func FenceUserResponseFromJSON(jsonIn []byte) *FenceUserResponse {
	var info FenceUserResponse
	err := json.Unmarshal(jsonIn, &info)
	if nil != err {
		return nil
	}
	return &info
}

// CheckAccess implements limitted admin access control
func (fence *FenceAuthzService) CheckAccess(token, action, resource string) (access bool, err error) {
	if cacheAccess, cacheOk := fence.Cache.Lookup(token, action, resource); cacheOk {
		return cacheAccess, nil
	}
	req, err := http.NewRequest("GET", fence.Config.EndPointURL+"/user", nil)
	if nil != err {
		return false, err
	}
	req.Header.Set("Authorization", "bearer "+token)
	resp, err := fence.Client.Do(req)
	if nil != err {
		return false, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		fence.Cache.Add(token, action, resource, false)
		return false, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return false, err
	}
	userInfo := FenceUserResponseFromJSON(body)

	if nil != userInfo && userInfo.IsAdmin {
		fence.Cache.Add(token, action, resource, true)
		return true, nil
	}
	fence.Cache.Add(token, action, resource, false)
	return false, nil
}
