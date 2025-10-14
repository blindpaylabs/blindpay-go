package blindpay

import (
	"fmt"
)

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

// Error returns a string representation of the APIError.
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
