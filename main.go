package main

import (
	"header-rev-proxy/handlers"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/go-chi/chi"
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

	r := chi.NewRouter()
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	r.HandleFunc(proxyConfig.Endpoint, func(w http.ResponseWriter, r *http.Request) {
		// check for cookie header-rev-proxy-username if it exists, set X-WEBAUTH-USER header else redirect to /login
		cookie, err := r.Cookie("header-rev-proxy-username")
		if err != nil {
			logger.Error("failed to get cookie")
			r.Header.Set("X-WEBAUTH-USER", "unknown")
		} else {
			cookieValue := cookie.Value
			r.Header.Set("X-WEBAUTH-USER", cookieValue)
		}

		proxy.ServeHTTP(w, r)
	})

	r.HandleFunc(proxyConfig.HealthCheckEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(200)
		_, err := w.Write([]byte("OK"))
		if err != nil {
			logger.
				With(zap.Error(err)).
				Error("failed to write health response body")
		}
	})

	r.Handle(proxyConfig.MetricsEndpoint, promhttp.Handler())

	r.Get("/", handlers.HomeGetHandler)
	r.Post("/login", handlers.LoginPostHandler)
	log.Printf("Listening on Port %s\n", proxyConfig.Bind)
	err = http.ListenAndServe(proxyConfig.Bind, r)
	if err != nil {
		log.Fatal(err)
	}
}
