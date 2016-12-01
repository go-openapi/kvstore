package middlewares

import (
	"net/http"

	"github.com/justinas/alice"
)

// DefaultStack sets up the default middlewares
func DefaultStack(appInfo AppInfo, lgr Logger, orig http.Handler) http.Handler {
	return alice.New(
		NewRecoveryMW(appInfo.Name, lgr),
		NewAuditMW(appInfo, lgr),
		NewProfiler,
		NewHealthChecksMW(appInfo.BasePath),
		GzipMW(DefaultCompression),
	).Then(orig)
}
