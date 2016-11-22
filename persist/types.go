package persist

import "errors"

// Common errors.
var (
	ErrNotFound         = errors.New("not found")
	ErrReadOnly         = errors.New("read-only mode")
	ErrSnapshotReleased = errors.New("snapshot released")
	ErrIterReleased     = errors.New("iterator released")
	ErrClosed           = errors.New("closed")
)

// Value returned from a persistence medium
type Value []byte

// KeyValue represents an entry with key name
type KeyValue struct {
	Key   string
	Value Value
	_     struct{}
}

// Store for values by key
type Store interface {
	Put(string, Value) error
	Get(string) (Value, error)
	FindByPrefix(string) ([]KeyValue, error)
	Delete(string) error
	Close() error
}
