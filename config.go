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
	ServerPort       int
}

// GlobalConfig configuration singleton
var GlobalConfig = RuntimeConfig{
	AuthzProvider:    "fence",
	FenceEndpointURL: "http://fence-service",
	ServerPort:       80,
}

// LoadEnv initialize the runtime configuration from
// environment variables: APROXY_FENCE_ENDPOINT
func (config *RuntimeConfig) LoadEnv() {
	if serverPortStr := os.Getenv("AUTHPROXY_PORT"); serverPortStr != "" {
		serverPort, err := strconv.Atoi(serverPortStr)
		if nil != err {
			config.ServerPort = serverPort
		}
	}
}

// String interface implementation
func (config RuntimeConfig) String() string {
	return fmt.Sprintf("(AuthzProvider: %v, FenceEndpointURL: %v, ServerPort: %v)",
		config.AuthzProvider, config.FenceEndpointURL, config.ServerPort,
	)
}
