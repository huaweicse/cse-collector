package metricsink

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/go-chassis/go-archaius"
	chassisRuntime "github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
	"runtime"
)

// IsMonitoringConnected is a boolean to keep an check if there exsist any succeful connection to monitoring Server
var IsMonitoringConnected bool

// Reporter is a struct to store the registry address and different monitoring information
type Reporter struct {
	CseMonitorAddr string
	Header         http.Header
	Interval       time.Duration
	Percentiles    []float64
	TLSConfig      *tls.Config
	app            string
	version        string
	service        string
	environment    string
	serviceID      string
	metricsAPI     *CseMonitorClient
}

// NewReporter creates a new monitoring object for CSE type collections
func NewReporter(addr string, header http.Header, interval time.Duration, tls *tls.Config, app, version, service, env string) *Reporter {
	reporter := &Reporter{
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
	metricsAPI, err := NewCseMonitorClient(reporter.Header, reporter.CseMonitorAddr, reporter.TLSConfig, "v2")
	if err != nil {
		openlogging.GetLogger().Errorf("Get cse monitor client failed:%s", err)
	}
	reporter.metricsAPI = metricsAPI
	IsMonitoringConnected = true
	return reporter
}

// Run creates a go_routine which runs continuously and capture the monitoring data
func (reporter *Reporter) Run(cb *hystrix.CircuitBreaker) {
	ticker := time.Tick(reporter.Interval)

	for range ticker {
		if archaius.GetBool("cse.monitor.client.enable", true) {
			reporter.serviceID = chassisRuntime.ServiceID
			monitorData := reporter.getData(cb, reporter.app, reporter.version,
				reporter.service, reporter.environment, reporter.serviceID, chassisRuntime.InstanceID)
			err := reporter.metricsAPI.PostMetrics(monitorData)
			if err != nil {
				openlogging.GetLogger().Warnf("Unable to report to monitoring server, err: %v", err)
			}
		}
	}
}

func (reporter *Reporter) getData(cb *hystrix.CircuitBreaker,
	app, version, service, env, serviceID, instanceID string) MonitorData {
	var monitorData = NewMonitorData()
	monitorData.AppID = app
	monitorData.Version = version
	monitorData.Name = service
	monitorData.ServiceID = serviceID
	monitorData.InstanceID = instanceID
	monitorData.Environment = env
	monitorData.Instance, _ = os.Hostname()
	monitorData.Memory = getProcessInfo()
	monitorData.Thread = threadCreateProfile.Count()
	monitorData.CPU = float64(runtime.NumCPU())
	monitorData.appendInterfaceInfo(cb.Name, cb.Metrics.DefaultCollector())
	return *monitorData
}
