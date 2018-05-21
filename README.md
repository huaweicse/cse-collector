### Metric Collector for Go-Chassis
[![Build Status](https://travis-ci.org/ServiceComb/cse-collector.svg?branch=master)](https://travis-ci.org/ServiceComb/cse-collector)
This a metric collector for Go-Chassis which collects metrics of the microservices. 
It can collect metrics for each api's exposed by the micro-services. The same data can be 
exposed to Huawei CSE Governance Dashboard.

The metrics collected by this collector is listed below:
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
You need to configure your microservice to send the data at regular interval to 
Huawei CSE Dashboard.

```
cse:
  monitor:
    client:
      serverUri: https://cse.cn-north-1.myhwclouds.com:443
      enable: true
```

#### Data Post

Every 2 sec data will be posted to monitoring server if serverUri is correct and enable is true

The data formate after post 
```
[
    [
        {
            "time": 1526905902632,
            "appId": "default",
            "version": "0.0.1",
            "qps": 0.16,
            "latency": 0,
            "failureRate": 0,
            "total": 2,
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
                "heapAlloc": 2294040,
                "heapIdle": 1671168,
                "heapInUse": 3768320,
                "heapObjects": 19609,
                "heapReleased": 0,
                "heapSys": 5439488
            },
            "functionCount": 1,
            "customs": null,
            "name": "root1-ThinkPad-T440p"
        }
    ]
]

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

If perticular micro service has more than one instance and if request has been sent to perticular instance then only for that instance all the value should change.
