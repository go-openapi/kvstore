package middlewares

import (
	"net/http"
	pprof "net/http/pprof"
	"strings"
)

// Profiler exposes net/http/pprof as a middleware
type Profiler struct {
	next http.Handler
}

// NewProfiler creates a middleware for profiling
func NewProfiler(next http.Handler) http.Handler {
	return &Profiler{next}
}

func (p *Profiler) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/debug/pprof") {
		p.next.ServeHTTP(rw, r)
		return
	}

	switch r.URL.Path {
	case "/debug/pprof/cmdline":
		pprof.Cmdline(rw, r)
		return
	case "/debug/pprof/profile":
		pprof.Profile(rw, r)
		return
	case "/debug/pprof/symbol":
		pprof.Symbol(rw, r)
		return
	default:
		pprof.Index(rw, r)
		return
	}
}
