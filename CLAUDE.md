# CLAUDE.md -- blindpay-go SDK

This file is for AI agents (Claude) who need to modify this codebase.
Read it fully before making any changes.

## 1. Project structure

```
blindpay-go/
  blindpay.go              # Root client. Creates all sub-clients. Holds `const Version`.
  types.go                 # Re-exports internal/types as public type aliases + constants.
  errors.go                # Re-exports APIError, ErrorItem, NetworkError, ErrInvalidIdempotencyKey.
  options.go               # Client options (WithHTTPClient, WithBaseURL, WithTimeout, WithMaxRetries)
                           # + root re-exports of request options (WithIdempotencyKey, WithHeader, WithRequestTimeout).
  webhooks.go              # VerifyWebhook / ConstructWebhookEvent (native Svix verifier, no client needed).
  webhook_events.go        # Typed webhook payloads (WebhookPayload union, UnknownWebhook) + ParseWebhookEvent.
  go.mod                   # Module: github.com/blindpaylabs/blindpay-go (go 1.21 -- never raise it)
  option/option.go         # Public per-request options: WithIdempotencyKey, WithHeader, WithTimeout.
  pagination/              # Page[T], Pager[T] (Next/Current/Err); all_go123.go adds All() behind //go:build go1.23.
  internal/
    config/config.go       # config.Config -- bridges top-level config to request.Config.
    request/request.go     # APIError, ErrorItem, ParseAPIError, request.Config.
    request/do.go          # request.Do[T] / DoMultipart[T]: retries (A2 table), timeouts, idempotency keys.
    request/options.go     # RequestOptions / RequestOption (what option.* builds).
    request/path.go        # request.Pathf: builds paths, percent-encoding every segment.
    request/query.go       # request.Query: query builder that skips unset values.
    request/list.go        # request.List[T] / DoList[T]: decodes bare-array and {data, pagination} lists.
    request/body.go        # request.BodyWith: params struct + injected members (e.g. a "type" discriminator).
    types/
      shared.go            # Shared enums: CurrencyType, Network, StablecoinToken, etc.
      types.go             # More shared enums: Rail, WebhookEvent, RecipientRelationship, etc.
    blindpaytest/
      roundtripper.go      # Test helper: mock HTTP RoundTripper (asserts method, path, query, headers, body).
      fixtures.go          # Test helper: loaders for testdata/fixtures (ResponseVariants, LoadErrorFixture).
  available/client.go      # GET /available/* -- rails, bank details, NAICS, SWIFT lookup.
  bankaccounts/client.go   # Bank account CRUD (PIX, ACH, Wire, SPEI, SWIFT, RTP, etc.).
  custodialwallets/...     # Custodial wallet operations.
  customers/client.go      # Customer CRUD (individual/business, standard/enhanced KYC).
  fees/client.go           # GET /instances/{id}/billing/fees.
  instances/...            # Instance management.
  partnerfees/...          # Partner fee configuration.
  payables/                # Payables client (Create/List/ListPage/ListAutoPaging/Get/Delete); types.go holds Payable (also the payable.* webhook payload).
  rfi/                     # RFI types shared by Customers.GetRfi/SubmitRfi and Instances.GetRfi/SubmitRfi.
  payins/
    client.go              # Payin CRUD (+ deprecated Export). Has sub-client: Quotes.
    quotes.go              # QuotesClient -- payin quote creation + FX rates.
  payouts/client.go        # Payout CRUD + document submission (+ deprecated Export).
  quotes/client.go         # Payout quote creation + FX rates.
  transfers/...            # Transfer operations.
  upload/...               # File upload.
  virtualaccounts/...      # Virtual account operations.
  wallets/client.go        # Blockchain wallet CRUD + OfframpClient (second client in same pkg).
  webhookendpoints/...     # Webhook endpoint management.
  cmd/contractcheck/       # `go run ./cmd/contractcheck` -- wire-contract check, see section 8.
  cmd/apisync/             # `go run ./cmd/apisync` -- spec sync tool (-check / -apply / -coverage).
  testdata/fixtures/       # Shared spec-shaped fixtures (verbatim copy of the cross-SDK set), see section 8.
  testdata/webhook-vectors.json # Shared Svix conformance vectors (14), run by webhooks_vectors_test.go.
  requestoptions_guard_test.go  # go/ast guard: every ctx method ends with opts ...option.RequestOption and forwards it.
  opcoverage_test.go       # One entry per x-sdk operation: method + path + decode of every fixture variant (94/94).
  requestbody_test.go      # Wire-body assertions for request fields.
  examples/                # Compile-checked workflow programs (webhooks, payouts, payins, reconciliation).
```

