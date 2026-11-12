package blindpay

import "github.com/blindpaylabs/blindpay-go/internal/request"

// ErrorCode is a stable, machine-readable error code from the API's error
// table; APIError.Code carries it as a plain string. The ErrorCode*
// constants (errorcodes_gen.go, generated from the spec) list every code
// this SDK version knows. The set is open: an unknown code still decodes
// into APIError.Code and compares as a plain string.
//
//	var apiErr *blindpay.APIError
//	if errors.As(err, &apiErr) {
//		switch blindpay.ErrorCode(apiErr.Code) {
//		case blindpay.ErrorCodeTransfersSwapAmountBelowMinimum:
//			// swaps between different tokens need at least 100 USDC
//		}
//	}
type ErrorCode = request.ErrorCode

// RetryClass is the API's retry guidance for an error code: retryable,
// terminal, actionable, or unknown (no code, or a code this SDK version does
// not know). APIError.RetryClass returns it.
type RetryClass = request.RetryClass

const (
	// RetryClassRetryable means the same request can be retried as is.
	RetryClassRetryable = request.RetryClassRetryable
	// RetryClassTerminal means do not retry.
	RetryClassTerminal = request.RetryClassTerminal
	// RetryClassActionable means the request succeeds only after you fix
	// something (fund a balance, correct bank details...).
	RetryClassActionable = request.RetryClassActionable
	// RetryClassUnknown means the error has no code, or a code this SDK
	// version does not know.
	RetryClassUnknown = request.RetryClassUnknown
)

// ErrorCodeInfo is the error table's metadata for one code: the HTTP status
// it typically responds with and its retry guidance.
type ErrorCodeInfo = request.ErrorCodeInfo

// LookupErrorCode returns the error table's metadata for code, which may be
// an ErrorCode constant or APIError.Code. ok is false for an empty or
// unknown code.
func LookupErrorCode[C ~string](code C) (ErrorCodeInfo, bool) {
	return request.LookupErrorCode(code)
}
