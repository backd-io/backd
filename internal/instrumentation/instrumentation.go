package instrumentation

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

// Instrumentation holds the instrumentation definitions and log on a package to
//  ensure every service behaves exactly the same.
type Instrumentation struct {
	ipPort     string
	httpServer *http.Server
	metrics    map[string]prometheus.Collector
	logger     *zap.Logger
	hostname   string
	sync.Mutex
}

// New returns an Instrumentation object
func New(ipPort string, debug bool) (*Instrumentation, error) {

	var (
		loggerLevel     zap.AtomicLevel
		loggerConf      zap.Config
		logger          *zap.Logger
		instrumentation Instrumentation
		err             error
	)

	if debug {
		loggerLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	} else {
		loggerLevel = zap.NewAtomicLevelAt(zap.InfoLevel) // Default logger level (as production)
	}

	loggerConf = zap.Config{
		Level:       loggerLevel,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err = loggerConf.Build()
	if err != nil {
		return &instrumentation, err
	}

	instrumentation.hostname, err = os.Hostname()
	if err != nil {
		return &instrumentation, err
	}

	instrumentation.ipPort = ipPort
	instrumentation.metrics = make(map[string]prometheus.Collector)
	instrumentation.logger = logger

	return &instrumentation, err

}

// Start the instrumentation server
func (i *Instrumentation) Start() error {
	i.httpServer = new(http.Server)
	handler := http.NewServeMux()
	handler.Handle("/metrics", promhttp.Handler())
	i.httpServer.Addr = i.ipPort
	i.httpServer.Handler = handler
	return i.httpServer.ListenAndServe()
}

// Shutdown stops the instrumentation server
func (i *Instrumentation) Shutdown() error {
	return i.httpServer.Shutdown(context.Background())
}

// RegisterGauge registers gauges
func (i *Instrumentation) RegisterGauge(name, help string, labels []string) {

	i.Lock()
	defer i.Unlock()

	switch len(labels) {
	case 0:
		i.metrics[name] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: name,
				Help: help,
			})
	default:
		i.metrics[name] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: name,
				Help: help,
			},
			labels)
	}

	prometheus.MustRegister(i.metrics[name])

}

// RegisterHistogram registers a new histogram to be consumed by Prometheus
func (i *Instrumentation) RegisterHistogram(name, help string, labels []string, buckets []float64) {

	i.Lock()
	defer i.Unlock()

	switch len(labels) {
	case 0:
		i.metrics[name] = prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Name:    name,
				Help:    help,
				Buckets: buckets,
			})
	default:
		i.metrics[name] = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    name,
				Help:    help,
				Buckets: buckets,
			},
			labels)
	}

	prometheus.MustRegister(i.metrics[name])

}

// Hostname returns the name of the host being instrumented
func (i *Instrumentation) Hostname() string {
	return i.hostname
}

// Metric returns a metric from the map safely
func (i *Instrumentation) Metric(name string) prometheus.Collector {
	i.Lock()
	defer i.Unlock()
	return i.metrics[name]
}

// Error uses the zap.Logger logging functionality
func (i *Instrumentation) Error(msg string, fields ...zap.Field) {
	i.logger.Error(msg, fields...)
}

// Info uses the zap.Logger logging functionality
func (i *Instrumentation) Info(msg string, fields ...zap.Field) {
	i.logger.Info(msg, fields...)
}

// Debug uses the zap.Logger logging functionality
func (i *Instrumentation) Debug(msg string, fields ...zap.Field) {
	i.logger.Debug(msg, fields...)
}
