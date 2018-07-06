package metricsink

import (
	"fmt"
	"math"
	"runtime"
	"runtime/pprof"
	"strconv"
	"strings"
	"time"

	"github.com/ServiceComb/go-chassis/core/lager"
	"github.com/ServiceComb/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/rcrowley/go-metrics"
)

var threadCreateProfile = pprof.Lookup("threadcreate")

// MonitorData is an object which stores the monitoring information for an application
type MonitorData struct {
	AppID       string                 `json:"appId"`
	Version     string                 `json:"version"`
	Name        string                 `json:"name"`
	Environment string                 `json:"environment"`
	Instance    string                 `json:"instance"`
	Thread      int                    `json:"thread"`
	Customs     map[string]interface{} `json:"customs"` // ?
	Interfaces  []*InterfaceInfo       `json:"interfaces"`
	CPU         float64                `json:"cpu"`
	Memory      map[string]interface{} `json:"memory"`
	ServiceID   string                 `json:"serviceId"`
	InstanceID  string                 `json:"instanceId"`
}

// InterfaceInfo is an object which store the monitoring information of a particular interface
type InterfaceInfo struct {
	Name                 string  `json:"name"`
	Desc                 string  `json:"desc"`
	QPS                  float64 `json:"qps"`
	Latency              int     `json:"latency"`
	L995                 int     `json:"l995"`
	L99                  int     `json:"l99"`
	L90                  int     `json:"l90"`
	L75                  int     `json:"l75"`
	L50                  int     `json:"l50"`
	L25                  int     `json:"l25"`
	L5                   int     `json:"l5"`
	Rate                 float64 `json:"rate"`
	Total                int64   `json:"total"`
	Failure              int64   `json:"failure"`
	ShortCircuited       int64   `json:"shortCircuited"`
	IsCircuitBreakerOpen bool    `json:"circuitBreakerOpen"`
	SemaphoreRejected    int64   `json:"semaphoreRejected"`
	ThreadPoolRejected   int64   `json:"threadPoolRejected"`
	CountTimeout         int64   `json:"countTimeout"`
	FailureRate          float64 `json:"failureRate"`
	successCount         int64
}

// NewMonitorData creates a new monitoring object
func NewMonitorData() *MonitorData {
	monitorData := new(MonitorData)
	monitorData.Interfaces = make([]*InterfaceInfo, 0)
	return monitorData
}

func (monitorData *MonitorData) getOrCreateInterfaceInfo(name string) *InterfaceInfo {
	for _, interfaceInfo := range monitorData.Interfaces {
		if interfaceInfo.Name == name {
			return interfaceInfo
		}
	}
	interfaceInfo := new(InterfaceInfo)
	interfaceInfo.Name = name
	interfaceInfo.Desc = name
	monitorData.Interfaces = append(monitorData.Interfaces, interfaceInfo)
	return interfaceInfo
}

func (monitorData *MonitorData) appendInterfaceInfo(name string, i interface{}) {
	var interfaceInfo = monitorData.getOrCreateInterfaceInfo(getInterfaceName(name))
	switch metric := i.(type) {
	case metrics.Counter:
		switch getEventType(name) {
		case "attempts":
			interfaceInfo.Total = metric.Count()
		case "errors":
			interfaceInfo.Failure = metric.Count()
		case "shortCircuits":
			interfaceInfo.ShortCircuited = metric.Count()
		case "successes":
			interfaceInfo.successCount = metric.Count()
		}

		if isCBOpen, err := hystrix.IsCircuitBreakerOpen(getInterfaceName(name)); err != nil {
			interfaceInfo.IsCircuitBreakerOpen = false
		} else {
			interfaceInfo.IsCircuitBreakerOpen = isCBOpen
		}

		qps := (float64(interfaceInfo.Total) * (1 - math.Exp(-5.0/60.0/1)))
		movingAverageFor3Precision, err := strconv.ParseFloat(fmt.Sprintf("%.3f", qps), 64)
		if err == nil {
			interfaceInfo.QPS = movingAverageFor3Precision
		} else {
			interfaceInfo.QPS = 0
		}

	case metrics.Timer:
		t := metric.Snapshot()
		ps := t.Percentiles([]float64{0.05, 0.25, 0.5, 0.75, 0.90, 0.99, 0.995})
		switch getEventType(name) {
		case "runDuration":
			interfaceInfo.L5 = int(ps[0] / float64(time.Millisecond))
			interfaceInfo.L25 = int(ps[1] / float64(time.Millisecond))
			interfaceInfo.L50 = int(ps[2] / float64(time.Millisecond))
			interfaceInfo.L75 = int(ps[3] / float64(time.Millisecond))
			interfaceInfo.L90 = int(ps[4] / float64(time.Millisecond))
			interfaceInfo.L99 = int(ps[5] / float64(time.Millisecond))
			interfaceInfo.L995 = int(ps[6] / float64(time.Millisecond))
			interfaceInfo.Latency = int(t.Mean() / float64(time.Millisecond))
		}
	}
	interfaceInfo.Rate = 1 //rate is no use any more and must be set to 1
	if interfaceInfo.Total == 0 {
		interfaceInfo.FailureRate = 0
	} else {
		totalErrorCount := interfaceInfo.Failure + interfaceInfo.SemaphoreRejected + interfaceInfo.ThreadPoolRejected + interfaceInfo.CountTimeout
		if totalErrorCount == 0 {
			interfaceInfo.FailureRate = 0
		} else {
			failureRate, err := strconv.ParseFloat(fmt.Sprintf("%.3f", float64(totalErrorCount)/float64(interfaceInfo.Total)), 64)
			if err == nil && failureRate > 0 {
				interfaceInfo.FailureRate = failureRate
			} else {
				lager.Logger.Warnf("Error in calculating the failureRate %v, default value(0) is assigned to failureRate", failureRate)
				interfaceInfo.FailureRate = 0
			}
		}
	}
}

func getInterfaceName(metricName string) string {
	command := strings.Split(metricName, ".")
	return strings.Join(command[:len(command)-1], ".")

}

func getEventType(metricName string) string {
	command := strings.Split(metricName, ".")
	return command[len(command)-1]
}

func getProcessInfo() map[string]interface{} {
	var memoryInfo = make(map[string]interface{})
	var memStats = runtime.MemStats{}
	runtime.ReadMemStats(&memStats)
	memoryInfo["heapAlloc"] = memStats.HeapAlloc
	memoryInfo["heapSys"] = memStats.HeapSys
	memoryInfo["heapIdle"] = memStats.HeapIdle
	memoryInfo["heapInUse"] = memStats.HeapInuse
	memoryInfo["heapReleased"] = memStats.HeapReleased
	memoryInfo["heapObjects"] = memStats.HeapObjects
	return memoryInfo
}
