package metricsink

import (
	"github.com/rcrowley/go-metrics"
	"runtime"
	"runtime/pprof"
	"strings"
	"time"
)

var threadCreateProfile = pprof.Lookup("threadcreate")

// MonitorData is an object which stores the monitoring information for an application
type MonitorData struct {
	AppID      string                 `json:"appId"`
	Version    string                 `json:"version"`
	Name       string                 `json:"name"`
	Instance   string                 `json:"instance"`
	Thread     int                    `json:"thread"`
	Customs    map[string]interface{} `json:"customs"` // ?
	Interfaces []*InterfaceInfo       `json:"interfaces"`
	CPU        float64                `json:"cpu"`
	Memory     map[string]interface{} `json:"memory"`
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
	IsCircuitBreakerOpen bool    `json:"isCircuitBreakerOpen"`
	SemaphoreRejected    int64   `json:"semaphoreRejected"`
	ThreadPoolRejected   int64   `json:"threadPoolRejected"`
	CountTimeout         int64   `json:"countTimeout"`
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
		case "failures":
			interfaceInfo.Failure = metric.Count()
		case "shortCircuits":
			interfaceInfo.ShortCircuited = metric.Count()
		case "successes":
			interfaceInfo.successCount = metric.Count()
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
			interfaceInfo.QPS = t.RateMean()
		}

	}
	if interfaceInfo.Total == 0 {
		interfaceInfo.Rate = 100
	} else {
		interfaceInfo.Rate = float64(interfaceInfo.successCount) / float64(interfaceInfo.Total)
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
