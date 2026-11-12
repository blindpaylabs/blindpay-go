package request

// ErrorCode is a stable, machine-readable error code from the API's error
// table (APIError.Code). The set is open: a code the SDK does not know yet
// is still a valid ErrorCode, it just has no ErrorCodeInfo.
type ErrorCode string

// RetryClass is the API's retry guidance for an error code.
type RetryClass string

const (
	// RetryClassRetryable: the same request can be retried as is.
	RetryClassRetryable RetryClass = "retryable"
	// RetryClassTerminal: do not retry.
	RetryClassTerminal RetryClass = "terminal"
	// RetryClassActionable: the request succeeds only after you fix
	// something (fund a balance, correct bank details...).
	RetryClassActionable RetryClass = "actionable"
	// RetryClassUnknown: the error has no code, or a code this SDK version
	// does not know.
	RetryClassUnknown RetryClass = "unknown"
)

// ErrorCodeInfo is the error table's metadata for one code.
type ErrorCodeInfo struct {
	// Status is the HTTP status the code typically responds with.
	Status int
	// Retry is the retry guidance.
	Retry RetryClass
}

// LookupErrorCode returns the error table's metadata for code. ok is false
// for an empty or unknown code.
func LookupErrorCode[C ~string](code C) (ErrorCodeInfo, bool) {
	info, ok := errorCodes[ErrorCode(code)]
	return info, ok
}

// RetryClass returns the error table's retry guidance for e.Code, or
// RetryClassUnknown when the response had no code or the code is unknown to
// this SDK version. It is guidance for the caller; the SDK's own automatic
// retries follow the transport rules documented on WithMaxRetries, and
// Retryable reports whether this failure is transient.
func (e *APIError) RetryClass() RetryClass {
	if info, ok := LookupErrorCode(e.Code); ok {
		return info.Retry
	}
	return RetryClassUnknown
}
