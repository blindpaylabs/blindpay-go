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
type APIError struct {
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Errors     []ErrorItem `json:"errors,omitempty"`
	TraceID    string      `json:"trace_id,omitempty"`
	RawBody    []byte      `json:"-"`
}

// ErrorItem represents an individual error in the errors array.
type ErrorItem struct {
	Code        string `json:"code,omitempty"`
	Message     string `json:"message"`
	LongMessage string `json:"long_message,omitempty"`
}

// Error implements the error interface for APIError.
func (e *APIError) Error() string {
	msg := fmt.Sprintf("blindpay: API error (status %d)", e.StatusCode)

	if e.Message != "" {
		msg += fmt.Sprintf(": %s", e.Message)
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
		return zero, parseAPIError(resp.StatusCode, respBody)
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

// parseAPIError attempts to parse an API error from the response body.
func parseAPIError(statusCode int, body []byte) *APIError {
	apiErr := &APIError{
		StatusCode: statusCode,
		RawBody:    body,
	}

	if err := json.Unmarshal(body, apiErr); err != nil {
		apiErr.Message = fmt.Sprintf("HTTP %d error", statusCode)
		if len(body) > 0 && len(body) < 1000 {
			apiErr.Message += fmt.Sprintf(": %s", string(body))
		}
	}

	return apiErr
}