Each resource lives in its own Go package (directory). The package name is the plural lowercase resource name (e.g. `payins`, `customers`, `wallets`). Multi-word resources use a single lowercase word (e.g. `bankaccounts`, `custodialwallets`, `webhookendpoints`, `virtualaccounts`, `partnerfees`). There is no `receivers` or `apikeys` package: `receivers` was fully replaced by `customers`, and API Keys is intentionally absent from the public API reference.

## 2. Conventions

### Naming

- **Packages**: Plural, lowercase, no underscores. `payins`, `payouts`, `bankaccounts`.
- **Client struct**: Always named `Client` (or a variant like `QuotesClient`, `OfframpClient` for sub-resources in the same package).
- **Constructor**: `NewClient(cfg *config.Config) *Client`. Always takes `*config.Config`.
- **Methods**: PascalCase verbs. `List`, `Get`, `Create`, `Update`, `Delete`, `Export`.
  - Variant create methods: `CreateEvm`, `CreateStellar`, `CreateSolana`, `CreateWithAddress`, `CreateWithHash`, `CreatePix`, `CreateAch`, etc.
  - Variant get methods: `GetTrack`, `GetLimits`, `GetFxRate`, `GetWalletMessage`.
- **Params structs**: `<Action>Params` or `<Action><Variant>Params`. E.g. `ListParams`, `CreateEvmParams`, `CreateIndividualStandardParams`.
- **Response structs**: `<Action>Response` or the resource struct itself. E.g. `ListResponse`, `CreateResponse`, `CreateEvmResponse`.
- **Resource structs**: Singular PascalCase. `Payin`, `Payout`, `Customer`, `BlockchainWallet`, `BankAccount`.

### Types and fields

