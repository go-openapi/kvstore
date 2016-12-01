package middlewares

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/justinas/alice"
)

// Recovery is a middleware that recovers from any panics and writes a 500 if there was one.
type Recovery struct {
	Logger     Logger
	PrintStack bool
	next       http.Handler
}

// NewRecoveryMW returns a new instance of Recovery middleware which traps panics
func NewRecoveryMW(appName string, lgr Logger) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return NewRecovery(appName, lgr, next)
	}
}

// NewRecovery returns a new instance of Recovery
func NewRecovery(appName string, lgr Logger, next http.Handler) *Recovery {
	if appName == "" {
		appName = "api"
	}
	return &Recovery{
		Logger:     lgr,
		PrintStack: true,
		next:       next,
	}
}

func (rec *Recovery) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			stack := debug.Stack()

			rec.Logger.Errorf("%s\n%s", err, stack)

			if rec.PrintStack {
				fmtStr := `{"message":"%s","stack":"%s","type":"error"}`
				fmt.Fprintf(rw, fmtStr, err, stack)
			}
		}
	}()

	rec.next.ServeHTTP(rw, r)
}
