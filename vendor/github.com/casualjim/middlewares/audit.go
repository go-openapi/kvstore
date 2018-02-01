package middlewares

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/rcrowley/go-metrics"
	// ensure pprof is loaded and its http hooks installed
	_ "net/http/pprof"

	"github.com/justinas/alice"
)

// Logger interface to allow different logging libaries to be used with this middleware
type Logger interface {
	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
}

func init() {
	metrics.RegisterDebugGCStats(metrics.DefaultRegistry)
	go metrics.CaptureDebugGCStats(metrics.DefaultRegistry, 5e9)

	metrics.RegisterRuntimeMemStats(metrics.DefaultRegistry)
	go metrics.CaptureRuntimeMemStats(metrics.DefaultRegistry, 5e9)
}

// Audit is a middleware handler that logs the request as it goes in and the response as it goes out.
type Audit struct {
	// logger is the Logger instance used to log messages with the Audit middleware
	Logger   Logger
	basePath string
	next     http.Handler
	info     AppInfo
}

// NewAuditMW returns a new Audit middleware
func NewAuditMW(info AppInfo, logger Logger) alice.Constructor {
	return func(hand http.Handler) http.Handler {
		return NewAudit(info, logger, hand)
	}
}

// NewAudit returns a new Audit instance
func NewAudit(info AppInfo, logger Logger, next http.Handler) *Audit {
	basePath := info.BasePath
	if basePath == "" {
		basePath = "/"
	}
	return &Audit{
		Logger:   logger,
		basePath: basePath,
		next:     next,
		info:     info,
	}
}

// CounterMetric represents a value that increases or decreases
type CounterMetric struct {
	Count int64 `json:"count"`
}

// GaugeMetric represents measurements.
type GaugeMetric struct {
	Value int64 `json:"value"`
}

// GaugeFloat64Metric represents a float64 measurement
type GaugeFloat64Metric struct {
	Value float64 `json:"value"`
}

// HealtCheckData shows error status if any for a health check
type HealtCheckData struct {
	Error string `json:"error"`
}

// HistogramMetric shows a histogram metric
type HistogramMetric struct {
	Count  int64   `json:"count"`
	Min    int64   `json:"min"`
	Max    int64   `json:"max"`
	Mean   float64 `json:"mean"`
	StdDev float64 `json:"stdDev"`
	Median float64 `json:"median"`
	P75    float64 `json:"p75"`
	P95    float64 `json:"p95"`
	P99    float64 `json:"p99"`
	P999   float64 `json:"p999"`
	Unit   string  `json:"unit,omitempty"`
}

// MeterMetric represents a metered value with 1-minute, 5-minute and 15 minute averages as well as a mean
type MeterMetric struct {
	Count int64   `json:"count"`
	M1    float64 `json:"m1"`
	M5    float64 `json:"m5"`
	M15   float64 `json:"m15"`
	Mean  float64 `json:"mean"`
}

// TimerMetric is a combination of a meter metric and a histogram metric, typically used for tracing time spent and call counts for example
type TimerMetric struct {
	Rate     MeterMetric     `json:"rate"`
	Duration HistogramMetric `json:"duration"`
}

// AppInfo the information describing the component for this API
type AppInfo struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Commit   string `json:"commit"`
	BasePath string `json:"basePath"`
	Pid      int    `json:"pid"`
}

func (l *Audit) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, filepath.Join(l.basePath, "audit", "metrics")) {
		data := make(map[string]interface{})
		metrics.DefaultRegistry.Each(func(name string, i interface{}) {
			var values interface{}
			switch metric := i.(type) {
			case metrics.Counter:
				values = &CounterMetric{metric.Count()}
			case metrics.Gauge:
				values = &GaugeMetric{metric.Value()}
			case metrics.GaugeFloat64:
				values = &GaugeFloat64Metric{metric.Value()}
			case metrics.Healthcheck:
				metric.Check()
				if err := metric.Error(); nil != err {
					values = &HealtCheckData{metric.Error().Error()}
				}
			case metrics.Histogram:
				h := metric.Snapshot()
				ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				values = &HistogramMetric{
					Count:  h.Count(),
					Min:    h.Min(),
					Max:    h.Max(),
					Mean:   h.Mean(),
					StdDev: h.StdDev(),
					Median: ps[0],
					P75:    ps[1],
					P95:    ps[2],
					P99:    ps[3],
					P999:   ps[4],
				}
			case metrics.Meter:
				m := metric.Snapshot()
				values = &MeterMetric{
					Count: m.Count(),
					M1:    m.Rate1(),
					M5:    m.Rate5(),
					M15:   m.Rate15(),
					Mean:  m.RateMean(),
				}
			case metrics.Timer:
				t := metric.Snapshot()
				ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
				duration := HistogramMetric{
					Count:  t.Count(),
					Min:    t.Min(),
					Max:    t.Max(),
					Mean:   t.Mean(),
					StdDev: t.StdDev(),
					Median: ps[0],
					P75:    ps[1],
					P95:    ps[2],
					P99:    ps[3],
					P999:   ps[4],
					Unit:   "nanoseconds",
				}
				rate := MeterMetric{
					Count: t.Count(),
					M1:    t.Rate1(),
					M5:    t.Rate5(),
					M15:   t.Rate15(),
					Mean:  t.RateMean(),
				}
				values = &TimerMetric{
					Rate:     rate,
					Duration: duration,
				}

			}
			data[name] = values
		})
		enc := json.NewEncoder(rw)
		rw.Header().Set("Content-Type", "application/json;charset=utf-8")
		enc.Encode(data)
	} else if strings.HasPrefix(r.URL.Path, filepath.Join(l.basePath, "audit", "info")) {
		enc := json.NewEncoder(rw)
		rw.Header().Set("Content-Type", "application/json;charset=utf-8")
		enc.Encode(l.info)
	} else {
		l.Logger.Debugf("Begin %s %s", r.Method, r.URL.Path)
		timer := metrics.GetOrRegisterTimer(r.URL.Path, metrics.DefaultRegistry)
		start := time.Now()
		timer.Time(func() {
			l.next.ServeHTTP(rw, r)
		})
		l.Logger.Infof("%s %s (took: %s)", r.Method, r.URL.Path, time.Now().Sub(start).String())
	}

}
