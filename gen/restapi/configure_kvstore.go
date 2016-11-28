package restapi

import (
	"crypto/tls"
	"net/http"

	errors "github.com/go-openapi/errors"
	runtime "github.com/go-openapi/runtime"
	middleware "github.com/go-openapi/runtime/middleware"

	"github.com/go-openapi/kvstore/gen/restapi/operations"
	"github.com/go-openapi/kvstore/gen/restapi/operations/kv"
)

// This file is safe to edit. Once it exists it will not be overwritten

//go:generate swagger generate server --target ../gen --name kvstore --spec ../swagger/swagger.yml --exclude-main

func configureFlags(api *operations.KvstoreAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.KvstoreAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// s.api.Logger = log.Printf

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	api.KvDeleteEntryHandler = kv.DeleteEntryHandlerFunc(func(params kv.DeleteEntryParams) middleware.Responder {
		return middleware.NotImplemented("operation kv.DeleteEntry has not yet been implemented")
	})
	api.KvFindKeysHandler = kv.FindKeysHandlerFunc(func(params kv.FindKeysParams) middleware.Responder {
		return middleware.NotImplemented("operation kv.FindKeys has not yet been implemented")
	})
	api.KvGetEntryHandler = kv.GetEntryHandlerFunc(func(params kv.GetEntryParams) middleware.Responder {
		return middleware.NotImplemented("operation kv.GetEntry has not yet been implemented")
	})
	api.KvPutEntryHandler = kv.PutEntryHandlerFunc(func(params kv.PutEntryParams) middleware.Responder {
		return middleware.NotImplemented("operation kv.PutEntry has not yet been implemented")
	})

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
