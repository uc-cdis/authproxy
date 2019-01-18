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

// Build a FenceAuthzService configured with this configuration
func (config *FenceConfig) Build() *FenceAuthzService {
	fence := &FenceAuthzService{Config: *config, Client: nil}
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	fence.Client = &http.Client{
		Transport: tr,
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
func (fence *FenceAuthzService) CheckAccess(token, permission, resource string) (access int, err error) {
	resp, err := http.Get(fence.Config.EndPointURL + "/user")
	if nil != err {
		return AUTHZ_NOTOK, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return AUTHZ_NOTOK, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if nil != err {
		return AUTHZ_NOTOK, err
	}
	userInfo := FenceUserResponseFromJSON(body)

	if nil != userInfo && userInfo.IsAdmin {
		return AUTHZ_OK, nil
	}
	return AUTHZ_NOTOK, nil
}
