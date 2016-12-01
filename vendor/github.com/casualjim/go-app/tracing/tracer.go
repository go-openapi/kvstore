package tracing

import (
	"runtime"
	"time"

	logrus "github.com/Sirupsen/logrus"
	metrics "github.com/rcrowley/go-metrics"
)

const (
	noMethodName = "<anonymous>"
	trace        = "trace"
)

// Tracer interface that represents a tracer in golang
type Tracer interface {
	// Trace the enter/leave of the method.
	// Record time spent in the method.
	// Returns a closure to close the method, best used in conjunction with defer, eg.: defer tr.Trace()()
	Trace(name ...string) func()
}

// NewTracer creates a new tracer object with the specified configuration
// When the config is nil the tracer will use default values for the config,
// this is equivalent to
//
//      name: trace
func New(name string, logger logrus.FieldLogger, registry metrics.Registry) Tracer {
	nm := name
	if nm == "" {
		nm = trace
	}

	var bl = logger
	if bl == nil {
		bl = logrus.WithField("name", nm)
	}

	reg := registry
	if reg == nil {
		reg = metrics.DefaultRegistry
	}
	return &defaultTracing{logger: bl, registry: reg}
}

type defaultTracing struct {
	logger   logrus.FieldLogger
	registry metrics.Registry
}

func (d *defaultTracing) Trace(methods ...string) func() {
	var method string
	if len(methods) == 0 || methods[0] == "" {
		method = noMethodName
		pc, _, _, ok := runtime.Caller(1)
		if ok {
			fun := runtime.FuncForPC(pc)
			if fun != nil {
				method = fun.Name()
			}
		}
	} else {
		method = methods[0]
	}

	timer := metrics.GetOrRegisterTimer(method, d.registry)
	d.logger.Debugf("Enter %s ", method)

	start := time.Now()
	return func() {
		timer.UpdateSince(start)
		d.logger.Debugf("Leave %s took %v", method, time.Now().Sub(start))
	}
}
