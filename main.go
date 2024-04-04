package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/itzg/zapconfigs"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

func main() {

	proxyConfig := NewProxyConfig()

	var logger *zap.Logger
	if proxyConfig.Debug {
		logger = zapconfigs.NewDebugLogger()
	} else {
		logger = zapconfigs.NewDefaultLogger()
	}
	defer logger.Sync()

	targetURL, err := url.Parse(proxyConfig.BackendUrl)
	if err != nil {
		log.Fatal(err)
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	http.HandleFunc(proxyConfig.Endpoint, func(w http.ResponseWriter, r *http.Request) {

		r.Header.Set("X-WEBAUTH-USER", "kevcoxe")

		proxy.ServeHTTP(w, r)
	})

	http.HandleFunc(proxyConfig.HealthCheckEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			logger.
				With(zap.Error(err)).
				Error("failed to write health response body")
		}
	})

	http.Handle(proxyConfig.MetricsEndpoint, promhttp.Handler())

	log.Fatal(http.ListenAndServe(proxyConfig.Bind, nil))
}
