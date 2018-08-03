package metricsink

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/go-chassis/go-chassis/pkg/httpclient"
)

// constant for the cse-collector
const (
	PostRemoteTimeout = 20 * time.Second
	IdleConnsPerHost  = 100
	DefaultTimeout    = time.Second * 20
	MetricsPath       = "/csemonitor/metric"
	EnvProjectID      = "CSE_PROJECT_ID"
)

// variables for cse-collector
var (
	MetricServerPath = ""
)

// CseMonitorClient is an object for storing client information
type CseMonitorClient struct {
	Header http.Header
	URL    string
	Client *httpclient.URLClient
}

// NewCseMonitorClient creates an new client for monitoring
func NewCseMonitorClient(header http.Header, url string, tlsConfig *tls.Config, version string) (*CseMonitorClient, error) {
	var apiVersion string

	switch version {
	case "v1":
		apiVersion = "v1"
	case "V1":
		apiVersion = "v1"
	case "v2":
		apiVersion = "v2"
	case "V2":
		apiVersion = "v2"
	default:
		apiVersion = "v2"
	}
	//Update the API Base Path based on the Version
	updateAPIPath(apiVersion)

	c, err := httpclient.GetURLClient(&httpclient.URLClientOption{
		SSLEnabled:            tlsConfig != nil,
		TLSConfig:             tlsConfig,
		ResponseHeaderTimeout: DefaultTimeout,
	})
	if err != nil {
		return nil, err
	}
	return &CseMonitorClient{
		Header: header,
		URL:    url,
		Client: c,
	}, nil
}

// updateAPIPath Update the Base PATH and HEADERS Based on the version of MetricServer used.
func updateAPIPath(apiVersion string) {

	//Check for the env Name in Container to get Domain Name
	//Default value is  "default"
	projectID, isExsist := os.LookupEnv(EnvProjectID)
	if !isExsist {
		projectID = "default"
	}
	switch apiVersion {
	case "v2":
		MetricServerPath = "/v2/" + projectID + MetricsPath

	case "v1":
		MetricServerPath = "/csemonitor/v1/metric"

	default:
		MetricServerPath = "/v2/" + projectID + MetricsPath
	}
}

// PostMetrics is a functions which sends the monintoring data to monitoring Server
func (cseMonitorClient *CseMonitorClient) PostMetrics(monitorData MonitorData) (err error) {
	var (
		js      []byte
		resp    *http.Response
		postURL string
		h       http.Header
	)

	if js, err = json.Marshal(monitorData); err != nil {
		return
	}

	postURL = cseMonitorClient.URL + MetricServerPath + "?service=" + monitorData.Name
	h = make(http.Header)
	for k, v := range cseMonitorClient.Header {
		h[k] = v
	}

	if resp, err = cseMonitorClient.Client.HTTPDo(http.MethodPost, postURL, h, js); err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var body []byte
		if body, err = ioutil.ReadAll(resp.Body); err != nil {
			body = []byte(fmt.Sprintf("(could not fetch response body for error: %s)", err))
		}
		err = fmt.Errorf("Unable to post to csemonitor: %d %s %s", resp.StatusCode, resp.Status, string(body))
	}
	return
}

// TransportFor creates an transport object with TLS information
func TransportFor(tlsconfig *tls.Config) http.RoundTripper {

	return &http.Transport{

		TLSClientConfig:     tlsconfig,
		MaxIdleConnsPerHost: IdleConnsPerHost,
		Dial: (&net.Dialer{
			Timeout:   10 * time.Second,
			KeepAlive: 30 * time.Second,
		}).Dial,
		ResponseHeaderTimeout: PostRemoteTimeout,
	}
}
