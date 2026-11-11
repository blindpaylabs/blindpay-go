package request

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/stretchr/testify/require"
)

func newTestConfig(rt http.RoundTripper) *Config {
	return &Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		HTTPClient: &http.Client{Transport: rt},
		UserAgent:  "blindpay-go/test",
	}
}

func doError(t *testing.T, status int, header http.Header, body []byte) *APIError {
	t.Helper()

	cfg := newTestConfig(&blindpaytest.RoundTripper{T: t, Status: status, Header: header, Out: body})
	_, err := Do[map[string]any](cfg, context.Background(), http.MethodGet, "/instances/in_000000000000/payouts", nil)
	require.Error(t, err)

	var apiErr *APIError
	require.True(t, errors.As(err, &apiErr), "expected *APIError, got %T: %v", err, err)
	return apiErr
}

func TestDo_ErrorFixtures(t *testing.T) {
	tests := []struct {
		fixture     string
		code        string
		message     string
		description string
		traceID     string
		errors      []string
	}{
		{
			fixture:     "api_error_envelope.json",
			code:        "CUSTOMERS_NOT_FOUND",
			message:     "customer_not_found",
			description: "The customer was not found.",
			traceID:     "0af7651916cd43dd8448eb211c80319c",
		},
		{
			fixture: "auth_401_no_trace_id.json",
			code:    "AUTH_INVALID_API_KEY",
			message: "unauthorized",
		},
		{
			fixture: "idempotency_in_flight_409.json",
			message: "idempotency_key_in_flight",
			errors:  []string{"A request with this Idempotency-Key is still being processed"},
		},
		{
			fixture: "idempotency_key_invalid_400.json",
			message: "idempotency_key_invalid",
			errors:  []string{"Idempotency-Key must be at most 255 characters"},
		},
		{
			fixture: "idempotency_payload_mismatch_422.json",
			message: "idempotency_key_payload_mismatch",
			errors:  []string{"This Idempotency-Key was already used with a different request body"},
		},
		{
			fixture: "internal_error_no_errors_500.json",
			code:    "INTERNAL_ERROR",
			message: "internal_error",
			traceID: "2cf7651916cd43dd8448eb211c80319c",
		},
		{
			fixture: "rate_limited_429.json",
			code:    "AUTH_RATE_LIMITED",
			message: "rate_limited",
			traceID: "1bf7651916cd43dd8448eb211c80319c",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.fixture, func(t *testing.T) {
			f := blindpaytest.LoadErrorFixture(t, tc.fixture)
			header := http.Header{}
			for k, v := range f.Headers {
				header.Set(k, v)
			}

			apiErr := doError(t, f.Status, header, f.ResponseBody())

			require.Equal(t, f.Status, apiErr.StatusCode)
			require.Equal(t, tc.code, apiErr.Code)
			require.Equal(t, tc.message, apiErr.Message)
			require.Equal(t, tc.description, apiErr.Description)
			require.Equal(t, tc.traceID, apiErr.TraceID)
			require.Equal(t, f.ResponseBody(), apiErr.RawBody)

			var got []string
			for _, item := range apiErr.Errors {
				got = append(got, item.Message)
			}
			require.Equal(t, tc.errors, got)
		})
	}
}

func TestDo_ErrorSpecExample(t *testing.T) {
	f := blindpaytest.LoadErrorFixture(t, "api_error_envelope.json")
	require.NotEmpty(t, f.SpecGeneratedExample)

	apiErr := doError(t, http.StatusBadRequest, nil, f.SpecGeneratedExample)

	require.Equal(t, "PAYOUTS_INSUFFICIENT_BALANCE", apiErr.Code)
	require.Equal(t, "not_enough_balance", apiErr.Message)
	require.Equal(t, "There is not enough balance to complete this payout. Fund your balance and try again.", apiErr.Description)
	require.Equal(t, "0af7651916cd43dd8448eb211c80319c", apiErr.TraceID)
	require.Len(t, apiErr.Errors, 1)
	require.Equal(t, "sample_errors", apiErr.Errors[0].Message)

	msg := apiErr.Error()
	require.Contains(t, msg, "status 400")
	require.Contains(t, msg, "PAYOUTS_INSUFFICIENT_BALANCE")
	require.Contains(t, msg, "not_enough_balance")
	require.Contains(t, msg, "[trace_id: 0af7651916cd43dd8448eb211c80319c]")
}

func TestDo_ErrorNonJSONBody(t *testing.T) {
	f := blindpaytest.LoadErrorFixture(t, "non_json_502.json")
	header := http.Header{}
	for k, v := range f.Headers {
		header.Set(k, v)
	}

	apiErr := doError(t, f.Status, header, f.ResponseBody())

	require.Equal(t, http.StatusBadGateway, apiErr.StatusCode)
	require.Empty(t, apiErr.Code)
	require.True(t, strings.HasPrefix(apiErr.Message, "HTTP 502 error"), apiErr.Message)
	require.Equal(t, f.ResponseBody(), apiErr.RawBody)
	require.NotContains(t, apiErr.Error(), "failed to decode")
}

