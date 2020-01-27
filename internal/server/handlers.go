package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khos2ow/ratelimiter/internal/version"
	"github.com/sirupsen/logrus"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, "OK: 'root' content.")
}
func versionHandler(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, version.Full())
}
func healthHandler(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, "OK!")
}

func resourcesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resource := vars["resource"]
	// TODO do rate limit here
	// if limiter.IsAllowed() {
	if true {
		write(w, http.StatusOK, fmt.Sprintf("OK: '%s' content.", resource))
	} else {
		write(w, http.StatusTooManyRequests, fmt.Sprintf("Error: '%s' rate limited!", resource))
	}
}

func write(w http.ResponseWriter, status int, content string) int {
	w.WriteHeader(status)
	b, err := w.Write([]byte(content))
	if err != nil {
		logrus.Error(err)
	}
	return b
}
