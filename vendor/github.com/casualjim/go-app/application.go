package app

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"sync"
	"syscall"

	"github.com/Sirupsen/logrus"
	"github.com/casualjim/go-app/logging"
	"github.com/casualjim/go-app/tracing"
	cjm "github.com/casualjim/middlewares"
	"github.com/fsnotify/fsnotify"
	"github.com/kardianos/osext"
	"github.com/spf13/viper"

	// we enable remote config providers by default
	_ "github.com/spf13/viper/remote"
)

var (
	// ErrModuleUnknown returned when no module can be found for the specified key
	ErrModuleUnknown error

	execName func() (string, error)

	// Version of the application
	Version string
)

func init() {
	ErrModuleUnknown = errors.New("unknown module")
	execName = osext.Executable
	log.SetOutput(logrus.StandardLogger().Writer())
	log.SetFlags(0)
}

// A Key represents a key for a module.
// Users of this package can define their own keys, this is just the type definition.
type Key string

// Application is an application level context package
// It can be used as a kind of dependency injection container
type Application interface {
	// Add modules to the application context
	Add(...Module) error

	// Get the module at the specified key, thread-safe
	Get(Key) interface{}

	// Get the module at the specified key, thread-safe
	GetOK(Key) (interface{}, bool)

	// Set the module at the specified key, this should be safe across multiple threads
	Set(Key, interface{}) error

	// Logger gets the root logger for this application
	Logger() logrus.FieldLogger

	// NewLogger creates a new named logger for this application
	NewLogger(string, logrus.Fields) logrus.FieldLogger

	// Tracer returns the root
	Tracer() tracing.Tracer

	// Config returns the viper config for this application
	Config() *viper.Viper

	// Info returns the app info object for this application
	Info() cjm.AppInfo

	// Init the application and its modules with the config.
	Init() error

	// Start the application an its enabled modules
	Start() error

	// Stop the application an its enabled modules
	Stop() error
}

func addDefaultConfigPaths(v *viper.Viper, name string) {
	v.SetConfigName("config")
	norm := strings.ToLower(name)
	paths := filepath.Join(os.Getenv("HOME"), ".config", norm) + ":" + filepath.Join("/etc", norm) + ":etc:."
	if os.Getenv("CONFIG_PATH") != "" {
		paths = os.Getenv("CONFIG_PATH")
	}
	for _, path := range filepath.SplitList(paths) {
		v.AddConfigPath(path)
	}
}

var viperLock *sync.Mutex

func init() {
	viperLock = new(sync.Mutex)
}

func createViper(name string, configPath string) (*viper.Viper, error) {
	viperLock.Lock()
	defer viperLock.Unlock()
	v := viper.New()
	if configPath == "" {
		addDefaultConfigPaths(v, name)
	} else {
		if _, err := os.Stat(configPath); os.IsNotExist(err) {
			return nil, fmt.Errorf("No config file found at %s", configPath)
		}
		dir, fname := filepath.Split(configPath)
		// viper wants the file name without extention...
		suffixes := []string{".json", ".yml", ".yaml", ".hcl", ".toml"}
		for _, suffix := range suffixes {
			if strings.HasSuffix(fname, suffix) {
				fname = strings.Split(fname, suffix)[0]
				break
			}
		}
		v.SetConfigName(fname)
		v.AddConfigPath(dir)
	}

	if err := addViperRemoteConfig(v); err != nil {
		return nil, err
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.UnsupportedConfigError); !ok && v.ConfigFileUsed() != "" {
			return nil, err
		}
	}
	v.SetEnvPrefix(name)
	v.AutomaticEnv()

	addViperDefaults(v)

	return v, nil
}

func addViperRemoteConfig(v *viper.Viper) error {
	// check if encryption is required CONFIG_KEYRING
	keyring := os.Getenv("CONFIG_KEYRING")

	// check for etcd env vars CONFIG_REMOTE_URL, eg:
	// etcd://localhost:2379/[app-name]/config.[type]
	// consul://localhost:8500/[app-name]/config.[type]
	remURL := os.Getenv("CONFIG_REMOTE_URL")
	if remURL == "" {
		return nil
	}
	var provider, host, path, tpe string
	u, err := url.Parse(remURL)
	if err != nil {
		return err
	}
	provider = strings.ToLower(u.Scheme)
	host = u.Host
	if provider == "etcd" {
		host = "http://" + host
	}
	path = u.Path
	tpe = strings.ToLower(strings.TrimLeft(filepath.Ext(path), "."))
	if tpe == "" {
		tpe = "json"
	}

	v.SetConfigType(tpe)
	if keyring != "" {
		if err := v.AddSecureRemoteProvider(provider, host, path, keyring); err != nil {
			return err
		}
	} else {

		if err := v.AddRemoteProvider(provider, host, path); err != nil {
			return err
		}
	}

	if err := v.ReadRemoteConfig(); err != nil {
		return fmt.Errorf("config is invalid as %s", tpe)
	}

	return nil
}

