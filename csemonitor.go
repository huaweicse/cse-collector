package metricsink

import (
	"crypto/tls"
	"github.com/ServiceComb/go-chassis/core/archaius"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/rcrowley/go-metrics"
	"net/http"
	"os"
	"runtime"
	"time"
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
}

// NewReporter creates a New reporter for monitoring
func NewReporter(r metrics.Registry, addr string, header http.Header, interval time.Duration, tls *tls.Config) *Reporter {
	return &Reporter{
		Registry:       r,
		CseMonitorAddr: addr,
		Header:         header,
		Interval:       interval,
		Percentiles:    []float64{0.5, 0.75, 0.95, 0.99, 0.999},
		TLSConfig:      tls,
	}
}

// CseMonitor creates a new monitoring object for CSE type collections
func CseMonitor(r metrics.Registry, addr string, header http.Header, interval time.Duration, tls *tls.Config) {
	NewReporter(r, addr, header, interval, tls).Run()
}

// Run creates a go_routine which runs continously and capture the monitoring data
func (reporter *Reporter) Run() {
	ticker := time.Tick(reporter.Interval)
	metricsAPI := NewCseMonitorClient(reporter.Header, reporter.CseMonitorAddr, reporter.TLSConfig)
	IsMonitoringConnected = true
	isConnctedForFirstTime := false

	for range ticker {

		//If monitoring is enabled then only try to connect to Monitoring Server
		if archaius.GetBool("cse.monitor.client.enable", true) {
			monitorData := reporter.getData()
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
func (reporter *Reporter) getData() MonitorData {
	var monitorData = NewMonitorData()
	monitorData.AppID = config.GlobalDefinition.AppID
	monitorData.Version = config.SelfVersion
	monitorData.Name = config.SelfServiceName
	monitorData.Instance, _ = os.Hostname()
	monitorData.Memory = getProcessInfo()
	monitorData.Thread = threadCreateProfile.Count()
	monitorData.CPU = float64(runtime.NumCPU())
	reporter.Registry.Each(func(name string, i interface{}) {
		monitorData.appendInterfaceInfo(name, i)
	})
	return *monitorData
}
