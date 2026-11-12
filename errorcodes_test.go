package blindpay

import (
	"context"
	"encoding/json"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/stretchr/testify/require"
)

// TestErrorCodes_MatchOracle compares the generated table with
// testdata/error-codes.json, the cross-SDK oracle extracted from the same
// spec (blindpay-v2@7e49addd9): every code, status and retry class, and
// one exported constant per code.
func TestErrorCodes_MatchOracle(t *testing.T) {
	var oracle struct {
		Codes []struct {
			Code   string `json:"code"`
			Status int    `json:"status"`
			Retry  string `json:"retry"`
		} `json:"codes"`
	}
	data, err := os.ReadFile("testdata/error-codes.json")
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(data, &oracle))
	require.Len(t, oracle.Codes, 175)

	classes := map[RetryClass]int{}
	for _, c := range oracle.Codes {
		info, ok := LookupErrorCode(c.Code)
		require.True(t, ok, c.Code)
		require.Equal(t, ErrorCodeInfo{Status: c.Status, Retry: RetryClass(c.Retry)}, info, c.Code)
		classes[info.Retry]++
	}
	require.Equal(t, map[RetryClass]int{RetryClassActionable: 96, RetryClassTerminal: 61, RetryClassRetryable: 18}, classes)

	consts := generatedErrorCodeConstants(t)
	require.Len(t, consts, len(oracle.Codes))
	for _, c := range oracle.Codes {
		require.Contains(t, consts, c.Code)
	}
}

// generatedErrorCodeConstants returns the string values of the ErrorCode*
// constants in errorcodes_gen.go.
func generatedErrorCodeConstants(t *testing.T) map[string]bool {
	t.Helper()
	f, err := parser.ParseFile(token.NewFileSet(), "errorcodes_gen.go", nil, 0)
	require.NoError(t, err)
	out := map[string]bool{}
	ast.Inspect(f, func(n ast.Node) bool {
		vs, ok := n.(*ast.ValueSpec)
		if !ok || len(vs.Values) != 1 {
			return true
		}
		if lit, ok := vs.Values[0].(*ast.BasicLit); ok {
			v, err := strconv.Unquote(lit.Value)
			require.NoError(t, err)
			out[v] = true
		}
		return true
	})
	return out
}

func TestLookupErrorCode(t *testing.T) {
	info, ok := LookupErrorCode(ErrorCodeTransfersSwapAmountBelowMinimum)
	require.True(t, ok)
	require.Equal(t, ErrorCodeInfo{Status: 400, Retry: RetryClassActionable}, info)

	info2, ok := LookupErrorCode("TRANSFERS_SWAP_AMOUNT_BELOW_MINIMUM")
	require.True(t, ok)
	require.Equal(t, info, info2)

	for _, code := range []string{"", "SOMETHING_NEW_NEXT_YEAR", "transfers_swap_amount_below_minimum"} {
		_, ok := LookupErrorCode(code)
		require.False(t, ok, code)
	}
}

func TestAPIError_RetryClass(t *testing.T) {
	for code, want := range map[string]RetryClass{
		"TRANSFERS_SWAP_AMOUNT_BELOW_MINIMUM": RetryClassActionable,
		"AUTH_RATE_LIMITED":                   RetryClassRetryable,
		"AUTH_FORBIDDEN":                      RetryClassTerminal,
		"SOMETHING_NEW_NEXT_YEAR":             RetryClassUnknown,
		"":                                    RetryClassUnknown,
	} {
		require.Equal(t, want, (&APIError{Code: code}).RetryClass(), code)
	}
}

// TestAPIError_BranchOnErrorCode is the README's "Errors" example: the swap
// minimum error, matched through errors.As, a switch on ErrorCode and the
// retry class.
func TestAPIError_BranchOnErrorCode(t *testing.T) {
	c, err := New("test_key", "in_000000000000", WithMaxRetries(0), WithHTTPClient(&http.Client{Transport: &blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodGet,
		Path:   "/v1/instances/in_000000000000/payouts/po_000000000000",
		Status: http.StatusBadRequest,
		Out:    json.RawMessage(`{"success":false,"code":"TRANSFERS_SWAP_AMOUNT_BELOW_MINIMUM","message":"fx_amount_below_minimum","description":"Swaps between different tokens require at least 100 USDC, measured on the USDC side.","errors":[]}`),
	}}))
	require.NoError(t, err)

	_, err = c.Payouts.Get(context.Background(), "po_000000000000")
	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr))

	matched := false
	switch ErrorCode(apiErr.Code) {
	case ErrorCodeTransfersSwapAmountBelowMinimum:
		matched = true
	}
	require.True(t, matched)
	require.Equal(t, "fx_amount_below_minimum", apiErr.Message, "the legacy message stays exposed")
	require.Equal(t, RetryClassActionable, apiErr.RetryClass())
	require.False(t, apiErr.Retryable)
}

// TestAPIError_RetryClassOnFixtures: error bodies without a code (the
// idempotency 4xx/409 and non-JSON bodies) report RetryClassUnknown.
func TestAPIError_RetryClassOnFixtures(t *testing.T) {
	for name, want := range map[string]RetryClass{
		"rate_limited_429.json":          RetryClassRetryable,
		"idempotency_in_flight_409.json": RetryClassUnknown,
		"non_json_502.json":              RetryClassUnknown,
	} {
		f := blindpaytest.LoadErrorFixture(t, name)
		apiErr := request.ParseAPIError(f.Status, f.ResponseBody())
		require.Equal(t, want, apiErr.RetryClass(), name)
	}
}
