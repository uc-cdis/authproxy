package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/uc-cdis/authproxy"
)

func main() {
	authproxy.GlobalConfig.LoadEnv()
	rxRun, _ := regexp.Compile(`^-*run$`)
	if len(os.Args) < 2 || !rxRun.MatchString(os.Args[1]) {
		fmt.Print(
			"Use: authProxyServer [--run]\n" +
				"  Launch the authproxy server with --run, otherwise print this help.\n" +
				"  Loads configuration from environment variables:\n" +
				"   * AUTHPROXY_PORT - port the	server listens on\n" +
				"   * AUTHPROXY_FENCE_ENDPOINT - URL of fence service\n" +
				"Current configuration:\n" +
				fmt.Sprintf("%v\n", authproxy.GlobalConfig),
		)

		os.Exit(0)
	}
	authproxy.AuthzFromRuntime(&authproxy.GlobalConfig)
	fmt.Println("ERROR: server run not yet implemented")
	os.Exit(1)
}
