package handlers

import (
	"fmt"
	"net/http"

	"github.com/emojify-app/face-detection/logging"
)

// Health is a HTTP Handler
type Health struct {
	l logging.Logger
}

// NewHealth creates a new Health handler
func NewHealth(l logging.Logger) *Health {
	return &Health{l}
}

// ServeHTTP is a http handler which reports the service health
func (h *Health) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	done := h.l.HealthHandlerCalled()
	defer done()

	fmt.Fprint(rw, "ok")
}
