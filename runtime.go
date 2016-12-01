package kvstore

import (
	"github.com/Sirupsen/logrus"
	app "github.com/casualjim/go-app"
	"github.com/casualjim/go-app/tracing"
	"github.com/go-openapi/kvstore/persist"
	"github.com/spf13/viper"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(app app.Application) (*Runtime, error) {
	db, err := persist.NewGoLevelDBStore(app.Config())
	if err != nil {
		return nil, err
	}
	return &Runtime{
		db: db,
	}, nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	db  persist.Store
	app app.Application
}

// DB returns the persistent store
func (r *Runtime) DB() persist.Store {
	return r.db
}

// Tracer returns the root tracer, this is typically the only one you need
func (r *Runtime) Tracer() tracing.Tracer {
	return r.app.Tracer()
}

// Logger gets the root logger for this application
func (r *Runtime) Logger() logrus.FieldLogger {
	return r.app.Logger()
}

// NewLogger creates a new named logger for this application
func (r *Runtime) NewLogger(name string, fields logrus.Fields) logrus.FieldLogger {
	return r.app.NewLogger(name, fields)
}

// Config returns the viper config for this application
func (r *Runtime) Config() *viper.Viper {
	return r.app.Config()
}
