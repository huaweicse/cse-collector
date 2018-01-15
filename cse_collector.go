package metricsink

import (
	"crypto/tls"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix/metric_collector"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"time"
)

// CseCollector is a struct to keeps metric information of Http requests
type CseCollector struct {
	attempts          string
	errors            string
	successes         string
	failures          string
	rejects           string
	shortCircuits     string
	timeouts          string
	fallbackSuccesses string
	fallbackFailures  string
	totalDuration     string
	runDuration       string
}

// CseCollectorConfig is a struct to keep monitoring information
type CseCollectorConfig struct {
	// CseMonitorAddr is the http address of the csemonitor server
	CseMonitorAddr string
	// Headers for csemonitor server
	Header http.Header
	// TickInterval spcifies the period that this collector will send metrics to the server.
	TimeInterval time.Duration
	// Config structure to configure a TLS client for sending Metric data
	TLSConfig *tls.Config
}

// InitializeCseCollector starts the CSE collector in a new Thread
func InitializeCseCollector(config *CseCollectorConfig, r metrics.Registry) {
	go CseMonitor(r, config.CseMonitorAddr, config.Header, config.TimeInterval, config.TLSConfig)
}

// NewCseCollector creates a new Collector Object
func NewCseCollector(name string) metricCollector.MetricCollector {
	return &CseCollector{
		attempts:          name + ".attempts",
		errors:            name + ".errors",
		successes:         name + ".successes",
		failures:          name + ".failures",
		rejects:           name + ".rejects",
		shortCircuits:     name + ".shortCircuits",
		timeouts:          name + ".timeouts",
		fallbackSuccesses: name + ".fallbackSuccesses",
		fallbackFailures:  name + ".fallbackFailures",
		totalDuration:     name + ".totalDuration",
		runDuration:       name + ".runDuration",
	}
}

func (c *CseCollector) incrementCounterMetric(prefix string) {
	count, ok := metrics.GetOrRegister(prefix, metrics.NewCounter).(metrics.Counter)
	if !ok {
		return
	}
	count.Inc(1)
}

func (c *CseCollector) updateTimerMetric(prefix string, dur time.Duration) {
	count, ok := metrics.GetOrRegister(prefix, metrics.NewTimer).(metrics.Timer)
	if !ok {
		return
	}
	count.Update(dur)
}

// IncrementAttempts function increments the number of calls to this circuit.
// This registers as a counter
func (c *CseCollector) IncrementAttempts() {
	c.incrementCounterMetric(c.attempts)
}

// IncrementErrors function increments the number of unsuccessful attempts.
// Attempts minus Errors will equal successes.
// Errors are result from an attempt that is not a success.
// This registers as a counter
func (c *CseCollector) IncrementErrors() {
	c.incrementCounterMetric(c.errors)

}

// IncrementSuccesses function increments the number of requests that succeed.
// This registers as a counter
func (c *CseCollector) IncrementSuccesses() {
	c.incrementCounterMetric(c.successes)

}

// IncrementFailures function increments the number of requests that fail.
// This registers as a counter
func (c *CseCollector) IncrementFailures() {
	c.incrementCounterMetric(c.failures)
}

// IncrementRejects function increments the number of requests that are rejected.
// This registers as a counter
func (c *CseCollector) IncrementRejects() {
	c.incrementCounterMetric(c.rejects)
}

// IncrementShortCircuits function increments the number of requests that short circuited due to the circuit being open.
// This registers as a counter
func (c *CseCollector) IncrementShortCircuits() {
	c.incrementCounterMetric(c.shortCircuits)
}

// IncrementTimeouts function increments the number of timeouts that occurred in the circuit breaker.
// This registers as a counter
func (c *CseCollector) IncrementTimeouts() {
	c.incrementCounterMetric(c.timeouts)
}

// IncrementFallbackSuccesses function increments the number of successes that occurred during the execution of the fallback function.
// This registers as a counter
func (c *CseCollector) IncrementFallbackSuccesses() {
	c.incrementCounterMetric(c.fallbackSuccesses)
}

// IncrementFallbackFailures function increments the number of failures that occurred during the execution of the fallback function.
// This registers as a counter
func (c *CseCollector) IncrementFallbackFailures() {
	c.incrementCounterMetric(c.fallbackFailures)
}

// UpdateTotalDuration function updates the internal counter of how long we've run for.
// This registers as a timer
func (c *CseCollector) UpdateTotalDuration(timeSinceStart time.Duration) {
	c.updateTimerMetric(c.totalDuration, timeSinceStart)
}

// UpdateRunDuration function updates the internal counter of how long the last run took.
// This registers as a timer
func (c *CseCollector) UpdateRunDuration(runDuration time.Duration) {
	c.updateTimerMetric(c.runDuration, runDuration)
}

// Reset function is a noop operation in this collector.
func (c *CseCollector) Reset() {}
