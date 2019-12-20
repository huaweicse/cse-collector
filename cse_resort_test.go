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

package metricsink_test

import (
	"encoding/json"
	"github.com/go-chassis/go-chassis/third_party/forked/afex/hystrix-go/hystrix"
	"github.com/huaweicse/cse-collector/pkg/monitoring"
	"testing"
	"time"
)

func TestGetReporter(t *testing.T) {
	hystrix.Do("cmd", func() error {
		time.Sleep(20 * time.Millisecond)
		return nil
	}, nil)
	hystrix.Do("cmd", func() error {
		time.Sleep(10 * time.Millisecond)
		return nil
	}, nil)
	hystrix.Do("cmd", func() error {
		time.Sleep(2 * time.Millisecond)
		return nil
	}, nil)
	time.Sleep(1 * time.Second)
	data := &monitoring.MonitorData{}
	cb, _, _ := hystrix.GetCircuit("cmd")
	data.AppendInterfaceInfo(cb)
	b, _ := json.MarshalIndent(data, "", " ")
	t.Log(string(b))
}
