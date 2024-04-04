package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/itzg/go-flagsfiller"
)

type ProxyConfig struct {
	Bind                string `default:":8000" usage:"[host:port] to bind for serving HTTP"`
	BackendUrl          string `usage:"[URL] of the backend being proxied"`
	Endpoint            string `default:"/" usage:"[path] to proxy to the backend"`
	MetricsEndpoint     string `default:"/metrics" usage:"[path] to expose Prometheus metrics"`
	HealthCheckEndpoint string `default:"/_health" usage:"[path] to expose health check"`
	Debug               bool   `usage:"Enable debug logs"`
}

func NewProxyConfig() *ProxyConfig {
	var proxyConfig ProxyConfig

	filler := flagsfiller.New(flagsfiller.WithEnv("HEADER_PROXY"))
	err := filler.Fill(flag.CommandLine, &proxyConfig)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	checkRequired(proxyConfig.BackendUrl, "backend-url")

	return &proxyConfig
}

func checkRequired(value string, name string) {
	if value == "" {
		_, _ = fmt.Fprintf(os.Stderr, "%s is required\n", name)
		flag.Usage()
		os.Exit(2)
	}
}
