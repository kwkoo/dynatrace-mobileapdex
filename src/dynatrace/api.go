package dynatrace

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// API contains reusable parameters like the server URL prefix and the authentication token.
type API struct {
	serverURL string
	token     string
	client    *http.Client
}

var tmpl *template.Template

const reportPostBody = `{
	"displayName" : "{{.DeviceDisplayName}}",
	"ipAddresses" : ["{{.IPAddress}}"],
	"type" : "{{.DeviceType}}",
	"favicon" : "https://dt-cdn.net/assets/images/brand/dynatrace-logo-33a874730e.svg",
	"series" : [
	  {
		"timeseriesId" : "custom:{{.MetricName}}",
		"dataPoints" : [ [ {{.Timestamp}}  , {{.Value}} ] ]
	  }
	]
}`

type DataPoint struct {
	CustomDeviceID    string
	IPAddress         string
	DeviceDisplayName string
	MetricDisplayName string
	DeviceType        string
	MetricName        string
	Timestamp         int64
	Value             float32
}

// NewAPI creates new Api object and initializes it.
func NewAPI(s, t string) API {
	log.Println("API initialized with server URL", s, "and API token", t)
	return API{
		serverURL: s,
		token:     t,
		client:    &http.Client{},
	}
}

// RegisterCustomMetric registers a custom metric.
func (api API) RegisterCustomMetric(name, displayName, unit, metricType string) error {
	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	writer.WriteString("{\"displayName\":\"")
	writer.Flush()
	escapeJSONString(&buf, displayName)
	writer.WriteString("\",\"unit\":\"")
	writer.Flush()
	escapeJSONString(&buf, unit)
	writer.WriteString("\",\"types\":[\"")
	writer.Flush()
	escapeJSONString(&buf, metricType)
	writer.WriteString("\"]}")
	writer.Flush()

	resp, err := api.ServerRequest("PUT", "/api/v1/timeseries/custom:"+name+"/", &buf)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	respBody := string(b[:])
	if strings.HasPrefix(respBody, "{\"error\"") {
		return errors.New("Received error response while trying to create custom metric: " + respBody)
	}
	return nil
}

// CustomMetricExists checks if a certain custom metric exists.
func (api API) CustomMetricExists(m string) (bool, error) {
	m = "custom:" + m
	metrics, err := api.GetCustomMetrics()
	if err != nil {
		return false, err
	}
	for _, v := range metrics {
		//log.Println("Checking", v, "against", m)
		if v == m {
			return true, nil
		}
	}
	return false, nil
}

// GetCustomMetrics returns a slice of all custom metrics.
func (api API) GetCustomMetrics() ([]string, error) {
	resp, err := api.ServerRequest("GET", "/api/v1/timeseries?filter=CUSTOM", nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New("server returned " + strconv.Itoa(resp.StatusCode) + " status code for GetCustomMetrics")
	}

	defer resp.Body.Close()

	dec := json.NewDecoder(resp.Body)
	metrics := make([]string, 0, 1)
	state := 0
	for {
		t, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		switch state {
		case 0:
			if v, ok := t.(json.Delim); ok && v.String() == "{" {
				state = 1
			}
		case 1:
			if v, ok := t.(string); ok && v == "timeseriesId" {
				state = 2
			} else {
				state = 0
			}
		case 2:
			if v, ok := t.(string); ok {
				metrics = append(metrics, v)
			}
			state = 0
		}
	}

	return metrics, nil
}

// ServerRequest makes an http request to the server and returns an http response.
func (api API) ServerRequest(method, uri string, body io.Reader) (*http.Response, error) {
	targetURL := api.serverURL + uri
	req, err := http.NewRequest(method, targetURL, body)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Api-Token "+api.token)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := api.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (api API) ReportDataPoint(dp DataPoint) (string, error) {
	if tmpl == nil {
		tmpl, _ = template.New("body").Parse(reportPostBody)
	}

	var buf bytes.Buffer
	writer := bufio.NewWriter(&buf)
	if err := tmpl.Execute(writer, dp); err != nil {
		return "", fmt.Errorf("error forming template: %v", err)
	}
	writer.Flush()

	resp, err := api.ServerRequest("POST", "/api/v1/entity/infrastructure/custom/"+dp.CustomDeviceID, &buf)
	if err != nil {
		return "", fmt.Errorf("error while trying to report custom metric: %v", err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error received while reading body: %v", err)
	}
	respBody := string(b[:])
	return respBody, nil
}

func escapeJSONString(buf *bytes.Buffer, s string) {
	json.HTMLEscape(buf, []byte(s))
}
