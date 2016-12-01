package logging

import (
	"os"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

func addLoggingDefaults(cfg *viper.Viper) {
	cfg.SetDefault("level", "info")
	cfg.SetDefault("writer", map[interface{}]interface{}{"stderr": nil})
}

func parseLevel(level string) logrus.Level {
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		if os.Getenv("DEBUG") != "" {
			logrus.Infof("%v, falling back to default of error", err)
		}
		return logrus.ErrorLevel
	}
	return lvl
}

func mergeConfig(child, parent *viper.Viper) *viper.Viper {
	// This merge is only a partial merge
	// the remaining keys are not used to configure a logger but
	// indicate children of the current logger
	child.SetDefault("format", parent.GetString("format"))
	child.SetDefault("level", parent.GetString("level"))
	child.SetDefault("writer", parent.Get("writer"))

	// hooks are "special" they get merged for real
	// so if you define hooks then the hooks from the parent logger trickle down
	mergeHooks(child, parent)

	return child
}

func mergeFields(child, parent logrus.Fields) logrus.Fields {
	data := make(logrus.Fields, len(parent)+len(child))
	for k, v := range parent {
		data[k] = v
	}
	for k, v := range child {
		data[k] = v
	}
	return data
}

var loggerLock = new(sync.Mutex)

func configureLogger(logger *logrus.Logger, fields logrus.Fields, cfg *viper.Viper) {
	loggerLock.Lock()
	defer loggerLock.Unlock()
	logger.Level = parseLevel(cfg.GetString("level"))
	logger.Formatter = parseFormatter(cfg.GetString("format"), cfg)

	// writer config can be a string key or a full fledged config.
	var wcfg *viper.Viper
	if cfg.IsSet("writer") {
		vv := cfg.Get("writer")
		switch tpe := vv.(type) {
		case string:
			wcfg = viper.New()
			wcfg.Set("name", tpe)
		default:
			wcfg = cfg.Sub("writer")
		}
	}
	logger.Out = parseWriter(wcfg)

	logger.Hooks = make(logrus.LevelHooks, len(logrus.AllLevels))
	for _, hook := range parseHooks(cfg) {
		logger.Hooks.Add(hook)
	}
}

func newNamedLogger(name string, fields logrus.Fields, cfg *viper.Viper, parent *defaultLogger) *defaultLogger {

	logger := logrus.New()

	configureLogger(logger, fields, cfg)

	var bpth []string
	if parent != nil {
		bpth = parent.path
	}

	return &defaultLogger{
		Entry: logrus.Entry{
			Logger: logger,
			Data:   fields,
		},
		config: cfg,
		path:   append(bpth, strings.ToLower(name)),
	}
}

// Logger is the interface that application use to log against.
// Ideally you use the logrus.FieldLogger or logrus.StdLogger interfaces in your own code.
type Logger interface {
	logrus.FieldLogger
	New(string, logrus.Fields) Logger
	Configure(v *viper.Viper)
	Fields() logrus.Fields
}

type defaultLogger struct {
	logrus.Entry

	config *viper.Viper
	path   []string
	reg    *Registry
}

func (d *defaultLogger) New(name string, fields logrus.Fields) Logger {
	nme := strings.ToLower(name)
	pth := strings.Join(append(d.path, nme), ".")
	if l, ok := d.reg.GetOK(pth); ok {
		return l
	}

	data := mergeFields(fields, d.Entry.Data)
	data["module"] = name

	if d.config.InConfig(nme) {
		// new config, so make a new logger
		cfg := mergeConfig(d.config.Sub(nme), d.config)
		l := newNamedLogger(name, data, cfg, d)
		l.reg = d.reg
		d.reg.Register(pth, l)
		return l
	}

	// Share the logger with the parent, same config
	l := &defaultLogger{
		Entry: logrus.Entry{
			Logger: d.Entry.Logger,
			Data:   data,
		},
		config: d.config,
		path:   append(d.path, nme),
	}
	l.reg = d.reg
	d.reg.Register(pth, l)
	return l
}

func (d *defaultLogger) Configure(cfg *viper.Viper) {
	configureLogger(d.Logger, d.Data, cfg)
	d.config = cfg
}

func (d *defaultLogger) Fields() logrus.Fields {
	return d.Data
}
