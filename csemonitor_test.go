package metricsink

import (
	"crypto/tls"
	"fmt"
	"github.com/ServiceComb/go-chassis/core/config"
	"github.com/ServiceComb/go-chassis/core/config/model"
	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/rcrowley/go-metrics"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var globalConf = `
---
APPLICATION_ID: CSE

cse:
  monitor:
      client:
        serverUri: https://10.21.209.37:30109   #monitor server url
        enable: true                    # if enable is false then it will not send the metric data to monitor server
        userName : weixing
        domainName : default
  protocols:
    highway:
      listenAddress: 127.0.0.1:8080
      advertiseAddress: 127.0.0.1:8080
      transport: tcp #optional 指定加载那个传输层
    rest:
      listenAddress: 127.0.0.1:8081
      advertiseAddress: 127.0.0.1:8081
      transport: tcp #optional 指定加载那个传输层
  handler:
    chain:
      provider:
        default: bizkeeper-provider

`

func initEnv() {
	lager.Initialize("", "INFO", "", "size", true, 1, 10, 7)
	config.GlobalDefinition = new(model.GlobalCfg)
	yaml.Unmarshal([]byte(globalConf), config.GlobalDefinition)
}
func TestNewReporter(t *testing.T) {
	initEnv()
	assert := assert.New(t)
	reporter := NewReporter(metrics.DefaultRegistry, "127.0.0.1:8080", http.Header{"Content-Type": []string{"application/json"}}, time.Second, &tls.Config{})
	assert.Equal(reporter.Interval, time.Second)
	assert.Equal(reporter.CseMonitorAddr, "127.0.0.1:8080")
	assert.Equal(reporter.Header, http.Header{"Content-Type": []string{"application/json"}})
}
func TestCseMonitor(t *testing.T) {
	initEnv()
	assert := assert.New(t)
	reporter := NewReporter(metrics.DefaultRegistry, "127.0.0.1:8080", http.Header{"Content-Type": []string{"application/json"}}, time.Second, &tls.Config{})
	config.SelfServiceName = "testService"
	monitorData := reporter.getData()
	assert.Equal(monitorData.Name, "testService")

}
func TestCseMonitor2(t *testing.T) {
	initEnv()
	assert := assert.New(t)
	reporter := NewReporter(metrics.DefaultRegistry, "127.0.0.1:8080", http.Header{"Content-Type": []string{"application/json"}}, time.Second, &tls.Config{})
	metricCollector := NewCseCollector("source.Provider.Microservice.SchemaID.OperationId")
	config.SelfServiceName = "testService"

	metricCollector.IncrementAttempts()
	metricCollector.IncrementErrors()
	metricCollector.IncrementFailures()
	metricCollector.IncrementFallbackFailures()
	metricCollector.IncrementFallbackSuccesses()
	metricCollector.IncrementShortCircuits()
	metricCollector.UpdateTotalDuration(time.Second)

	monitorData := reporter.getData()
	assert.Equal(monitorData.Interfaces[0].Total, int64(1))
	assert.Equal(monitorData.Interfaces[0].Failure, int64(1))
	assert.Equal(monitorData.Interfaces[0].ShortCircuited, int64(1))

}

func TestCseMonitorClient_PostMetrics(t *testing.T) {
	initEnv()
	assert := assert.New(t)
	config.SelfServiceName = "testService"
	reporter := NewReporter(metrics.DefaultRegistry, "127.0.0.1:8080", http.Header{"Content-Type": []string{"application/json"}}, time.Second, &tls.Config{})
	cseMonitClient := NewCseMonitorClient(http.Header{"Content-Type": []string{"application/json"}}, "http://127.0.0.1:9098", &tls.Config{})
	assert.Equal(cseMonitClient.URL, "http://127.0.0.1:9098")
	assert.Equal(cseMonitClient.Header, http.Header{"Content-Type": []string{"application/json"}})
	config.GlobalDefinition.Cse.Monitor.Client.Enable = false
	err := cseMonitClient.PostMetrics(reporter.getData())
	assert.NotNil(err)

	config.GlobalDefinition.Cse.Monitor.Client.Enable = true

	err = cseMonitClient.PostMetrics(MonitorData{
		Name:     "testService",
		Instance: "BLRY23283",
	})
	assert.NotNil(err)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(w, "Hello client")
	}))
	defer ts.Close()
	cseMonitClient = NewCseMonitorClient(http.Header{"Content-Type": []string{"application/json"}}, ts.URL, &tls.Config{})
	err = cseMonitClient.PostMetrics(MonitorData{
		Name:     "testService",
		Instance: "BLRY23283",
	})
	assert.Nil(err)
	ts1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Hello client")
	}))
	defer ts1.Close()
	cseMonitClient = NewCseMonitorClient(http.Header{"Content-Type": []string{"application/json"}}, ts1.URL, &tls.Config{})
	err = cseMonitClient.PostMetrics(MonitorData{
		Name:     "testService",
		Instance: "BLRY23283",
	})
	assert.NotNil(err)
}
