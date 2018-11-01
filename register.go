package metricsink

import (
	"time"

	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/core/config"
	"github.com/go-chassis/go-chassis/core/lager"
	chassisMetrics "github.com/go-chassis/go-chassis/metrics"
	"github.com/rcrowley/go-metrics"
)

func init() {
	chassisMetrics.InstallReporter("CSE Monitoring", reportMetricsToCSEDashboard)
}

//reportMetricsToCSEDashboard use go-metrics to send metrics to cse dashboard
func reportMetricsToCSEDashboard(r metrics.Registry) error {

	monitorServerURL, err := getMonitorEndpoint()
	if err != nil {
		lager.Logger.Warnf("Get Monitoring URL failed, CSE monitoring function disabled, err: %v", err)
		return nil
	}

	tlsConfig, tlsError := getTLSForClient(monitorServerURL)
	if tlsError != nil {
		lager.Logger.Errorf("Get %s.%s TLS config failed,error : %s", monitorServerURL, common.Consumer, tlsError)
		return tlsError
	}

	InitializeCseCollector(&CseCollectorConfig{
		CseMonitorAddr: monitorServerURL,
		Header:         getAuthHeaders(),
		TimeInterval:   time.Second * 2,
		TLSConfig:      tlsConfig,
	}, r, config.GlobalDefinition.AppID, config.SelfVersion, config.SelfServiceName,
		config.MicroserviceDefinition.ServiceDescription.Environment)
	lager.Logger.Infof("Started sending metric Data to Monitor Server : %s", monitorServerURL)
	return nil
}
