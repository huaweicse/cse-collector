package metricsink

import (
	"runtime"
	"runtime/pprof"
	"strings"
)

var threadCreateProfile = pprof.Lookup("threadcreate")

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
