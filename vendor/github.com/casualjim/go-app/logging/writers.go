package logging

import (
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

// CreateWriter based on a viper config
type CreateWriter func(*viper.Viper) io.Writer

var (
	knownWriters map[string]CreateWriter
	writersLock  *sync.Mutex

	// DefaultWriter is used as fallback when no other writer can be found.
	// This defaults to writing to stderr, but you can replace it to do something more useful
	DefaultWriter io.Writer
)

func init() {
	writersLock = new(sync.Mutex)
	knownWriters = make(map[string]CreateWriter, 50)
	DefaultWriter = os.Stderr

	knownWriters["discard"] = func(_ *viper.Viper) io.Writer {
		return ioutil.Discard
	}

	knownWriters["stdout"] = func(_ *viper.Viper) io.Writer {
		return os.Stdout
	}

	knownWriters["stderr"] = func(_ *viper.Viper) io.Writer {
		return os.Stderr
	}
}

// KnownWriters returns the list of keys for the registered writers
func KnownWriters() []string {
	writersLock.Lock()

	var writers []string
	for k := range knownWriters {
		writers = append(writers, k)
	}
	sort.Strings(writers)

	writersLock.Unlock()
	return writers
}

// RegisterWriter for use through the configuration system
// When you register a writer with a name that was already present
// then that writer will get overwritten
func RegisterWriter(name string, factory CreateWriter) {
	writersLock.Lock()
	knownWriters[strings.ToLower(name)] = factory
	writersLock.Unlock()
}

func parseWriter(cfg *viper.Viper) io.Writer {
	if cfg == nil || !cfg.IsSet("name") {
		return DefaultWriter
	}

	name := strings.ToLower(cfg.GetString("name"))
	writersLock.Lock()
	defer writersLock.Unlock()

	if create, ok := knownWriters[name]; ok {
		return create(cfg)
	}
	return DefaultWriter
}
