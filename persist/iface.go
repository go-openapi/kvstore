package persist

import (
	"errors"
	"unsafe"

	"github.com/OneOfOne/xxhash"
)

// Common errors.
var (
	ErrNotFound         = errors.New("not found")
	ErrGone             = errors.New("gone")
	ErrReadOnly         = errors.New("read-only mode")
	ErrSnapshotReleased = errors.New("snapshot released")
	ErrIterReleased     = errors.New("iterator released")
	ErrClosed           = errors.New("closed")
	ErrVersionMismatch  = errors.New("version mismatch")
)

// UnsafeStringToBytes converts strings to []byte without memcopy
func UnsafeStringToBytes(s string) []byte {
	/* #nosec */
	return *(*[]byte)(unsafe.Pointer(&s))
}

// UnsafeBytesToString converts []byte to string without a memcopy
func UnsafeBytesToString(b []byte) string {
	/* #nosec */
	return *(*string)(unsafe.Pointer(&b))
}

// VersionOf calculates the version of the value to store
func VersionOf(data []byte) uint64 {
	h := xxhash.New64()
	_, _ = h.Write(data)
	return h.Sum64()
}

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
