// Package option provides per-request options accepted by every BlindPay
// SDK method as trailing variadic arguments:
//
//	payout, err := client.Payouts.CreateEvm(ctx, params,
//		option.WithIdempotencyKey("payout-order-1234"),
//		option.WithTimeout(30*time.Second),
//	)
//
// The same options are re-exported from the root package
// (blindpay.WithIdempotencyKey, blindpay.WithHeader, blindpay.WithRequestTimeout).
// Options never mutate the shared client, so a client stays safe for
// concurrent use.
package option

import (
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/request"
)

// RequestOption configures a single request.
type RequestOption = request.RequestOption

// WithIdempotencyKey sends key as the Idempotency-Key header on
// POST/PUT/PATCH/DELETE JSON requests (it is ignored on GET and on
// multipart uploads, which the API does not protect). The API replays the
// stored response for 24 hours when the same key and body are resent.
// The key must be 1-255 characters with no leading or trailing whitespace;
// otherwise the method returns blindpay.ErrInvalidIdempotencyKey without
// sending anything. An Idempotency-Key passed through WithHeader is treated
// the same way.
//
// When the client has retries enabled (the default) and no key is given,
// the SDK generates one before the first attempt; it is available on
// *blindpay.APIError and *blindpay.NetworkError as IdempotencyKey.
func WithIdempotencyKey(key string) RequestOption {
	return request.WithIdempotencyKey(key)
}

// WithHeader adds an extra request header. Calling it again with the same
// name adds another value.
func WithHeader(name, value string) RequestOption {
	return request.WithHeader(name, value)
}

// WithTimeout overrides the client's per-attempt timeout for this request.
func WithTimeout(d time.Duration) RequestOption {
	return request.WithTimeout(d)
}
