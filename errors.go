package blindpay

import (
	"github.com/blindpaylabs/blindpay-go/internal/request"
)

// APIError represents an error response from the BlindPay API.
// It is an alias to the internal request.APIError for convenience.
type APIError = request.APIError

// ErrorItem represents an individual error in the errors array.
// It is an alias to the internal request.ErrorItem for convenience.
type ErrorItem = request.ErrorItem

// NetworkError is returned when no HTTP response was received (connection
// failure, timeout, cancellation). It unwraps to the underlying error.
type NetworkError = request.NetworkError

// ErrInvalidIdempotencyKey is returned, without sending the request, when a
// caller-supplied idempotency key is empty or longer than 255 characters.
var ErrInvalidIdempotencyKey = request.ErrInvalidIdempotencyKey

// ErrInvalidTimeout is returned, without sending the request, when a
// per-request timeout is zero or negative.
var ErrInvalidTimeout = request.ErrInvalidTimeout

// DecodeError is returned when a 2xx response body cannot be decoded. It
// keeps the status, request ID and idempotency key.
type DecodeError = request.DecodeError
