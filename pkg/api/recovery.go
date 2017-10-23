package api

import (
	"fmt"
	"github.com/Aptomi/aptomi/pkg/api/reqresp"
	"github.com/Aptomi/aptomi/pkg/lang/yaml"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"runtime/debug"
)

// NewPanicHandler returns HTTP handler for Panics processing
func NewPanicHandler(handler http.Handler) http.Handler {
	return &panicHandler{handler}
}

type panicHandler struct {
	handler http.Handler
}

func (h *panicHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)

			log.WithField("request", req).Errorf("Error while serving request: %s", err)

			if log.GetLevel() >= log.DebugLevel {
				log.Debug(string(debug.Stack()))
			}

			data := yaml.SerializeObject(reqresp.NewError(err))
			_, wErr := fmt.Fprint(w, data)
			if wErr != nil {
				log.Errorf("Error while writing error to response: %s", err)
			}
		}
	}()

	h.handler.ServeHTTP(w, req)
}
