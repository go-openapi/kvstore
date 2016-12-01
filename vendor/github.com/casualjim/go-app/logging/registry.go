package logging

import (
	"io"
	"sort"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// RootName of the root logger, defaults to root
	RootName string

	logKeys map[string]struct{}
)

func init() {
	RootName = "root"

	logKeys = map[string]struct{}{
		"level":  struct{}{},
		"format": struct{}{},
		"writer": struct{}{},
		"hooks":  struct{}{},
	}
}

// LoggerRegistry represents a registry for known loggers
type Registry struct {
	config *viper.Viper
	store  map[string]Logger
	lock   *sync.Mutex
}

// NewRegistry creates a new logger registry
func NewRegistry(cfg *viper.Viper, context logrus.Fields) *Registry {
	if cfg == nil {
		cfg = viper.New()
	}

	c := cfg
	if c.IsSet("logging") {
		c = cfg.Sub("logging")
	}

	var keys []string
	for _, kn := range c.AllKeys() {
		ln := strings.SplitN(kn, ".", 2)
		if _, ok := logKeys[ln[0]]; !ok {
			keys = append(keys, ln[0])
		}
	}

	store := make(map[string]Logger, len(keys))
	reg := &Registry{
		store:  store,
		config: c,
		lock:   new(sync.Mutex),
	}

	for _, k := range keys {

		// no sharing of context, so copy
		fields := make(logrus.Fields, len(context)+1)
		for kk, vv := range context {
			fields[kk] = vv
		}

		v := c
		if c.IsSet(k) {
			v = c.Sub(k)
		}

		addLoggingDefaults(v)
		fields["module"] = k
		if v.IsSet("name") {
			fields["module"] = v.GetString("name")
		}

		l := newNamedLogger(k, fields, v, nil)
		l.reg = reg
		reg.store[k] = l

	}

	if len(keys) == 0 {
		fields := make(logrus.Fields, len(context)+1)
		for k, v := range context {
			fields[k] = v
		}

		fields["module"] = RootName
		l := newNamedLogger(RootName, fields, c, nil)
		l.reg = reg
		reg.store[RootName] = l
	}

	return reg
}

// Get a logger by name, returns nil when logger doesn't exist.
// GetOK is the safe method to use.
func (r *Registry) Get(name string) Logger {
	l, ok := r.GetOK(name)
	if !ok {
		return nil
	}
	return l
}

// GetOK a logger by name, boolean is true when a logger was found
func (r *Registry) GetOK(name string) (Logger, bool) {
	r.lock.Lock()
	res, ok := r.store[strings.ToLower(name)]
	r.lock.Unlock()
	return res, ok
}

// Register a logger in this registry, overrides existing keys
func (r *Registry) Register(path string, logger Logger) {
	r.lock.Lock()
	r.store[strings.ToLower(path)] = logger
	r.lock.Unlock()
}

// Root returns the root logger, the name is configurable through the RootName variable
func (r *Registry) Root() Logger {
	return r.Get(RootName)
}

// Writer returns the pipe writer for the root logger
func (r *Registry) Writer() *io.PipeWriter {
	return r.Root().(*defaultLogger).Logger.Writer()
}

// Reload all the loggers with the new config
func (r *Registry) Reload() {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Get all keys, sorted by name and shortest to longest
	var keys []string
	for key := range r.store {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// find matching config
	// for each key find the longest possible path that has a config
	// if no more path found or parts are exhausted use last config and stop searching
	configs := make(map[string]*viper.Viper, len(keys))
	for _, key := range keys {
		configs[key] = findLongestMatchingPath(key, r.config)
	}

	// call reconfigure on logger
	for _, k := range keys {
		logger := r.store[k]
		if cfg, ok := configs[k]; ok {
			logger.Configure(cfg)
		}
	}
}

func findLongestMatchingPath(path string, cfg *viper.Viper) *viper.Viper {
	parts := strings.Split(path, ".")
	pl := len(parts)
	for i := range parts {
		mod := pl - i
		k := strings.Join(parts[:mod], ".")
		if cfg.IsSet(k) {
			return cfg.Sub(k)
		}
	}
	return nil
}
