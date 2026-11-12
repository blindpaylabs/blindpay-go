package blindpay

import (
	"net/http"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/option"
)

// DefaultTimeout is the default per-attempt timeout.
const DefaultTimeout = 60 * time.Second

// DefaultMaxRetries is the default number of retries after the first attempt.
const DefaultMaxRetries = 2

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets the HTTP client used for every request (e.g. to add a
// custom RoundTripper for logging or tests). The SDK never changes its
// Timeout; the per-attempt timeout is applied through the request context.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		if httpClient != nil {
			c.httpClient = httpClient
		}
	}
}

// WithBaseURL overrides the API base URL (default https://api.blindpay.com/v1),
// for tests or proxies. The URL receives the API key.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		if baseURL != "" {
			c.baseURL = baseURL
		}
	}
}

// WithTimeout sets the per-attempt timeout (default 60s). Zero or a negative
// value disables it; the caller's context still applies.
func WithTimeout(d time.Duration) Option {
	return func(c *Client) {
		if d < 0 {
			d = 0
		}
		c.timeout = d
	}
}

// WithMaxRetries sets how many times a failed request may be retried
// (default 2, so up to 3 attempts). Zero disables retries and automatic
// idempotency keys.
//
// GET requests are retried on network errors, timeouts and HTTP 408, 429,
// 500, 502, 503 and 504. POST/PUT/PATCH/DELETE requests are retried only on
// network errors, timeouts, 429 and 409 idempotency_key_in_flight, and only
// with an idempotency key. They are never retried on 5xx: if a create or
// payout returns a 5xx, look the resource up before resending, and resend
// with the IdempotencyKey from the error.
func WithMaxRetries(n int) Option {
	return func(c *Client) {
		if n < 0 {
			n = 0
		}
		c.maxRetries = n
	}
}

// DefaultMaxResponseBodyBytes is the largest response body the client reads
// unless WithMaxResponseBodyBytes changes it (32 MiB).
const DefaultMaxResponseBodyBytes = request.DefaultMaxResponseBodyBytes

// WithMaxResponseBodyBytes caps how many bytes of a response body the client
// reads (default DefaultMaxResponseBodyBytes). A larger response fails with
// an error that wraps ErrResponseTooLarge. Zero or a negative value keeps
// the default.
func WithMaxResponseBodyBytes(n int64) Option {
	return func(c *Client) {
		if n > 0 {
			c.maxResponseBodyBytes = n
		}
	}
}

// RequestOption configures a single request. See package option.
type RequestOption = option.RequestOption

// WithIdempotencyKey is option.WithIdempotencyKey.
func WithIdempotencyKey(key string) RequestOption {
	return option.WithIdempotencyKey(key)
}

// WithHeader is option.WithHeader.
func WithHeader(name, value string) RequestOption {
	return option.WithHeader(name, value)
}

// WithRequestTimeout is option.WithTimeout: it overrides the per-attempt
// timeout for one request. (The root WithTimeout is the client-wide option.)
func WithRequestTimeout(d time.Duration) RequestOption {
	return option.WithTimeout(d)
}
