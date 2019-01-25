package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/uc-cdis/authproxy"
)

func main() {
	rxRun, _ := regexp.Compile(`^-*run$`)
	if len(os.Args) < 2 || !rxRun.MatchString(os.Args[1]) {
		fmt.Print(
			"Use: authProxyServer [--run]\n" +
				"  Launch the authproxy server with --run, otherwise print this help.\n" +
				"  Loads configuration from environment variables:\n" +
				"   * AUTHPROXY_PORT - port the	server listens on\n" +
				"   * AUTHPROXY_FENCE_ENDPOINT - URL of fence service\n" +
				"   * AUTHPROXY_CACHE_SIZE\n" +
				"   * AUTHPROXY_CACHE_TIMEOUT_SECS\n" +
				"Current configuration:\n" +
				fmt.Sprintf("%v\n", authproxy.GlobalConfig),
		)

		os.Exit(0)
	}
	var config = &authproxy.GlobalConfig
	config.LoadEnv()
	var authz = authproxy.AuthzFromRuntime(
		config,
		authproxy.PermCacheFromRuntime(config),
	)
	httpLogger := log.New(os.Stdout, "", log.LstdFlags)
	handler := authproxy.NewAuthzHandler(authz, config.ServerRoot, httpLogger)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.ServerPort),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     httpLogger,
		Handler:      handler,
	}
	httpLogger.Println(fmt.Sprintf("serving at %s with config %v", httpServer.Addr, config))
	httpLogger.Fatal(httpServer.ListenAndServe())
}
