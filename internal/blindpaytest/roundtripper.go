package blindpaytest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"testing"
)

// RoundTripper is a mock http.RoundTripper that asserts on the outgoing
// request and replies with a canned response.
//
// Request-side assertions (each is skipped when left at its zero value):
//   - Method, Path: exact match on req.Method / req.URL.Path.
//   - Query: exact match on the parsed query string. A non-nil empty
//     url.Values{} asserts that no query string is sent at all.
//   - RequestHeader: every listed header must be present with exactly these
//     values. Headers not listed are not checked.
//   - NoRequestHeader: every listed header must be absent.
//   - In: the JSON request body, compared structurally.
//
// Response-side fields: Out (body), Status (default 200), StatusText and
// Header (response headers; Content-Type defaults to application/json).
type RoundTripper struct {
	T               *testing.T
	In              json.RawMessage
	Out             json.RawMessage
	Method          string
	Path            string
	Query           url.Values
	RequestHeader   http.Header
	NoRequestHeader []string
	Status          int
	StatusText      string
	Header          http.Header
}

func (rt *RoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	rt.T.Helper()

	if rt.Method != "" && req.Method != rt.Method {
		rt.T.Errorf("expected method %s, got %s", rt.Method, req.Method)
	}

	if rt.Path != "" && req.URL.Path != rt.Path {
		rt.T.Errorf("expected path %s, got %s", rt.Path, req.URL.Path)
	}

	if rt.Query != nil {
		got := req.URL.Query()
		if len(rt.Query) == 0 {
			if req.URL.RawQuery != "" {
				rt.T.Errorf("expected no query string, got %q", req.URL.RawQuery)
			}
		} else if !reflect.DeepEqual(map[string][]string(rt.Query), map[string][]string(got)) {
			rt.T.Errorf("query mismatch:\nexpected: %s\nactual:   %s", rt.Query.Encode(), got.Encode())
		}
	}

	for name, want := range rt.RequestHeader {
		got := req.Header.Values(name)
		if !reflect.DeepEqual(want, got) {
			rt.T.Errorf("request header %s mismatch: expected %q, got %q", name, want, got)
		}
	}

	for _, name := range rt.NoRequestHeader {
		if got := req.Header.Values(name); len(got) > 0 {
			rt.T.Errorf("expected request header %s to be absent, got %q", name, got)
		}
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
