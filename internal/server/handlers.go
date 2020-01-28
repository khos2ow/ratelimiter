package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khos2ow/ratelimiter"
	"github.com/sirupsen/logrus"
)

type resourcesHandler struct {
	limiter *ratelimiter.Limiter
}

func (rh *resourcesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resource := vars["resource"]
	logrus.Info("Registering new rate-limited resource: ", resource)
	if err := rh.limiter.Register(resource); err != nil {
		logrus.Error(err)
		write(w, http.StatusInternalServerError, "Interal Error!")
	}
	if rh.limiter.IsAllowed(resource) {
		write(w, http.StatusOK, fmt.Sprintf("OK: '%s' content.", resource))
	} else {
		write(w, http.StatusTooManyRequests, fmt.Sprintf("Error: '%s' rate limited!", resource))
	}
}
