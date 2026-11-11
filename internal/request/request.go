package request

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Config holds the configuration for making API requests.
type Config struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
	UserAgent  string
}

// APIError represents an error response from the BlindPay API.
//
// Only Message is guaranteed by the API. Code is the stable, UPPERCASE
// machine-readable key (e.g. "PAYOUTS_INSUFFICIENT_BALANCE") and is the field
// to branch on; it is empty for idempotency errors and for non-JSON bodies.
// Description is a human-readable explanation safe to show to end users.
// TraceID is absent on auth and rate-limit errors.
type APIError struct {
	StatusCode  int         `json:"-"`
	Code        string      `json:"code,omitempty"`
	Message     string      `json:"message"`
	Description string      `json:"description,omitempty"`
	Errors      []ErrorItem `json:"errors,omitempty"`
	TraceID     string      `json:"trace_id,omitempty"`
	RawBody     []byte      `json:"-"`
}

// ErrorItem represents an individual entry of the errors array.
//
// The API sends the entries as strings (the spec's `errors: string[]`) but may
// forward arbitrary JSON values. A string entry is stored in Message. An
// object entry with string code/message/long_message members fills those
// fields. Any other entry (numbers, booleans, null, objects without a string
// message) keeps its raw JSON text in Message, so no detail is lost and
// decoding never fails.
type ErrorItem struct {
	Code        string `json:"code,omitempty"`
	Message     string `json:"message"`
	LongMessage string `json:"long_message,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler. It never returns an error.
func (e *ErrorItem) UnmarshalJSON(data []byte) error {
	*e = ErrorItem{}
	data = bytes.TrimSpace(data)

	// json.Unmarshal treats null as a no-op for strings and structs, so it
	// must be handled before the typed attempts below.
	if bytes.Equal(data, []byte("null")) {
		e.Message = "null"
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		e.Message = s
		return nil
	}

	var obj struct {
		Code        string `json:"code"`
		Message     string `json:"message"`
		LongMessage string `json:"long_message"`
	}
	if err := json.Unmarshal(data, &obj); err == nil && obj.Message != "" {
		e.Code = obj.Code
		e.Message = obj.Message
		e.LongMessage = obj.LongMessage
		return nil
	}

	e.Message = string(data)
	return nil
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	msg := fmt.Sprintf("blindpay: API error (status %d)", e.StatusCode)

	if e.Code != "" {
		msg += fmt.Sprintf(" %s", e.Code)
	}

	if e.Message != "" {
		msg += fmt.Sprintf(": %s", e.Message)
	}

	if e.Description != "" {
		msg += fmt.Sprintf(" (%s)", e.Description)
	}

	if e.TraceID != "" {
		msg += fmt.Sprintf(" [trace_id: %s]", e.TraceID)
	}

	if len(e.Errors) > 0 {
		msg += " - errors:"
		for _, err := range e.Errors {
			if err.Code != "" {
				msg += fmt.Sprintf(" [%s: %s]", err.Code, err.Message)
			} else {
				msg += fmt.Sprintf(" [%s]", err.Message)
			}
		}
	}

	return msg
}

// Do performs an HTTP request and decodes the response into the given type T.
func Do[T any](cfg *Config, ctx context.Context, method, path string, body any) (T, error) {
	var zero T

	url := cfg.BaseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return zero, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return zero, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)
	req.Header.Set("User-Agent", cfg.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := cfg.HTTPClient.Do(req)
	if err != nil {
		return zero, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return zero, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return zero, ParseAPIError(resp.StatusCode, respBody)
	}

	// For DELETE requests that return 204 No Content, return zero value
	if resp.StatusCode == http.StatusNoContent {
		return zero, nil
	}

	var result T
	if err := json.Unmarshal(respBody, &result); err != nil {
		return zero, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}

// ParseAPIError builds an APIError from a non-2xx response. A JSON object
// body is decoded field by field (ErrorItem decoding never fails, so the
// API's code and message always survive). Any other body (HTML from a proxy,
// empty, truncated JSON) yields a generic "HTTP <status> error" message that
// carries the status; the raw body is always kept in RawBody. Code paths
// that cannot use Do (multipart upload) must use it too.
func ParseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		RawBody:    body,
	}

	// Decode member by member so one oddly typed member (e.g. errors sent
	// as a string) cannot discard code, message or trace_id.
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(body, &fields); err != nil || fields == nil {
		apiErr.Message = fmt.Sprintf("HTTP %d error", statusCode)
		if len(body) > 0 && len(body) < 1000 {
			apiErr.Message += fmt.Sprintf(": %s", string(body))
		}
		return apiErr
	}

	apiErr.Code = stringMember(fields["code"])
	apiErr.Message = stringMember(fields["message"])
	apiErr.Description = stringMember(fields["description"])
	apiErr.TraceID = stringMember(fields["trace_id"])
	apiErr.Errors = errorItems(fields["errors"])

	if apiErr.Message == "" {
		apiErr.Message = fmt.Sprintf("HTTP %d error", statusCode)
	}

	return apiErr
}

// stringMember returns raw as a string, or "" when it is absent or not a
// JSON string.
func stringMember(raw json.RawMessage) string {
	var s string
	if len(raw) == 0 || json.Unmarshal(raw, &s) != nil {
		return ""
	}
	return s
}

// errorItems decodes the errors member. A JSON array decodes item by item
// (ErrorItem never fails); any other non-null value becomes a single item
// holding its raw text.
func errorItems(raw json.RawMessage) []ErrorItem {
	raw = bytes.TrimSpace(raw)
	if len(raw) == 0 || bytes.Equal(raw, []byte("null")) {
		return nil
	}

	var items []ErrorItem
	if err := json.Unmarshal(raw, &items); err == nil {
		return items
	}

	var item ErrorItem
	_ = item.UnmarshalJSON(raw)
	return []ErrorItem{item}
}
