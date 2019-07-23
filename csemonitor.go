package metricsink

import (
	"github.com/go-chassis/go-archaius"
	chassisRuntime "github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
	"github.com/huaweicse/cse-collector/pkg/monitoring"
	"os"
	"runtime"
)

// IsMonitoringConnected is a boolean to keep an check if there exsist any succeful connection to monitoring Server
var IsMonitoringConnected bool

// Reporter is a struct to store the registry address and different monitoring information
type Reporter struct {
	environment string
	c           *monitoring.CseMonitorClient
}

// NewReporter creates a new monitoring object for CSE type collections
func NewReporter(config *CseCollectorConfig) (*Reporter, error) {
	c, err := monitoring.NewCseMonitorClient(config.Header, config.CseMonitorAddr, config.TLSConfig)
	if err != nil {
		openlogging.GetLogger().Errorf("Get cse monitor client failed:%s", err)
		return nil, err
	}
	IsMonitoringConnected = true
	return &Reporter{
		environment: config.Env,
		c:           c,
	}, nil
}

//Send send metrics to monitoring service
func (reporter *Reporter) Send(cb *hystrix.CircuitBreaker) {
	if archaius.GetBool("cse.monitor.client.enable", true) {
		monitorData := reporter.getData(cb)
		openlogging.Debug("send metrics", openlogging.WithTags(openlogging.Tags{
			"data": monitorData,
		}))
		err := reporter.c.PostMetrics(monitorData)
		if err != nil {
			openlogging.GetLogger().Warnf("unable to report to monitoring server, err: %v", err)
		}
	}
}

func (reporter *Reporter) getData(cb *hystrix.CircuitBreaker) monitoring.MonitorData {
	var monitorData = monitoring.NewMonitorData()
	monitorData.AppID = chassisRuntime.App
	monitorData.Version = chassisRuntime.Version
	monitorData.Name = chassisRuntime.ServiceName
	monitorData.ServiceID = chassisRuntime.ServiceID
	monitorData.InstanceID = chassisRuntime.InstanceID

	monitorData.Environment = reporter.environment
	monitorData.Instance, _ = os.Hostname()
	monitorData.Memory = getProcessInfo()
	monitorData.Thread = threadCreateProfile.Count()
	monitorData.CPU = float64(runtime.NumCPU())
	monitorData.AppendInterfaceInfo(cb)
	return *monitorData
}
