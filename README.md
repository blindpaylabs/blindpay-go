# BlindPay Go SDK

The official Go SDK for [BlindPay](https://blindpay.com), the stablecoin API for global payments.

## Installation

```bash
go get github.com/blindpaylabs/blindpay-go
```

Requires Go 1.21 or later, and a BlindPay API key and instance ID.

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blindpaylabs/blindpay-go"
)

func main() {
	client, err := blindpay.New("your-api-key", "your-instance-id")
	if err != nil {
		log.Fatal(err)
	}

	rails, err := client.Available.GetRails(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	for _, rail := range rails {
		fmt.Printf("%s (%s) in %s\n", rail.Label, rail.Value, rail.Country)
	}
}
```

Every method takes a `context.Context` first; cancelling it stops the request
and any retry wait.

## Resources

| Field on `*blindpay.Client` | API |
|---|---|
| `Available` | Rails, bank details, NAICS codes, SWIFT lookup |
| `BankAccounts` | Customer bank accounts (one `Create<Rail>` per rail, incl. `CreateFasterPayments`), Plaid link |
| `Customers` | Customers (KYC/KYB), limits, limit increases, RFIs (`GetRfi`, `SubmitRfi`) |
| `CustodialWallets` | Managed wallets, balances, deposits (`ListDeposits`, `ListDepositsAutoPaging`) |
| `Wallets` / `OfframpWallets` | Customer blockchain wallets (`SetSelfCustody`), offramp wallets, Solana/Stellar helpers |
| `Instances` | The instance, members, onboarding details, instance RFIs |
| `Payables` | Bills to pay (boletos, PIX codes, invoices) |
| `Quotes` / `Payouts` | Payout quotes and payouts |
| `Payins` (`Payins.Quotes`) | Payin quotes and payins |
| `Transfers` (`Transfers.Quotes`) | Wallet-to-wallet transfers |
| `VirtualAccounts`, `PartnerFees`, `Fees`, `WebhookEndpoints`, `Tos`, `Upload` | The matching resources |

Request and response types live in each resource's package (`payouts.CreateEvmParams`,
`customers.Customer`, ...). Shared enums are re-exported from the root package
(`blindpay.NetworkBase`, `blindpay.TransactionStatusCompleted`). Enum types are plain
strings: a value newer than this SDK decodes as-is instead of failing.

Inputs that take "exactly one of" two fields document it on the field. Payouts
take `SenderWalletAddress` (the customer's wallet) or `WalletID` (a managed
wallet); payout quotes take `BankAccountID` or `PayableID`. The API answers 400
otherwise.

## Configuration

```go
client, err := blindpay.New(apiKey, instanceID,
	blindpay.WithTimeout(30*time.Second), // per attempt, default 60s
	blindpay.WithMaxRetries(3),           // default 2 (3 attempts); 0 disables retries
	blindpay.WithHTTPClient(httpClient),  // custom transport, proxy, logging RoundTripper
	blindpay.WithBaseURL("https://proxy.example.com/v1"),
)
```

- `WithBaseURL` changes where requests go (default `https://api.blindpay.com/v1`).
  The URL receives your API key, so only point it at a host you trust.
- `WithHTTPClient` lets you inject an `http.Client` (for example with a
  `RoundTripper` that logs or mocks requests). The SDK never sets its
  `Timeout`; the per-attempt timeout is applied through the request context.
- A `*blindpay.Client` is immutable and safe for concurrent use. Printing it
  (`%v`, `%#v`) never shows the API key.

### Per-request options

Every method accepts trailing request options:

```go
payout, err := client.Payouts.CreateEvm(ctx, params,
	blindpay.WithIdempotencyKey("payout-order-1234"),
	blindpay.WithRequestTimeout(90*time.Second),
	blindpay.WithHeader("X-Request-Id", "req-1234"),
)
```

The same options live in package `option` (`option.WithIdempotencyKey`,
`option.WithTimeout`, `option.WithHeader`).

## Errors

Every failed call returns an `error`. Use `errors.As` to inspect it:

```go
_, err := client.Payouts.Get(ctx, "po_000000000000")

var apiErr *blindpay.APIError
var netErr *blindpay.NetworkError
switch {
case errors.As(err, &apiErr):
	// The API answered with a non-2xx status.
	log.Printf("status=%d code=%s message=%s trace_id=%s request_id=%s",
		apiErr.StatusCode, apiErr.Code, apiErr.Message, apiErr.TraceID, apiErr.RequestID)
	if apiErr.Code == "PAYOUTS_INSUFFICIENT_BALANCE" {
		// ...
	}
case errors.As(err, &netErr):
	// No response: connection error, timeout or cancellation.
	// errors.Is(err, context.DeadlineExceeded) works too.
case err != nil:
	// Validation error raised before sending (e.g. an empty ID).
}
```

`APIError` fields:

| Field | Meaning |
|---|---|
| `StatusCode` | HTTP status |
| `Code` | Stable UPPERCASE error code; branch on this. Empty for idempotency errors and non-JSON bodies |
| `Message`, `Description` | Human-readable text |
| `Errors` | Extra details (`[]ErrorItem`) |
| `TraceID` | API trace id (absent on auth and rate-limit errors) |
| `RequestID` | The `X-Request-Id` response header; quote it to support |
| `IdempotencyKey` | The key the request carried (yours or generated) |
| `Retryable` | The failure is transient (408, 429, in-flight 409, most 5xx) |
| `Header`, `RawBody` | The raw response |

Non-JSON error bodies (for example a proxy's HTML 502) still produce an
`*APIError` with the status.

## Idempotency and retries

`POST`, `PUT`, `PATCH` and `DELETE` calls accept an idempotency key. The API
stores the response for 24 hours and replays it when the same key and body are
sent again, so a retried payout is not executed twice.

```go
quote, err := client.Quotes.Create(ctx, quoteParams,
	blindpay.WithIdempotencyKey("quote-order-1234"))
```

- The key must be 1 to 255 characters. An empty or longer key returns
  `blindpay.ErrInvalidIdempotencyKey` without sending anything.
- With retries on (the default) and no key given, the SDK generates a UUIDv4
  key before the first attempt and reuses it on every retry.
- Multipart uploads are never protected by the API and never get a key.

What is retried (up to `WithMaxRetries`, with exponential backoff and jitter,
honouring `Retry-After` up to 60 seconds):

| Request | Retried on |
|---|---|
| `GET` | network errors, timeouts, HTTP 408, 429, 500, 502, 503, 504 |
| `POST`/`PUT`/`PATCH`/`DELETE` | network errors, timeouts, 429, 409 `idempotency_key_in_flight` |

Mutations are **never** retried on a 5xx, because the API releases the
idempotency key on a 5xx and a retry would run the operation again. You get an
`*APIError` with `Retryable: true` and the `IdempotencyKey` instead.

> If a create or payout returns a 5xx, look the resource up (by quote, external
> id or a list filter) before resending. Resend with the same
> `apiErr.IdempotencyKey`.

## Pagination

Payouts, payins, customers, transfers and payables have three list methods
(wallet deposits, which are always paginated, have `ListDeposits` and
`ListDepositsAutoPaging`):

```go
// Every row (the API returns all rows when no page size is sent).
all, err := client.Payouts.List(ctx, &payouts.ListParams{Status: blindpay.TransactionStatusCompleted})

// One page with its metadata.
page, err := client.Payouts.ListPage(ctx, &payouts.ListParams{Limit: 50})
if page.HasMore() {
	next, err := client.Payouts.ListPage(ctx, &payouts.ListParams{Limit: 50, StartingAfter: page.NextCursor()})
	_, _ = next, err
}

// Every row, fetched one page at a time (100 per page by default).
pager := client.Payouts.ListAutoPaging(ctx, &payouts.ListParams{CustomerID: "cu_000000000000"})
for pager.Next() {
	p := pager.Current()
	fmt.Println(p.ID, p.Status)
}
if err := pager.Err(); err != nil {
	log.Fatal(err)
}
```

`Limit` accepts 10, 50, 100, 200, 500 or 1000. The auto-pager only goes
forward and rejects `EndingBefore`. On Go 1.23+ you can also range over
`pager.All()`. Prefer the auto-pager for large lists: an unpaginated `List`
loads every row.

## Webhooks

Verify every webhook against the **raw** request body before trusting it. No
client or API key is needed, only the endpoint's signing secret
(`whsec_...`, from `client.WebhookEndpoints.GetSecret`).

```go
func handleWebhook(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body) // raw bytes; never re-encode the JSON
	if err != nil {
		http.Error(w, "read error", http.StatusBadRequest)
		return
	}

	event, err := blindpay.ConstructWebhookEvent(body, r.Header, os.Getenv("BLINDPAY_WEBHOOK_SECRET"))
	if err != nil {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	// Deliveries can repeat: deduplicate on the svix-id header.
	deliveryID := r.Header.Get("svix-id")
	_ = deliveryID

	switch e := event.(type) {
	case *blindpay.PayoutWebhook:
		// A failed payout can arrive on payout.update or payout.complete:
		// branch on the status, not on the event name.
		if e.Status == blindpay.TransactionStatusFailed {
			// ...
		}
	case *blindpay.PayinWebhook:
		// ...
	case *blindpay.UnknownWebhook:
		// An event newer than this SDK; e.Raw holds the payload.
	}
	w.WriteHeader(http.StatusOK)
}
```

- `blindpay.VerifyWebhook(body, r.Header, secret)` returns only a `bool`.
- `blindpay.ParseWebhookEvent(body)` decodes without verifying; use it only on
  trusted input.
- The timestamp must be within 5 minutes; change it with
  `blindpay.WithWebhookTolerance(d)`.
- Never log the webhook secret.

## File upload

```go
f, err := os.Open("passport.pdf")
if err != nil {
	log.Fatal(err)
}
defer f.Close()

res, err := client.Upload.Upload(ctx, &upload.UploadParams{
	File:     f,
	FileName: "passport.pdf",
	Bucket:   upload.BucketOnboarding,
})
if err != nil {
	log.Fatal(err)
}
fmt.Println(res.FileURL)
```

Uploads are sent as `multipart/form-data` and are not retried. Set
`ContentType` to the file's MIME type (default `application/octet-stream`).

`client.Upload.Analyze` checks a document before you submit it and returns how
likely it is to be approved; `client.Upload.Presign` returns a temporary signed
URL for an uploaded file.

## Amounts

Payout, payin, quote and limit amounts are **integer cents** (`1000` = 10.00).
Partner percentage fees are basis points (`100` = 1%), partner flat fees are
cents, and wallet balances are whole tokens (`10.5` = 10.5 USDC).

## Testing against a development instance

Development instances settle with the test token USDB and move no real money:

- A payout of `66600` cents fails and `77700` is refunded.
- Payins complete after about 30 seconds.
- Quotes expire after 5 minutes.
- Development instances are limited to 40 requests per minute (HTTP 429; the
  SDK retries it).

## Examples

[`examples/`](examples) has compile-checked programs for the common workflows:
webhook setup and verification, a managed-wallet payout, a fiat payin, and
reconciliation with the auto-pager.

## Upgrading to 1.21

1.21 is source-compatible with 1.20, but some runtime behaviour changes:

- Requests now have a 60-second per-attempt timeout (`WithTimeout`) and are
  retried (`WithMaxRetries`, default 2); see "Idempotency and retries". Use
  `WithMaxRetries(0)` for the old single-attempt behaviour.
- Mutations carry an auto-generated `Idempotency-Key` when retries are on.
- Optional request fields the spec marks optional are no longer sent when
  unset (for example `sender_wallet_address`, quote `bank_account_id`,
  `cover_fees: false`).
- Errors also expose `RequestID`, `IdempotencyKey` and
  `Retryable`; network failures are `*blindpay.NetworkError`.
- Path parameters are percent-encoded.
- The webhook verifier is native; the `svix` dependency is gone.
- `Payins.Export` and `Payouts.Export` are deprecated: their route does not
  exist and they always fail.

## Documentation

- [BlindPay documentation](https://blindpay.com/docs/getting-started/overview)
- [API reference](https://api.blindpay.com/reference)

## Support

- Email: [support@blindpay.com](mailto:support@blindpay.com)
- Issues: [GitHub Issues](https://github.com/blindpaylabs/blindpay-go/issues)

## License

MIT. See [LICENSE](LICENSE).
