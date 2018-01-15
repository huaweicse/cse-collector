### Metric Collector for Go-Chassis
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

