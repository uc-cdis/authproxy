package authproxy

import (
	"fmt"
	"os"
	"strconv"
)

// RuntimeConfig captures the authproxy runtime configuration
type RuntimeConfig struct {
	AuthzProvider    string
	FenceEndpointURL string
	ServerRoot       string
	ServerPort       int
	CacheSize        int
	CacheTimeoutSecs int
}

// GlobalConfig configuration singleton
var GlobalConfig = RuntimeConfig{
	AuthzProvider:    "fence",
	FenceEndpointURL: "http://fence-service",
	ServerRoot:       "",
	ServerPort:       7780,
	CacheSize:        1000,
	CacheTimeoutSecs: 60,
}

// EnvGet environment key lookup with fallback value
func EnvGet(key, fallback string) string {
	if envStr := os.Getenv(key); envStr != "" {
		return envStr
	}
	return fallback
}

// EnvGetInt environment key lookup with convert to int with fallback value
func EnvGetInt(key string, fallback int) int {
	if envStr := os.Getenv(key); envStr != "" {
		envNum, err := strconv.Atoi(envStr)
		if nil == err {
			return envNum
		}
	}
	return fallback
}

// LoadEnv initialize the runtime configuration from
// environment variables: APROXY_FENCE_ENDPOINT
func (config *RuntimeConfig) LoadEnv() {
	config.ServerRoot = EnvGet("AUTHPROXY_ROOT", "")
	config.ServerPort = EnvGetInt("AUTHPROXY_PORT", config.ServerPort)
	config.FenceEndpointURL = EnvGet("AUTHPROXY_FENCE_URL", config.FenceEndpointURL)
	config.CacheSize = EnvGetInt("AUTHPROXY_CACHE_SIZE", config.CacheSize)
	config.CacheTimeoutSecs = EnvGetInt("AUTHPROXY_CACHE_TIMEOUT_SECS", config.CacheTimeoutSecs)
}

// String interface implementation
func (config RuntimeConfig) String() string {
	return fmt.Sprintf(
		"(AuthzProvider: %v, FenceEndpointURL: %v, ServerRoot %v, ServerPort: %v, CacheSize: %v, CacheTimeoutSecs: %v)",
		config.AuthzProvider, config.FenceEndpointURL, config.ServerRoot, config.ServerPort,
		config.CacheSize, config.CacheTimeoutSecs,
	)
}
