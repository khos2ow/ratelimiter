package server

import (
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/khos2ow/ratelimiter/internal/version"
	"github.com/khos2ow/ratelimiter/pkg/ratelimiter"
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
	handler := &resourcesHandler{
		limiter: limiter,
	}

	router.HandleFunc("/", rootHandler)
	router.HandleFunc("/version", versionHandler)
	router.HandleFunc("/healthz", healthHandler)
	router.Handle("/{resource}", handler)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5000 * time.Millisecond,
		WriteTimeout: 5000 * time.Millisecond,
	}
	logrus.Info("Server is running http://localhost:8080 ...")
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, "OK: 'root' content.")
}
func versionHandler(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, version.Full())
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, "OK!")
}

func write(w http.ResponseWriter, status int, content string) int {
	w.WriteHeader(status)
	b, err := w.Write([]byte(content))
	if err != nil {
		logrus.Error(err)
	}
	return b
}
