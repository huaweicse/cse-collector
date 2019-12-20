/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package monitoring

import (
	"fmt"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/go-mesh/openlogging"
	"math"
	"strconv"
	"strings"
	"time"
)

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

// NewMonitorData creates a new monitoring object
func NewMonitorData() *MonitorData {
	monitorData := new(MonitorData)
	monitorData.Interfaces = make([]*InterfaceInfo, 0)
	return monitorData
}

func (monitorData *MonitorData) getOrCreateInterfaceInfo(name string) *InterfaceInfo {
	interfaceName := GetInterfaceName(name)
	for _, interfaceInfo := range monitorData.Interfaces {
		if interfaceInfo.Name == interfaceName {
			return interfaceInfo
		}
	}
	interfaceInfo := new(InterfaceInfo)
	interfaceInfo.Name = interfaceName
	interfaceInfo.Desc = interfaceName
	monitorData.Interfaces = append(monitorData.Interfaces, interfaceInfo)
	return interfaceInfo
}
func GetInterfaceName(metricName string) string {
	command := strings.Split(metricName, ".")
	return strings.Join(command[:len(command)-1], ".")

}
func (monitorData *MonitorData) AppendInterfaceInfo(cb *hystrix.CircuitBreaker) {
	var interfaceInfo = monitorData.getOrCreateInterfaceInfo(cb.Name)
	now := time.Now()
	c := cb.Metrics.DefaultCollector()
	//attempts:
	interfaceInfo.Total = int64(c.NumRequests().Sum(now))
	//errors
	interfaceInfo.Failure = int64(c.Failures().Sum(now))
	//shortCircuits
	interfaceInfo.ShortCircuited = int64(c.ShortCircuits().Sum(now))
	//successes
	interfaceInfo.successCount = int64(c.Successes().Sum(now))

	interfaceInfo.IsCircuitBreakerOpen = cb.IsOpen()

	qps := float64(interfaceInfo.Total) * (1 - math.Exp(-5.0/60.0/1))
	movingAverageFor3Precision, err := strconv.ParseFloat(fmt.Sprintf("%.3f", qps), 64)
	if err == nil {
		interfaceInfo.QPS = movingAverageFor3Precision
	} else {
		interfaceInfo.QPS = 0
	}
	runDuration := c.RunDuration()
	interfaceInfo.L5 = int(runDuration.Percentile(5))
	interfaceInfo.L25 = int(runDuration.Percentile(25))
	interfaceInfo.L50 = int(runDuration.Percentile(5))
	interfaceInfo.L75 = int(runDuration.Percentile(75))
	interfaceInfo.L90 = int(runDuration.Percentile(90))
	interfaceInfo.L99 = int(runDuration.Percentile(99))
	interfaceInfo.L995 = int(runDuration.Percentile(99.5))
	interfaceInfo.Latency = int(runDuration.Mean())
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
				openlogging.GetLogger().Warnf("Error in calculating the failureRate %v, default value(0) is assigned to failureRate", failureRate)
				interfaceInfo.FailureRate = 0
			}
		}
	}
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
