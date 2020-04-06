package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/khos2ow/ratelimiter/pkg/ratelimiter"
	"github.com/sirupsen/logrus"
)

type resourcesHandler struct {
	limiter *ratelimiter.Limiter
}

func (rh *resourcesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	resource := vars["resource"]
	logrus.Info("checking rate-limit for resource: ", resource)
	ok, err := rh.limiter.IsAllowed(resource)
	if err != nil {
		write(w, http.StatusInternalServerError, err.Error())
		return
	}
	if ok {
		logrus.Info("serving resource: ", resource)
		write(w, http.StatusOK, fmt.Sprintf("OK: '%s' content.", resource))
	} else {
		logrus.Error("too many request for resource: ", resource)
		write(w, http.StatusTooManyRequests, fmt.Sprintf("Error: '%s' rate limited!", resource))
	}
}
