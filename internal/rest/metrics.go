package rest

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func (rr *REST) registerMetrics() {
	// rr.inst.RegisterHistogram("backd_rest_ops", "REST durations", []string{"hostname", "method", "uri", "code"}, []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5, 1, 5, 10})
	rr.inst.RegisterSummary("backd_rest_ops", "REST durations", []string{"hostname", "method", "uri", "code"})
	rr.inst.RegisterCounter("backd_rest_counter", "REST counters", []string{"hostname", "method", "uri", "code"})
}

// addOperationDuration
// name  = name of the operation on the map
// op    = label for operation (select, update, etc)
// app   = appID
// table = table name
// dur   = duration for the operation
func (rr *REST) addOperationDuration(name, method, uri, code string, dur time.Duration) {
	// rr.inst.Metric(name).(*prometheus.HistogramVec).WithLabelValues(rr.inst.Hostname(), method, uri, code).Observe(dur.Seconds())
	rr.inst.Metric(name).(*prometheus.SummaryVec).WithLabelValues(rr.inst.Hostname(), method, uri, code).Observe(dur.Seconds())
}

func (rr *REST) addOperationCounter(name, method, uri, code string) {
	rr.inst.Metric(name).(*prometheus.CounterVec).WithLabelValues(rr.inst.Hostname(), method, uri, code).Add(1)
}