func addViperDefaults(v *viper.Viper) {
	v.SetDefault("tracer", map[interface{}]interface{}{"enable": true})
	v.SetDefault("logging", map[interface{}]interface{}{"root": map[interface{}]interface{}{"level": "info"}})
}

func ensureDefaults(name string) (string, string, error) {
	// configure version defaults
	version := "dev"
	if Version != "" {
		version = Version
	}

	// configure name defaults
	if name == "" {
		name = os.Getenv("APP_NAME")
		if name == "" {
			exe, err := execName()
			if err != nil {
				return "", "", err
			}
			name = filepath.Base(exe)
		}
	}

	return name, version, nil
}

func newWithCallback(nme string, configPath string, reload func(fsnotify.Event)) (Application, error) {
	name, version, err := ensureDefaults(nme)
	if err != nil {
		return nil, err
	}
	appInfo := cjm.AppInfo{
		Name:     name,
		BasePath: "/",
		Version:  version,
		Pid:      os.Getpid(),
	}

	cfg, err := createViper(name, configPath)
	if err != nil {
		return nil, err
	}

	allLoggers := logging.NewRegistry(cfg, logrus.Fields{"app": appInfo.Name})

	log.SetOutput(allLoggers.Writer())

	tracer := allLoggers.Root().WithField("module", "trace")
	trace := tracing.New("", tracer, nil)

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGQUIT)
		buf := make([]byte, 1<<20)

		for {
			<-sigs
			ln := goruntime.Stack(buf, true)
			allLoggers.Root().Println(string(buf[:ln]))
		}
	}()

	app := &defaultApplication{
		appInfo:    appInfo,
		allLoggers: allLoggers,
		rootTracer: trace,
		config:     cfg,
		registry:   make(map[Key]interface{}, 100),
		regLock:    new(sync.Mutex),
	}
	app.watchConfigurations(func(in fsnotify.Event) {
		if reload != nil {
			reload(in)
		}
		allLoggers.Reload()
		for _, mod := range app.modules {
			if err := mod.Reload(app); err != nil {
				allLoggers.Root().Errorf("reload config: %v", err)
			}
		}
		allLoggers.Root().Infoln("config file changed:", in.Name)
	})
	return app, nil
}

// New application with the specified name, at the specified basepath
func New(nme string) (Application, error) {
	return newWithCallback(nme, "", nil)
}

// NewWithConfig application with the specified name, with a specific config file path
func NewWithConfig(nme string, configPath string) (Application, error) {
	return newWithCallback(nme, configPath, nil)
}

type defaultApplication struct {
	appInfo    cjm.AppInfo
	allLoggers *logging.Registry
	rootTracer tracing.Tracer
	config     *viper.Viper
	modules    []Module

	registry map[Key]interface{}
	regLock  *sync.Mutex
}

func (d *defaultApplication) watchConfigurations(reload func(fsnotify.Event)) {
	viperLock.Lock()
	defer viperLock.Unlock()
	d.config.WatchConfig()
	d.config.OnConfigChange(reload)

	// we made it this far, it's clear the url means we're also connecting remotely
	if os.Getenv("CONFIG_REMOTE_URL") != "" {
		go func() {
			for {
				err := d.config.WatchRemoteConfig()
				if err != nil {
					d.Logger().Errorf("watching remote config: %v", err)
					continue
				}
				reload(fsnotify.Event{Name: os.Getenv("CONFIG_REMOTE_URL"), Op: fsnotify.Write})
			}
		}()
	}
}

func (d *defaultApplication) Add(modules ...Module) error {
	if len(modules) > 0 {
		d.modules = append(d.modules, modules...)
	}
	return nil
}

// Get the module at the specified key, return nil when the component doesn't exist
func (d *defaultApplication) Get(key Key) interface{} {
	mod, _ := d.GetOK(key)
	return mod
}

// Get the module at the specified key, return false when the component doesn't exist
func (d *defaultApplication) GetOK(key Key) (interface{}, bool) {
	d.regLock.Lock()
	defer d.regLock.Unlock()

	mod, ok := d.registry[key]
	if !ok {
		return nil, ok
	}
	return mod, ok
}

func (d *defaultApplication) Set(key Key, module interface{}) error {
	d.regLock.Lock()
	d.registry[key] = module
	d.regLock.Unlock()
	return nil
}

func (d *defaultApplication) Logger() logrus.FieldLogger {
	return d.allLoggers.Root()
}

func (d *defaultApplication) NewLogger(name string, ctx logrus.Fields) logrus.FieldLogger {
	return d.allLoggers.Root().New(name, ctx)
}

func (d *defaultApplication) Tracer() tracing.Tracer {
	return d.rootTracer
}

func (d *defaultApplication) Config() *viper.Viper {
	return d.config
}

func (d *defaultApplication) Info() cjm.AppInfo {
	return d.appInfo
}

func (d *defaultApplication) Init() error {
	for _, mod := range d.modules {
		if err := mod.Init(d); err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultApplication) Start() error {
	for _, mod := range d.modules {
		if err := mod.Start(d); err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultApplication) Stop() error {
	for _, mod := range d.modules {
		if err := mod.Stop(d); err != nil {
			return err
		}
	}
	return nil
}
