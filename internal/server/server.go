package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/khos2ow/ratelimiter"
	"github.com/sirupsen/logrus"
)

// Start starts internal http server and processes the requests.
// If `options.BackendServers` is passed it will act as a proxy
// and forwards requests to the backend server(s).
func Start(backends []string, limiter *ratelimiter.Limiter) error {
	logrus.Info("Starting server...")
	if len(backends) > 0 {
		// TODO add LB functionality against provided list of `backends`
		logrus.Warn("[Not implemented] Proxy requests to backend servers: ", strings.Join(backends, ", "))
	}

	router := mux.NewRouter()

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthHandler)
	router.HandleFunc("/{resource}", resourcesHandler)

	server := &http.Server{
		Addr:         ":8000",
		Handler:      router,
		ReadTimeout:  5000 * time.Millisecond,
		WriteTimeout: 5000 * time.Millisecond,
	}
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}
