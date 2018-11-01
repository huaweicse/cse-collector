package metricsink

// Forked from github.com/afex/hystrix-go
// Some parts of this file have been modified to make it functional in this package

import (
	"crypto/tls"
	"net/http"
	"time"

	"github.com/rcrowley/go-metrics"
)

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
func InitializeCseCollector(config *CseCollectorConfig, r metrics.Registry, app, version, service, env string) {
	go NewReporter(r, config.CseMonitorAddr, config.Header, config.TimeInterval, config.TLSConfig, app, version, service, env).Run()
}
