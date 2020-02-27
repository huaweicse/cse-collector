package metricsink

// Forked from github.com/afex/hystrix-go
// Some parts of this file have been modified to make it functional in this package

import (
	"crypto/tls"
	"fmt"
	"github.com/go-chassis/go-chassis/core/common"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
	"net/http"
	"sync"
	"time"
)

var initOnce = sync.Once{}

// CseCollectorConfig is a struct to keep monitoring information
type CseCollectorConfig struct {
	// CseMonitorAddr is the http address of the csemonitor server
	CseMonitorAddr string
	// Headers for csemonitor server
	Header http.Header
	// TickInterval specifies the period that this collector will send metrics to the server.
	TimeInterval time.Duration
	// Config structure to configure a TLS client for sending Metric data
	TLSConfig *tls.Config

	Env string
}

func init() {
	hystrix.InstallReporter("CSE Monitoring", reportMetricsToCSEDashboard)
}

var reporter *Reporter

//GetReporter get reporter singleton
func GetReporter() (*Reporter, error) {
	var errResult error
	initOnce.Do(func() {
		monitorServerURL, err := getMonitorEndpoint()
		if err != nil {
			openlogging.GetLogger().Warnf("Get Monitoring URL failed, CSE monitoring function disabled, err: %v", err)
			errResult = err
			return

		}
		openlogging.GetLogger().Infof("init monitoring client : %s", monitorServerURL)
		tlsConfig, err := getTLSForClient(monitorServerURL)
		if err != nil {
			openlogging.GetLogger().Errorf("Get %s.%s TLS config failed,error : %s", monitorServerURL, common.Consumer, err)
			errResult = err
		}
		reporter, err = NewReporter(&CseCollectorConfig{
			CseMonitorAddr: monitorServerURL,
			Header:         getAuthHeaders(),
			TLSConfig:      tlsConfig,
		})
		if err != nil {
			openlogging.Error("new reporter failed", openlogging.WithTags(openlogging.Tags{
				"err": err.Error(),
			}))
			errResult = err
		}
	})

	if reporter == nil {
		errResult = fmt.Errorf("Reporter is nil")
	}
	return reporter, errResult
}

//reportMetricsToCSEDashboard use send metrics to cse dashboard
func reportMetricsToCSEDashboard(cb *hystrix.CircuitBreaker) error {
	r, err := GetReporter()
	if err != nil {
		return err
	}
	r.Send(cb)
	return nil
}
