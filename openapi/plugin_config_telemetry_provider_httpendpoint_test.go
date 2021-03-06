package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTelemetryProviderHttpEndpoint_Validate(t *testing.T) {
	testCases := []struct {
		testName    string
		url         string
		expectedErr error
	}{
		{
			testName:    "happy path - host and port populated",
			url:         "http://telemetry.myhost.com/v1/metrics",
			expectedErr: nil,
		},
		{
			testName:    "url is empty",
			url:         "",
			expectedErr: errors.New("http endpoint telemetry configuration is missing a value for the 'url property'"),
		},
		{
			testName:    "url is wrongly formatter",
			url:         "htop://something-wrong.com",
			expectedErr: errors.New("http endpoint telemetry configuration does not have a valid URL 'htop://something-wrong.com'"),
		},
	}

	for _, tc := range testCases {
		tpg := TelemetryProviderHTTPEndpoint{
			URL: tc.url,
		}
		err := tpg.Validate()
		assert.Equal(t, tc.expectedErr, err, tc.testName)
	}
}

func TestCreateNewCounterMetric(t *testing.T) {
	testCases := []struct {
		testName       string
		prefix         string
		expectedMetric telemetryMetric
	}{
		{
			testName:       "prefix is not empty",
			prefix:         "prefix",
			expectedMetric: telemetryMetric{metricTypeCounter, "prefix.metric_name"},
		},
		{
			testName:       "prefix is empty",
			prefix:         "",
			expectedMetric: telemetryMetric{metricTypeCounter, "metric_name"},
		},
	}

	for _, tc := range testCases {
		telemetryMetric := createNewCounterMetric(tc.prefix, "metric_name")
		assert.Equal(t, tc.expectedMetric, telemetryMetric, tc.testName)
	}
}

func TestCreateNewRequest(t *testing.T) {
	testCases := []struct {
		testName              string
		url                   string
		expectedCounterMetric telemetryMetric
		expectedContentType   string
		expectedUserAgent     string
		expectedErr           error
	}{
		{
			testName: "happy path - request is created with the expected Header and telemetryMetric",
			expectedCounterMetric: telemetryMetric{
				MetricType: metricTypeCounter,
				MetricName: "prefix.terraform.openapi_plugin_version.version.total_runs",
			},
			expectedContentType: "application/json",
			expectedUserAgent:   "OpenAPI Terraform Provider",
			expectedErr:         nil,
		},
		{
			testName:    "crappy path - bad url",
			url:         "&^%",
			expectedErr: errors.New("parse &^%: invalid URL escape \"%\""),
		},
	}

	for _, tc := range testCases {
		var err error
		var request *http.Request
		var reqBody []byte
		telemetryMetric := telemetryMetric{}
		tph := TelemetryProviderHTTPEndpoint{
			URL: tc.url,
		}

		request, err = tph.createNewRequest(tc.expectedCounterMetric)
		if tc.expectedErr == nil {
			reqBody, err = ioutil.ReadAll(request.Body)
			err = json.Unmarshal(reqBody, &telemetryMetric)
			assert.NoError(t, err, tc.testName)
			assert.Equal(t, tc.expectedContentType, request.Header.Get(contentType), tc.testName)
			assert.Contains(t, request.Header.Get(userAgentHeader), tc.expectedUserAgent, tc.testName)
			assert.Equal(t, tc.expectedCounterMetric, telemetryMetric, tc.testName)
		} else {
			assert.EqualError(t, err, tc.expectedErr.Error(), tc.testName)
		}
	}
}

func TestTelemetryProviderHttpEndpointSubmitMetric(t *testing.T) {
	testCases := []struct {
		testName             string
		returnedResponseCode int
		expectedErr          error
	}{
		{
			testName:             "happy path",
			returnedResponseCode: http.StatusOK,
			expectedErr:          nil,
		},
		{
			testName:             "api server returns non 2xx code",
			returnedResponseCode: http.StatusNotFound,
			expectedErr:          errors.New("/v1/metrics' returned a non expected status code 404"),
		},
	}

	for _, tc := range testCases {

		expectedCounterMetric := telemetryMetric{
			MetricType: metricTypeCounter,
			MetricName: "prefix.terraform.openapi_plugin_version.version.total_runs",
		}

		api := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(t, req.Method, http.MethodPost, tc.testName)
			assert.Equal(t, "/v1/metrics", req.URL.String(), tc.testName)
			assert.Equal(t, req.Header.Get(contentType), "application/json", tc.testName)
			assert.Contains(t, req.Header.Get(userAgentHeader), "OpenAPI Terraform Provider", tc.testName)
			reqBody, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err, tc.testName)
			telemetryMetric := telemetryMetric{}
			err = json.Unmarshal(reqBody, &telemetryMetric)
			assert.Nil(t, err, tc.testName)
			assert.Equal(t, expectedCounterMetric.MetricType, telemetryMetric.MetricType, tc.testName)
			assert.Equal(t, expectedCounterMetric.MetricName, telemetryMetric.MetricName, tc.testName)
			rw.WriteHeader(tc.returnedResponseCode)
		}))
		// Close the server when test finishes
		defer api.Close()

		tph := TelemetryProviderHTTPEndpoint{
			URL: fmt.Sprintf("%s/v1/metrics", api.URL),
		}
		err := tph.submitMetric(expectedCounterMetric)
		if tc.expectedErr == nil {
			assert.NoError(t, err, tc.testName)
		} else {
			assert.Error(t, err, tc.testName)
			assert.Contains(t, err.Error(), tc.expectedErr.Error(), tc.testName)
		}
	}
}

