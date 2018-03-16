package metricsink

import (
	"crypto/tls"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/rcrowley/go-metrics"
)

// IsMonitoringConnected is a boolean to keep an check if there exsist any succeful connection to monitoring Server
var IsMonitoringConnected bool

// Reporter is a struct to store the registry address and different monitoring information
type Reporter struct {
	Registry       metrics.Registry
	CseMonitorAddr string
	Header         http.Header
	Interval       time.Duration
	Percentiles    []float64
	TLSConfig      *tls.Config
	app            string
	version        string
	service        string
	environment    string
}

// NewReporter creates a new monitoring object for CSE type collections
func NewReporter(r metrics.Registry, addr string, header http.Header, interval time.Duration, tls *tls.Config, app, version, service, env string) *Reporter {
	reporter := &Reporter{
		Registry:       r,
		CseMonitorAddr: addr,
		Header:         header,
		Interval:       interval,
		Percentiles:    []float64{0.5, 0.75, 0.95, 0.99, 0.999},
		TLSConfig:      tls,
		app:            app,
		version:        version,
		service:        service,
		environment:    env,
	}
	return reporter
}

// Run creates a go_routine which runs continuously and capture the monitoring data
func (reporter *Reporter) Run() {
	ticker := time.Tick(reporter.Interval)
	metricsAPI := NewCseMonitorClient(reporter.Header, reporter.CseMonitorAddr, reporter.TLSConfig, "v2")
	IsMonitoringConnected = true
	isConnctedForFirstTime := false

	for range ticker {

		//If monitoring is enabled then only try to connect to Monitoring Server
		if archaius.GetBool("cse.monitor.client.enable", true) {
			monitorData := reporter.getData(reporter.app, reporter.version, reporter.service, reporter.environment)
			err := metricsAPI.PostMetrics(monitorData)
			if err != nil {
				//If the connection fails for the first time then print Warn Logs
				if IsMonitoringConnected {
					lager.Logger.Warnf(err, "Unable to connect to monitoring server")
				}
				IsMonitoringConnected = false
			} else {
				//If Connection is established for first time then Print the Info logs
				if !isConnctedForFirstTime {
					lager.Logger.Infof("Connection to monitoring server established successfully")
					isConnctedForFirstTime = true
				}
				//If connection is recovered first time then print Info Logs
				if !IsMonitoringConnected {
					lager.Logger.Infof("Connection recovered successfully to monitoring server")
				}
				IsMonitoringConnected = true
			}
		}
	}
}
func (reporter *Reporter) getData(app, version, service, env string) MonitorData {
	var monitorData = NewMonitorData()
	monitorData.AppID = app
	monitorData.Version = version
	monitorData.Name = service
	monitorData.Environment = env
	monitorData.Instance, _ = os.Hostname()
	monitorData.Memory = getProcessInfo()
	monitorData.Thread = threadCreateProfile.Count()
	monitorData.CPU = float64(runtime.NumCPU())
	reporter.Registry.Each(func(name string, i interface{}) {
		monitorData.appendInterfaceInfo(name, i)
	})
	return *monitorData
}
