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
