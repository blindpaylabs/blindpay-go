package blindpay

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// Placeholder ids used both to call the SDK and to build expected paths.
const (
	tInstance    = "in_000000000000"
	tCustomer    = "cu_000000000000"
	tID          = "id_000000000000"
	tBankAccount = "ba_000000000000"
	tUser        = "us_000000000000"
	tSwift       = "AAAABBCCXXX"
)

// rtClient returns a root client that sends every request to rt, with
// retries off.
func rtClient(t *testing.T, rt http.RoundTripper) *Client {
	t.Helper()
	c, err := New("test_key", tInstance, WithBaseURL("https://api.blindpay.com/v1"), WithHTTPClient(&http.Client{Transport: rt}), WithMaxRetries(0))
	require.NoError(t, err)
	return c
}
