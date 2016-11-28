package kvstore

import (
	"github.com/go-openapi/kvstore/persist"
	"github.com/spf13/viper"
)

// NewRuntime creates a new application level runtime that encapsulates the shared services for this application
func NewRuntime(config *viper.Viper) (*Runtime, error) {
	db, err := persist.NewGoLevelDBStore(config)
	if err != nil {
		return nil, err
	}
	return &Runtime{
		db: db,
	}, nil
}

// Runtime encapsulates the shared services for this application
type Runtime struct {
	db persist.Store
}

// DB returns the persistent store
func (r *Runtime) DB() persist.Store {
	return r.db
}
