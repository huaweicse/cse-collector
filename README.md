### Metric Collector for Go-Chassis
[![Build Status](https://travis-ci.org/ServiceComb/cse-collector.svg?branch=master)](https://travis-ci.org/ServiceComb/cse-collector)   
This is a reporter plugin for go-chassis 
which report circuit breaker metrics to Huaweicloud.

# How to use 

in main.go
```go
import _ "github.com/huaweicse/cse-collector"
```

# Introdction
The metrics reported by this collector is listed below:
```
attempts
errors
successes
failures
rejects
shortCircuits
timeouts
fallbackSuccesses
fallbackFailures
totalDuration
runDuration
```
It also collects data for each api's:
```
Name                 string  `json:"name"`
Desc                 string  `json:"desc"`
Qps                  float64 `json:"qps"`
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
```

#### Configurations
in chassis.yaml file. 
You need to configure your micro service to send the data to 
Huaweicloud ServiceStage.

```
cse:
  monitor:
    client:
      serverUri: https://cse.cn-north-1.myhwclouds.com:443
      enable: true
```

#### Data Post

Every 2 sec data will be posted to monitoring server if serverUri is correct and enable is true

The data format
```
{
  "data": {
    "appId": "default",
    "version": "1.0.0",
    "name": "order",
    "environment": "",
    "instance": "order-c3bbef-8457f585dc-c69cz",
    "thread": 14,
    "customs": null,
    "interfaces": [
      {
        "name": "Consumer.restaurant.rest",
        "desc": "Consumer.restaurant.rest",
        "qps": 0,
        "latency": 1,
        "l995": 0,
        "l99": 0,
        "l90": 0,
        "l75": 0,
        "l50": 0,
        "l25": 0,
        "l5": 0,
        "rate": 1,
        "total": 0,
        "failure": 0,
        "shortCircuited": 0,
        "circuitBreakerOpen": false,
        "semaphoreRejected": 0,
        "threadPoolRejected": 0,
        "countTimeout": 0,
        "failureRate": 0
      }
    ],
    "cpu": 4,
    "memory": {
      "heapAlloc": 3041688,
      "heapIdle": 61628416,
      "heapInUse": 4399104,
      "heapObjects": 28501,
      "heapReleased": 61595648,
      "heapSys": 66027520
    },
    "serviceId": "32370e6ed839834a0493dddc3e7a7d5cb6d5db59",
    "instanceId": "7c91fbb922ef11ea99b30255ac1004a6"
  }
}

```

#### Data Flush

Every 10sec data will be flushed
```
[
    [
        {
            "time": 1526905933600,
            "appId": "default",
            "version": "0.0.1",
            "qps": 0,
            "latency": 0,
            "failureRate": 0,
            "total": 0,
            "breakerRateAgg": 0,
            "circuitBreakerOpen": false,
            "failure": 0,
            "shortCircuited": 0,
            "semaphoreRejected": 0,
            "threadPoolRejected": 0,
            "countTimeout": 0,
            "l995": 0,
            "l99": 0,
            "l90": 0,
            "l75": 0,
            "l50": 0,
            "l25": 0,
            "l5": 0,
            "instanceId": "6a0895085cf211e8bb850255ac105551",
            "thread": 11,
            "cpu": 4,
            "memory": {
                "heapAlloc": 2649632,
                "heapIdle": 1597440,
                "heapInUse": 3874816,
                "heapObjects": 24737,
                "heapReleased": 0,
                "heapSys": 5472256
            },
            "functionCount": 1,
            "customs": null,
            "name": "root1-ThinkPad-T440p"
        }
    ]
]
```
