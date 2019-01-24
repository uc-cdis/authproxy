package authproxy

import (
	"fmt"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var testURL = os.Getenv("AUTHPROXY_FENCE_URL")
	if testURL == "" {
		// avoid changing already set environment - can screw up other tests
		testURL = "https://fence.frickjack"
		os.Setenv("AUTHPROXY_FENCE_URL", testURL)
	}
	os.Setenv("AUTHPROXY_PORT", "6680")
	var config = RuntimeConfig{}
	config.LoadEnv()
	if config.ServerPort != 6680 {
		t.Error(fmt.Sprintf("config.LoadEnv did not configure port 6680 != %v", config.ServerPort))
	}
	if config.FenceEndpointURL != testURL {
		t.Error("config.LoadEnv did not configure fenceEndpointURL " + testURL + " != " + config.FenceEndpointURL)
	}
}

func TestEnvGet(t *testing.T) {
	result := EnvGet("TEST_BOGUS", "bogus")
	if result != "bogus" {
		t.Error(fmt.Sprintf("EnvGet did not give fallback %v != %v", "bogus", result))
	}
	os.Setenv("TEST_BOGUS", "frickjack")
	result = EnvGet("TEST_BOGUS", "bogus")
	if result != "frickjack" {
		t.Error(fmt.Sprintf("EnvGet did not get value %v != %v", "frickjack", result))
	}
	resultInt := EnvGetInt("TEST_BOGUS", 44)
	if resultInt != 44 {
		t.Error(fmt.Sprintf("EnvGetInt did not give default on bad value %v != %v", 44, resultInt))
	}
	os.Setenv("TEST_BOGUS", "321")
	resultInt = EnvGetInt("TEST_BOGUS", 44)
	if resultInt != 321 {
		t.Error(fmt.Sprintf("EnvGetInt did not give env value %v != %v", 321, resultInt))
	}
}
