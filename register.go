package metricsink

import (
	"time"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/pkg/runtime"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
)

func init() {
	hystrix.InstallReporter("CSE Monitoring", reportMetricsToCSEDashboard)
}

//reportMetricsToCSEDashboard use go-metrics to send metrics to cse dashboard
func reportMetricsToCSEDashboard(cb *hystrix.CircuitBreaker) error {
	monitorServerURL, err := getMonitorEndpoint()
	if err != nil {
		openlogging.GetLogger().Warnf("Get Monitoring URL failed, CSE monitoring function disabled, err: %v", err)
		return nil
	}

	tlsConfig, tlsError := getTLSForClient(monitorServerURL)
	if tlsError != nil {
		openlogging.GetLogger().Errorf("Get %s.%s TLS config failed,error : %s", monitorServerURL, common.Consumer, tlsError)
		return tlsError
	}

	InitializeCseCollector(&CseCollectorConfig{
		CseMonitorAddr: monitorServerURL,
		Header:         getAuthHeaders(),
		TimeInterval:   time.Second * 2,
		TLSConfig:      tlsConfig,
	}, runtime.App, runtime.Version, runtime.Version,
		config.MicroserviceDefinition.ServiceDescription.Environment)
	openlogging.GetLogger().Infof("Started sending metric Data to Monitor Server : %s", monitorServerURL)
	return nil
}