// TestDo_ErrorItemsTolerant covers the server's `errors: unknown[]` (the
// error handler forwards any sanitized value): strings, objects, objects with
// oddly typed members, numbers and nulls. None of them may make the decode
// fail, and the API's message/code must always survive.
func TestDo_ErrorItemsTolerant(t *testing.T) {
	tests := []struct {
		name  string
		items string
		want  []ErrorItem
	}{
		{
			name:  "strings",
			items: `["Bad Request","name is required"]`,
			want:  []ErrorItem{{Message: "Bad Request"}, {Message: "name is required"}},
		},
		{
			name:  "objects",
			items: `[{"code":"invalid_type","message":"Expected string","long_message":"name: Expected string"}]`,
			want:  []ErrorItem{{Code: "invalid_type", Message: "Expected string", LongMessage: "name: Expected string"}},
		},
		{
			name:  "object with oddly typed members",
			items: `[{"code":123,"message":["a"],"path":["name"]}]`,
			want:  []ErrorItem{{Message: `{"code":123,"message":["a"],"path":["name"]}`}},
		},
		{
			name:  "object without message",
			items: `[{"path":["name"],"expected":"string"}]`,
			want:  []ErrorItem{{Message: `{"path":["name"],"expected":"string"}`}},
		},
		{
			name:  "mixed scalars",
			items: `["text",42,true,null,{"message":"m"}]`,
			want:  []ErrorItem{{Message: "text"}, {Message: "42"}, {Message: "true"}, {Message: "null"}, {Message: "m"}},
		},
		{
			name:  "empty",
			items: `[]`,
			want:  []ErrorItem{},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			body := `{"success":false,"code":"VALIDATION_ERROR","message":"validation_error","errors":` + tc.items + `,"trace_id":"t"}`

			apiErr := doError(t, http.StatusBadRequest, nil, []byte(body))

			require.Equal(t, "VALIDATION_ERROR", apiErr.Code)
			require.Equal(t, "validation_error", apiErr.Message)
			require.Equal(t, "t", apiErr.TraceID)
			require.Len(t, apiErr.Errors, len(tc.want))
			for i := range tc.want {
				require.Equal(t, tc.want[i].Code, apiErr.Errors[i].Code)
				require.Equal(t, tc.want[i].Message, apiErr.Errors[i].Message)
				require.Equal(t, tc.want[i].LongMessage, apiErr.Errors[i].LongMessage)
			}
		})
	}
}

func TestDo_ErrorBodyNotAnObject(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "empty", body: ``},
		{name: "json array", body: `["x"]`},
		{name: "json string", body: `"x"`},
		{name: "truncated json", body: `{"message":`},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			apiErr := doError(t, http.StatusServiceUnavailable, nil, []byte(tc.body))
			require.Equal(t, http.StatusServiceUnavailable, apiErr.StatusCode)
			require.True(t, strings.HasPrefix(apiErr.Message, "HTTP 503 error"), apiErr.Message)
		})
	}
}

func TestDo_ErrorOddlyTypedMembers(t *testing.T) {
	body := `{"success":false,"code":"VALIDATION_ERROR","message":"validation_error","description":42,"errors":"name is required","trace_id":"t"}`

	apiErr := doError(t, http.StatusBadRequest, nil, []byte(body))

	require.Equal(t, "VALIDATION_ERROR", apiErr.Code)
	require.Equal(t, "validation_error", apiErr.Message)
	require.Empty(t, apiErr.Description)
	require.Equal(t, "t", apiErr.TraceID)
	require.Len(t, apiErr.Errors, 1)
	require.Equal(t, "name is required", apiErr.Errors[0].Message)
}

func TestDo_ErrorObjectWithoutMessage(t *testing.T) {
	apiErr := doError(t, http.StatusNotFound, nil, []byte(`{"success":false,"code":"NOT_FOUND"}`))
	require.Equal(t, "NOT_FOUND", apiErr.Code)
	require.Equal(t, "HTTP 404 error", apiErr.Message)
}

func TestDo_Empty204(t *testing.T) {
	f := blindpaytest.LoadErrorFixture(t, "empty_204.json")
	cfg := newTestConfig(&blindpaytest.RoundTripper{T: t, Status: f.Status, Out: f.ResponseBody()})

	got, err := Do[*struct{ ID string }](cfg, context.Background(), http.MethodDelete, "/instances/in_000000000000/payables/pa_000000000000", nil)
	require.NoError(t, err)
	require.Nil(t, got)
}

func TestErrorItem_MarshalRoundTrip(t *testing.T) {
	in := ErrorItem{Code: "c", Message: "m", LongMessage: "l"}
	b, err := json.Marshal(in)
	require.NoError(t, err)

	var out ErrorItem
	require.NoError(t, json.Unmarshal(b, &out))
	require.Equal(t, in, out)
}
