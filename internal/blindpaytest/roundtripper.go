package blindpaytest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type RoundTripper struct {
	T          *testing.T
	In         json.RawMessage
	Out        json.RawMessage
	Method     string
	Path       string
	Status     int
	StatusText string
	Header     http.Header
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if rt.Method != "" && req.Method != rt.Method {
		rt.T.Errorf("expected method %s, got %s", rt.Method, req.Method)
	}

	if rt.Path != "" && req.URL.Path != rt.Path {
		rt.T.Errorf("expected path %s, got %s", rt.Path, req.URL.Path)
	}

	if rt.In != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			rt.T.Fatalf("failed to read request body: %v", err)
		}

		var expected, actual any
		if err := json.Unmarshal(rt.In, &expected); err != nil {
			rt.T.Fatalf("failed to unmarshal expected body: %v", err)
		}
		if err := json.Unmarshal(body, &actual); err != nil {
			rt.T.Fatalf("failed to unmarshal actual body: %v", err)
		}

		expectedJSON, _ := json.Marshal(expected)
		actualJSON, _ := json.Marshal(actual)

		if string(expectedJSON) != string(actualJSON) {
			rt.T.Errorf("request body mismatch:\nexpected: %s\nactual:   %s", expectedJSON, actualJSON)
		}
	}

	status := rt.Status
	if status == 0 {
		status = http.StatusOK
	}

	statusText := rt.StatusText
	if statusText == "" {
		statusText = http.StatusText(status)
	}

	header := rt.Header
	if header == nil {
		header = http.Header{}
	}
	if header.Get("Content-Type") == "" {
		header.Set("Content-Type", "application/json")
	}

	return &http.Response{
		StatusCode: status,
		Status:     statusText,
		Header:     header,
		Body:       io.NopCloser(bytes.NewReader(rt.Out)),
		Request:    req,
	}, nil
}