- Required fields (in the spec's `required` list): bare type (`string`, `float64`, `types.Country`).
- Optional or nullable fields: pointer type + `omitempty` (`*string`, `*float64`, `*types.Country`);
  slices and maps need no pointer. This applies to every new field, request or response
  (apisync generates them this way). Older structs still have bare optional fields; leave
  them (changing a field type breaks callers).
- ID fields that are path parameters (not sent in JSON body): tagged `json:"-"`. E.g. `CustomerID string \`json:"-"\``.
- All JSON tags use `snake_case` matching the API.
- Optional JSON fields include `omitempty`.
- Shared enum types are defined in `internal/types/` and re-exported in root `types.go`.
- Package-specific enum types are defined directly in the package file (e.g. `customers.ProofOfAddressDocType`).

### Go idioms

- Every method takes `context.Context` as the first parameter and ends with
  `opts ...option.RequestOption`, which it forwards (`opts...`) to every
  `request.Do*` call. `requestoptions_guard_test.go` enforces this for every
  exported method whose first parameter is `ctx`.
- All HTTP calls go through `request.Do[T](cfg, ctx, method, path, body, opts...)`
  (JSON) or `request.DoMultipart[T]` (multipart). Never call `http.Client`
  directly: retries, timeouts, idempotency keys and error parsing live there.
- Build paths with `request.Pathf("/instances/%s/payouts/%s", c.instanceID, id)`,
  which percent-encodes every segment. Never `fmt.Sprintf` user input into a path.
- Build query strings with `request.Query` (`q.String`, `q.Int`, `q.Apply(path)`),
  which skips unset values. Never send `""`, `0` or `null` for an unset param.
- For DELETE returning 204: return `error` only (no response struct). Call `request.Do[struct{}]`.
- For PUT returning no body: same pattern, return `error` only.
- Send the params struct itself as the body (`request.Do(..., params, opts...)`) so
  every json-tagged field reaches the wire; never copy fields into a
  `map[string]any` or an anonymous struct (a field added later, by hand or by
  apisync, would silently never be sent). When the method injects members
  (a rail `type`, `kyc_type`), use `request.BodyWith(params, map[string]any{...})`.
  Path parameters on params structs are tagged `json:"-"`.
- Spec-optional request fields get `omitempty` (or a pointer + `omitempty` when
  the zero value is meaningful). "Exactly one of" inputs (payout
  `SenderWalletAddress`/`WalletID`, quote `BankAccountID`/`PayableID`) are both
  optional and documented on the field; never validate them client-side.
- New response fields are pointers (or slices/maps) with `omitempty`. A field
  the spec leaves untyped uses `types.LooseString` (or `json.RawMessage`), never
  a type the shared fixtures can break.
- Provider-named spec fields (`checkbook_*`, `onemoney_*`, `circle_wallet_id`,
  `jpm_*`, `aiprise_*`) are not modeled; they are excused in
  `.api-sync/unmodeled.json`.
- Validate required parameters at the top of each method. Return `fmt.Errorf(...)` for validation.
- Code must build and vet on Go 1.21 (`GOTOOLCHAIN=go1.21.13 go build ./... && go vet ./...`):
  no range-over-func and no `iter` package outside files guarded by
  `//go:build go1.23` (see `pagination/all_go123.go`).
- Compatibility policy (the module stays v1 forever, see section 9): changes are
  additive. Never rename or remove an exported identifier and never change a
  field's type or a method's parameters; add a new field or method and mark
  the old one `// Deprecated:`. Two exceptions:
  - A trailing variadic `opts ...option.RequestOption` is the sanctioned way to
    extend a method. Call sites keep compiling; method values, interfaces and
    mocks written against the old signature may not.
  - A type change that cannot be avoided (the old type can never work, or
    adding a slice/map field makes a struct non-comparable) ships in a MINOR
    and is listed in the README "Upgrading" section.

### Internal wiring

- `config.Config` holds `BaseURL`, `APIKey`, `InstanceID`, `HTTPClient`, `UserAgent`, `Timeout`, `MaxRetries`.
- Each sub-client stores `cfg *request.Config` and `instanceID string` (extracted via `cfg.ToRequestConfig()`).
- `request.Config` is the same minus `InstanceID`, plus test seams (`Sleep`, `NewKey`, `Jitter`). It is shared and must never be mutated after construction; per-request settings travel in `RequestOptions`.
- Exception: `available.Client` does not need `instanceID` (its endpoints are not instance-scoped).

## 3. How to: Add a new resource

Example: adding a resource called `invoices` for `GET/POST /instances/{id}/invoices`.

### Step 1: Create the package

Create directory `invoices/` with file `client.go`:

```go
package invoices

import (
    "context"
    "fmt"

    "github.com/blindpaylabs/blindpay-go/internal/config"
    "github.com/blindpaylabs/blindpay-go/internal/request"
    "github.com/blindpaylabs/blindpay-go/option"
)

// Invoice represents an invoice.
type Invoice struct {
    ID         string  `json:"id"`
    Amount     float64 `json:"amount"`
    Currency   string  `json:"currency"`
    CustomerID string  `json:"customer_id"`
}

// CreateParams represents parameters for creating an invoice.
type CreateParams struct {
    CustomerID string  `json:"customer_id"`
    Amount     float64 `json:"amount"`
    Currency   string  `json:"currency"`
}

// CreateResponse represents the response when creating an invoice.
type CreateResponse struct {
    ID string `json:"id"`
}

// Client handles invoice-related operations.
type Client struct {
    cfg        *request.Config
    instanceID string
}

// NewClient creates a new invoices client.
func NewClient(cfg *config.Config) *Client {
    return &Client{
        cfg:        cfg.ToRequestConfig(),
        instanceID: cfg.InstanceID,
    }
}

// List retrieves all invoices.
func (c *Client) List(ctx context.Context, opts ...option.RequestOption) ([]Invoice, error) {
    path := request.Pathf("/instances/%s/invoices", c.instanceID)
    return request.Do[[]Invoice](c.cfg, ctx, "GET", path, nil, opts...)
}

// Create creates a new invoice.
func (c *Client) Create(ctx context.Context, params *CreateParams, opts ...option.RequestOption) (*CreateResponse, error) {
    if params == nil {
        return nil, fmt.Errorf("params cannot be nil")
    }

    path := request.Pathf("/instances/%s/invoices", c.instanceID)
    return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params, opts...)
}
```

### Step 2: Wire into blindpay.go

1. Add the import:
```go
"github.com/blindpaylabs/blindpay-go/invoices"
```

2. Add the field to the `Client` struct (alphabetical order):
```go
Invoices *invoices.Client
```

3. Initialize in `New()` after `cfg` is built:
```go
c.Invoices = invoices.NewClient(cfg)
```

### Step 3: Add tests

Create `invoices/invoices_test.go`. Use `blindpaytest.RoundTripper` to mock HTTP:

```go
package invoices

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "testing"

    "github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
    "github.com/blindpaylabs/blindpay-go/internal/config"
    "github.com/stretchr/testify/require"
)

func TestInvoices_List(t *testing.T) {
    instanceID := "in_000000000000"

    cfg := &config.Config{
        BaseURL:    "https://api.blindpay.com",
        APIKey:     "test-key",
        InstanceID: instanceID,
        HTTPClient: &http.Client{
            Transport: &blindpaytest.RoundTripper{
                T:      t,
                Out:    json.RawMessage(`[{"id":"inv_001","amount":100,"currency":"USD","customer_id":"re_001"}]`),
                Method: http.MethodGet,
                Path:   fmt.Sprintf("/instances/%s/invoices", instanceID),
            },
        },
        UserAgent: "test",
    }

    client := NewClient(cfg)
    invoices, err := client.List(context.Background())
    require.NoError(t, err)
    require.Len(t, invoices, 1)
    require.Equal(t, "inv_001", invoices[0].ID)
}
```

### Step 4: Register the operation

Add one entry per new operation to `opTable` in `opcoverage_test.go` (it fails
for any operation in `testdata/fixtures/index.json` without one) and, for list
operations with cursors, `ListPage`/`ListAutoPaging` (section 10, "Cursor-paginated
lists"). Map new spec schemas in `.api-sync/spec-map.json` (section 11).

### Step 5: Bump version

Update `const Version` in `blindpay.go` (MINOR bump for new resource).

## 4. How to: Add a method to an existing resource

1. Define any new params/response structs in the resource's `client.go`.
2. Add the method to the `Client` receiver.
3. Follow the existing method patterns in the same file.
4. Add a test case in the corresponding `_test.go` file.

Example -- adding `GetByExternalID` to `customers`:

```go
// In customers/client.go

// GetByExternalID retrieves a customer by external ID.
func (c *Client) GetByExternalID(ctx context.Context, externalID string, opts ...option.RequestOption) (*Customer, error) {
    if externalID == "" {
        return nil, fmt.Errorf("external ID cannot be empty")
    }

    path := request.Pathf("/instances/%s/customers/external/%s", c.instanceID, externalID)
    return request.Do[*Customer](c.cfg, ctx, "GET", path, nil, opts...)
}
```

No changes to `blindpay.go` needed. Bump version (MINOR).

## 5. How to: Modify types

### Adding a field to a struct

Add the field with the correct JSON tag. Use pointer for optional, bare type for required:

```go
// Required field
NewField string `json:"new_field"`

// Optional field
NewField *string `json:"new_field,omitempty"`
```

Bump version: MINOR.

### Adding a shared enum value

1. Add the constant in `internal/types/shared.go` or `internal/types/types.go`.
2. Re-export in root `types.go` using `const NewValue = types.NewValue`.

Bump version: MINOR.

### Adding a package-specific enum value

Add the constant directly in the resource package file:

```go
const NewEnumValue MyEnumType = "new_value"
```

Bump version: MINOR.

### Renaming a field or changing its type

Don't (section 2, compatibility policy). Add a new field next to the old one
(new JSON tag, or a new Go name for the same tag on a different struct) and
mark the old one `// Deprecated:`. Only an unavoidable type change ships, in a
MINOR, listed in the README "Upgrading" section.

### Fixing a JSON tag typo

Fix the tag. Bump version: PATCH.

## 6. How to: Retire a resource or method

Never delete it (the module stays v1). Mark the package, client field or
method `// Deprecated:` with the replacement, keep it compiling, and bump
MINOR. A method whose route no longer exists stays too, documented as
always failing (see `Payouts.Export`).

## 7. How to: Add a sub-resource (nested client pattern)

See `payins/quotes.go` and `wallets/client.go` (OfframpClient) as examples.

### Pattern A: Sub-client in a separate file (same package)

Used by payins -> Quotes. The sub-client struct has a different name.

1. Create a new file in the same package (e.g. `payins/quotes.go`).
2. Define the sub-client:

```go
// QuotesClient handles payin quote-related operations.
type QuotesClient struct {
    cfg        *request.Config
    instanceID string
}

func (c *QuotesClient) Create(ctx context.Context, params *CreateQuoteParams, opts ...option.RequestOption) (*CreateQuoteResponse, error) {
    // ...
}
```

3. In the parent client struct, add the field:

```go
type Client struct {
    cfg        *request.Config
    instanceID string
    Quotes     *QuotesClient   // <-- exported sub-client
}
```

4. Initialize in `NewClient`:

```go
func NewClient(cfg *config.Config) *Client {
    reqCfg := cfg.ToRequestConfig()
    return &Client{
        cfg:        reqCfg,
        instanceID: cfg.InstanceID,
        Quotes:     &QuotesClient{cfg: reqCfg, instanceID: cfg.InstanceID},
    }
}
```

Usage: `client.Payins.Quotes.Create(ctx, params)`

### Pattern B: Second client in same package

Used by wallets -> OfframpClient. Exposed as a separate field on the root client.

1. Define `OfframpClient` in `wallets/client.go` with `NewOfframpClient`.
2. In `blindpay.go`, add a separate field:

```go
Wallets        *wallets.Client
OfframpWallets *wallets.OfframpClient
```

3. Initialize both in `New()`:

```go
c.Wallets = wallets.NewClient(cfg)
c.OfframpWallets = wallets.NewOfframpClient(cfg)
```

Usage: `client.OfframpWallets.Create(ctx, params)`

## 8. Testing

### Running tests

```bash
go test ./...
```

### Build check

```bash
go build ./...
```

### Contract check

```bash
go run ./cmd/contractcheck
```

Flags any json struct tag the SDK declares that doesn't exist anywhere in the
committed spec snapshot (`.api-sync/spec-snapshot.json`), plus any spec
webhook-event enum member missing from `internal/types.WebhookEvent`. Genuine
pre-existing divergences that can't be fixed in the current PR go in
`.api-sync/contract-allowlist.json` with a reason and owner; do not use it to
silence a mismatch the PR itself introduced.

### Test structure

- Each resource package has its own `_test.go` file (e.g. `payins/payins_test.go`).
- Tests use `blindpaytest.RoundTripper` as a mock HTTP transport.
- The mock verifies: HTTP method, URL path, and optionally the request body (via the `In` field).
- The mock returns a canned JSON response (via the `Out` field).
- Tests use `github.com/stretchr/testify/require` for assertions.
- Tests construct `config.Config` directly (no real HTTP calls).
- Test function naming: `Test<Resource>_<Method>` (e.g. `TestPayins_List`, `TestPayins_CreateQuote`).

### Mock RoundTripper fields

```go
blindpaytest.RoundTripper{
    T:               t,                          // *testing.T
    In:              json.RawMessage(`{...}`),   // Expected request body (optional, nil to skip)
    Out:             json.RawMessage(`{...}`),   // Canned response body
    Method:          http.MethodGet,             // Expected HTTP method
    Path:            "/instances/in_xxx/...",    // Expected URL path
    Query:           url.Values{"limit": {"10"}}, // Exact query (nil skips; url.Values{} asserts "no query")
    RequestHeader:   http.Header{...},           // Request headers that must match exactly (optional)
    NoRequestHeader: []string{"Idempotency-Key"}, // Request headers that must be absent (optional)
    Status:          http.StatusOK,              // Response status (default 200)
    Header:          http.Header{...},           // Response headers (optional)
}
```

### Shared spec-shaped fixtures

`testdata/fixtures/` is a verbatim copy of the cross-SDK reference fixture set
(see its `README.md`; regenerate by copying the directory, never hand-edit).
Load them with `blindpaytest.ResponseVariants(t, "GET /v1/instances/{instance_id}/payouts")`
(every variant of an operation's success body; list ops have `.b0` envelope and
`.b1` bare-array variants) and `blindpaytest.LoadErrorFixture(t, "rate_limited_429.json")`.
Response-decoding tests must use these instead of hand-written JSON.

## 9. Versioning

- The SDK version is stored in `blindpay.go` as `const Version = "X.Y.Z"` (no `v` prefix).
- Git tags use `v` prefix: `v1.4.0`.
- The SDK stays on **v1.x forever**. Never create v2+ tags (Go module path would need `/v2`).
- MINOR bump: new features, methods, fields and enum values, and the rare
  unavoidable type change (section 2), which must be listed in the README
  "Upgrading" section.
- PATCH bump: bug fixes, typo corrections, small non-breaking fixes.
- Every PR that changes SDK behavior must bump the Version constant.
- Releases are automated: merging a PR labelled `release` makes
  `.github/workflows/release.yml` read `const Version`, push tag `vX.Y.Z` and
  create the GitHub release. Never tag by hand.
- Verify on `https://pkg.go.dev/github.com/blindpaylabs/blindpay-go@vX.Y.Z`.

## 10. OpenAPI to SDK mapping rules

### Path to package

| API Path Pattern | Package | Notes |
|---|---|---|
| `/instances/{id}/payins` | `payins` | |
| `/instances/{id}/payouts` | `payouts` | |
| `/instances/{id}/quotes` | `quotes` | Payout quotes |
| `/instances/{id}/payin-quotes` | `payins` (sub-client `Quotes`) | Payin quotes |
| `/instances/{id}/customers` | `customers` | |
| `/instances/{id}/customers/{cid}/blockchain-wallets` | `wallets` | |
| `/instances/{id}/customers/{cid}/bank-accounts` | `bankaccounts` | |
| `/instances/{id}/customers/{cid}/bank-accounts/{bid}/offramp-wallets` | `wallets` (`OfframpClient`) | |
| `/instances/{id}/billing/fees` | `fees` | |
| `/available/*` | `available` | No instanceID in path |
| `/e/payins/{id}` | `payins` (method `GetTrack`) | External/public tracking |
| `/e/payouts/{id}` | `payouts` (method `GetTrack`) | External/public tracking |
| `/instances/{id}` (GET), `/instances/{id}/onboarding/*`, `/instances/{id}/rfi` | `instances` (`Get`, `GetBusinessDetails`/`GetBusinessProfile`/`GetOwnershipDocuments`, `GetRfi`/`SubmitRfi`) | RFI types in `rfi` |
| `/instances/{id}/customers/{cid}/rfi` | `customers` (`GetRfi`/`SubmitRfi`) | RFI types in `rfi` |
| `/instances/{id}/payables` | `payables` | |
| `/instances/{id}/customers/{cid}/wallets/{wid}/deposits` | `custodialwallets` (`ListDeposits`, `ListDepositsAutoPaging`) | Always an envelope |
| `/upload`, `/upload/analyze`, `/presign` | `upload` (`Upload`, `Analyze`, `Presign`) | Send `?instance_id=` = the client's instance |
| `/instances/{id}/export/payins` | `payins` (method `Export`) | Deprecated: route does not exist (always 404); kept for compatibility (v1 forever) |
| `/instances/{id}/export/payouts` | `payouts` (method `Export`) | Deprecated: route does not exist (always 404); kept for compatibility (v1 forever) |

### HTTP method to Go function

| HTTP Method | Go Method Name | Return Type |
|---|---|---|
| `GET` (single) | `Get` | `(*Resource, error)` |
| `GET` (list) | `List` | `(*ListResponse, error)` or `([]Resource, error)` |
| `POST` (create) | `Create` or `Create<Variant>` | `(*CreateResponse, error)` |
| `PUT` / `PATCH` (update) | `Update`, or an action name (`SetSelfCustody`) | `error` when the API returns only `{success}`, else the resource |
| `DELETE` | `Delete` | `error` when the API returns nothing or `{success}`, else the resource (`Payables.Delete`) |
| `POST` (action) | Descriptive name: `Export`, `SubmitDocuments`, `AuthorizeStellarToken` | Varies |

### Schema to struct

| OpenAPI Schema | Go Type |
|---|---|
| `string` (required) | `string` |
| `string` (optional) | `*string` |
| `number` (required) | `float64` |
| `number` (optional) | `*float64` |
| `integer` (required) | `int` |
| `boolean` (required) | `bool` |
| `boolean` (optional) | `*bool` |
| `string` enum | Named `type X string` + `const` block |
| `object` (nested) | Named struct or inline `struct{}` |
| `array` | `[]T` |
| `date-time` | `time.Time` |
| Property name `foo_bar` | JSON tag `json:"foo_bar"` or `json:"foo_bar,omitempty"` |
| Path parameter (e.g. `{customer_id}`) | Struct field with `json:"-"`, passed as method arg or in params struct |

### Pagination

For list endpoints that return paginated data, use:

```go
type ListResponse struct {
    Data       []Resource               `json:"data"`
    Pagination types.PaginationMetadata `json:"pagination"`
}
```

For list endpoints that return a simple array, use `[]Resource` directly.

Payouts, payins, payables and customers return a bare array when none of
`limit` / `starting_after` / `ending_before` is sent, and the
`{data, pagination}` envelope when any is (transfers and wallet deposits
always return the envelope). Decode these with
`request.DoList[Resource](cfg, ctx, path)` (rows only) or
`request.Do[request.List[Resource]](...)` (rows plus optional `Pagination`),
which detect the shape from the body; never infer the shape from the params
you sent.

### Errors

Every non-2xx response becomes `*request.APIError` (re-exported as
`blindpay.APIError`) via `request.ParseAPIError`. Code paths that cannot use
`request.Do` (multipart upload) must call `request.ParseAPIError` too.
`ErrorItem` decoding accepts any JSON value; never change `Errors`' type.

### Query parameters

Build query params with `request.Query`, which skips unset values. Do not use
struct tags for query encoding:

```go
var q request.Query
q.Int("limit", params.Limit)                 // skipped when 0
q.String("starting_after", params.StartingAfter)
q.String("status", string(params.Status))    // skipped when ""
path = q.Apply(path)
```

### Cursor-paginated lists

Ops whose spec accepts `starting_after` (payouts, payins, customers, transfers,
payables, wallet deposits) expose `ListPage(ctx, *ListParams, opts...)
(*pagination.Page[T], error)` and `ListAutoPaging(ctx, *ListParams, opts...)
*pagination.Pager[T]`. Copy `payouts/client.go`: `ListAutoPaging` defaults
`Limit` to `pagination.DefaultLimit` (100), rejects `EndingBefore` with
`pagination.ErrEndingBeforeUnsupported`, and fetches each page through
`ListPage` with `opts...`. Never change an existing `List` return type.

### Webhooks

`webhooks.go` verifies Svix signatures natively (no dependency): keep it
stdlib-only. `webhook_events.go` maps every spec `webhook_event` to a typed
payload (`webhookTypes`); a payload embeds `WebhookHeader` plus the resource
struct. When the spec adds an event, add it to `webhookTypes`, the
`WebhookEvent` enum, and the fixture copy; `webhook_events_test.go` fails
otherwise. Never give an embedded resource struct an `UnmarshalJSON` method:
it would be promoted to the webhook wrapper and drop `webhook_event`.

## 11. api-sync (`cmd/apisync`)

```bash
go run ./cmd/apisync -validate-map                 # spec-map.json anchors resolve
go run ./cmd/apisync -check                        # vs the committed snapshot (CI)
go run ./cmd/apisync -check -spec <openapi.json>   # vs another spec, read-only
go run ./cmd/apisync -apply -spec <openapi.json>   # WRITES code, Version, snapshot
go run ./cmd/apisync -coverage [-spec <openapi.json>]   # METHOD+path coverage report
```

- Only `x-sdk: true` operations count; `x-sdk: false` ones (yield, card
  payment, simulators, tos passthrough, upload/extract) are ignored, and their
  schemas need no mapping. `-coverage` reports them if the SDK ever sends one.
- Coverage is METHOD + path: it pairs each `request.Do*` call's HTTP method with
  the path literals of the function and of same-package helpers it calls.
- A new spec operation the SDK already sends (hand-written) is not
  regenerated. Generated identifiers use Go initialisms (`PartnerFeeID`,
  `LogoURL`); existing identifiers are never renamed.
- `.api-sync/spec-snapshot.json` must be byte-identical to the spec the SDK was
  synced against; `-apply` copies it. `-check` must be silent before merging.
- `-apply` bumps `Version`; inside a PR that already bumped it, reset it to the
  PR's version.
- Classifying a new schema: map it in `spec-map.json` `types[]` (union-policy
  types need a `property_sites` entry per new property), tie enum properties
  to an `enums[]` symbol, or excuse it in `unmodeled.json` with a real reason.
  Both files are hand-formatted: edit them in place, keeping one-line entries.
- `cmd/apisync/golden_test.go` deletes and regenerates `partnerfees.Get`/`Create`
  in a copy of the repo; keep its quoted source in sync when you edit them.