func TestTelemetryProviderHttpEndpointSubmitMetricFailureScenarios(t *testing.T) {
	testCases := []struct {
		testName    string
		inputURL    string
		expectedErr error
	}{
		{
			testName:    "url is missing the protocol",
			inputURL:    "?",
			expectedErr: errors.New("request POST ? failed. Response Error: 'Post ?: unsupported protocol scheme \"\"'"),
		},
		{
			testName:    "url contains invalid characters",
			inputURL:    "&^%",
			expectedErr: errors.New("parse &^%: invalid URL escape \"%\""),
		},
	}

	for _, tc := range testCases {
		tph := TelemetryProviderHTTPEndpoint{
			URL: tc.inputURL,
		}
		err := tph.submitMetric(telemetryMetric{metricTypeCounter, "prefix.terraform.openapi_plugin_version.version.total_runs"})
		assert.EqualError(t, err, tc.expectedErr.Error())
	}
}

func TestTelemetryProviderHttpEndpointIncOpenAPIPluginVersionTotalRunsCounter(t *testing.T) {
	testCases := []struct {
		testName             string
		returnedResponseCode int
		expectedErr          error
	}{
		{
			testName:             "happy path",
			returnedResponseCode: http.StatusOK,
			expectedErr:          nil,
		},
		{
			testName:             "metric submission fails",
			returnedResponseCode: http.StatusNotFound,
			expectedErr:          errors.New("/v1/metrics' returned a non expected status code 404"),
		},
	}

	for _, tc := range testCases {

		api := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			reqBody, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err, tc.testName)
			telemetryMetric := telemetryMetric{}
			err = json.Unmarshal(reqBody, &telemetryMetric)
			assert.Nil(t, err, tc.testName)
			assert.Equal(t, metricTypeCounter, telemetryMetric.MetricType, tc.testName)
			assert.Equal(t, "terraform.openapi_plugin_version.0_26_0.total_runs", telemetryMetric.MetricName, tc.testName)
			rw.WriteHeader(tc.returnedResponseCode)
		}))
		// Close the server when test finishes
		defer api.Close()

		tph := TelemetryProviderHTTPEndpoint{
			URL: fmt.Sprintf("%s/v1/metrics", api.URL),
		}
		err := tph.IncOpenAPIPluginVersionTotalRunsCounter("0.26.0")
		if tc.expectedErr == nil {
			assert.NoError(t, err, tc.testName)
		} else {
			assert.Error(t, err, tc.testName)
			assert.Contains(t, err.Error(), tc.expectedErr.Error(), tc.testName)
		}
	}
}

func TestTelemetryProviderHttpEndpointIncServiceProviderTotalRunsCounter(t *testing.T) {
	testCases := []struct {
		testName             string
		returnedResponseCode int
		expectedErr          error
	}{
		{
			testName:             "happy path",
			returnedResponseCode: http.StatusOK,
			expectedErr:          nil,
		},
		{
			testName:             "metric submission fails",
			returnedResponseCode: http.StatusNotFound,
			expectedErr:          errors.New("/v1/metrics' returned a non expected status code 404"),
		},
	}

	for _, tc := range testCases {

		api := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			reqBody, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err, tc.testName)
			telemetryMetric := telemetryMetric{}
			err = json.Unmarshal(reqBody, &telemetryMetric)
			assert.Nil(t, err, tc.testName)
			assert.Equal(t, metricTypeCounter, telemetryMetric.MetricType, tc.testName)
			assert.Equal(t, "terraform.providers.cdn.total_runs", telemetryMetric.MetricName, tc.testName)
			rw.WriteHeader(tc.returnedResponseCode)
		}))
		// Close the server when test finishes
		defer api.Close()

		tph := TelemetryProviderHTTPEndpoint{
			URL: fmt.Sprintf("%s/v1/metrics", api.URL),
		}
		err := tph.IncServiceProviderTotalRunsCounter("cdn")
		if tc.expectedErr == nil {
			assert.NoError(t, err, tc.testName)
		} else {
			assert.Error(t, err, tc.testName)
			assert.Contains(t, err.Error(), tc.expectedErr.Error(), tc.testName)
		}
	}
}
