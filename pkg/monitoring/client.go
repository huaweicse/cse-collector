package monitoring

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/go-chassis/foundation/httpclient"
)

// constant for the cse-collector
const (
	DefaultTimeout = time.Second * 20
	MetricsPath    = "/csemonitor/metric"
	EnvProjectID   = "CSE_PROJECT_ID"
)

// variables for cse-collector
var (
	// /v2/{project}/csemonitor/metric
	MetricServerPath = ""
)

// CseMonitorClient is an object for storing client information
type CseMonitorClient struct {
	Header http.Header
	URL    string
	Client *httpclient.URLClient
}

// NewCseMonitorClient creates an new client for monitoring
func NewCseMonitorClient(header http.Header, url string, tlsConfig *tls.Config) (*CseMonitorClient, error) {
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
func updateAPIPath() {
	//Check for the env Name in Container to get Domain Name
	//Default value is  "default"
	projectID, isExist := os.LookupEnv(EnvProjectID)
	if !isExist {
		projectID = "default"
	}
	MetricServerPath = "/v2/" + projectID + MetricsPath
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
		err = fmt.Errorf("can't post to csemonitor: %d %s %s", resp.StatusCode, resp.Status, string(body))
	}
	return
}
