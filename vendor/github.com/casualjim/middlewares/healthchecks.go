package middlewares

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/justinas/alice"
	"github.com/rcrowley/go-metrics"
)

// HealthChecks is a middleware that serves healthcheck information
type HealthChecks struct {
	basePath string
	next     http.Handler
}

// NewHealthChecksMW creates a new health check middleware at the specified path
func NewHealthChecksMW(basePath string) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return NewHealthChecks(basePath, next)
	}
}

// NewHealthChecks creates a new health check middleware at the specified path
func NewHealthChecks(basePath string, next http.Handler) *HealthChecks {
	if basePath == "" {
		basePath = "/"
	}

	return &HealthChecks{basePath: basePath, next: next}
}

// ServeHTTP is the middleware interface implementation
func (h *HealthChecks) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, filepath.Join(h.basePath, "health")) && !strings.HasPrefix(r.URL.Path, filepath.Join(h.basePath, "audit", "health")) {
		h.next.ServeHTTP(rw, r)
		return
	}

	metrics.RunHealthchecks()
	var errors []string
	metrics.DefaultRegistry.Each(func(name string, metric interface{}) {
		if hc, ok := metric.(metrics.Healthcheck); ok {
			if hc.Error() != nil {
				errors = append(errors, name+": failed, because: "+hc.Error().Error()+"\n")
			}
		}
	})

	if len(errors) > 0 {
		rw.WriteHeader(http.StatusInternalServerError)
		for _, err := range errors {
			rw.Write([]byte(err))
		}
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("OK"))
}
