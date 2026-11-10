package request

import (
	"net/http"
	"time"
)

// RequestOptions holds per-request settings. Build it only through
// RequestOption values (see the public option package).
type RequestOptions struct {
	// IdempotencyKey, when non-nil, is sent as Idempotency-Key on
	// POST/PUT/PATCH/DELETE JSON requests. It must be 1-255 characters.
	IdempotencyKey *string
	// Header holds extra request headers.
	Header http.Header
	// Timeout overrides the client's per-attempt timeout when > 0.
	Timeout time.Duration
}

// RequestOption configures a single request.
type RequestOption func(*RequestOptions)

// WithIdempotencyKey sets the Idempotency-Key header.
func WithIdempotencyKey(key string) RequestOption {
	return func(o *RequestOptions) {
		o.IdempotencyKey = &key
	}
}

// WithHeader adds a request header. Repeated calls for the same name add
// values.
func WithHeader(name, value string) RequestOption {
	return func(o *RequestOptions) {
		if o.Header == nil {
			o.Header = http.Header{}
		}
		o.Header.Add(name, value)
	}
}

// WithTimeout overrides the client's per-attempt timeout for one request.
func WithTimeout(d time.Duration) RequestOption {
	return func(o *RequestOptions) {
		o.Timeout = d
	}
}

func applyOptions(opts []RequestOption) RequestOptions {
	var o RequestOptions
	for _, opt := range opts {
		if opt != nil {
			opt(&o)
		}
	}
	return o
}
