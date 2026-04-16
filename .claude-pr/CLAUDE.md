# CLAUDE.md -- blindpay-go SDK

This file is for AI agents (Claude) who need to modify this codebase.
Read it fully before making any changes.

## 1. Project structure

```
blindpay-go/
  blindpay.go              # Root client. Creates all sub-clients. Holds `const Version`.
  types.go                 # Re-exports internal/types as public type aliases + constants.
  errors.go                # Re-exports internal/request.APIError as public type alias.
  options.go               # Functional options: WithHTTPClient, WithBaseURL.
  webhooks.go              # VerifyWebhookSignature (standalone, uses svix).
  go.mod                   # Module: github.com/blindpaylabs/blindpay-go (go 1.21)
  VERSIONING.md            # Release process docs.
  internal/
    config/config.go       # config.Config -- bridges top-level config to request.Config.
    request/request.go     # Generic HTTP helper: request.Do[T](...). APIError type.
    types/
      shared.go            # Shared enums: CurrencyType, Network, StablecoinToken, etc.
      types.go             # More shared enums: Rail, WebhookEvent, RecipientRelationship, etc.
    blindpaytest/
      roundtripper.go      # Test helper: mock HTTP RoundTripper for unit tests.
  available/client.go      # GET /available/* -- rails, bank details, NAICS, SWIFT lookup.
  apikeys/client.go        # API key management.
  bankaccounts/client.go   # Bank account CRUD (PIX, ACH, Wire, SPEI, SWIFT, RTP, etc.).
  custodialwallets/...     # Custodial wallet operations.
  fees/client.go           # GET /instances/{id}/billing/fees.
  instances/...            # Instance management.
  partnerfees/...          # Partner fee configuration.
  payins/
    client.go              # Payin CRUD + export. Has sub-client: Quotes.
    quotes.go              # QuotesClient -- payin quote creation + FX rates.
  payouts/client.go        # Payout CRUD + export + document submission.
  quotes/client.go         # Payout quote creation + FX rates.
  receivers/client.go      # Receiver CRUD (individual/business, standard/enhanced KYC).
  transfers/...            # Transfer operations.
  upload/...               # File upload.
  virtualaccounts/...      # Virtual account operations.
  wallets/client.go        # Blockchain wallet CRUD + OfframpClient (second client in same pkg).
  webhookendpoints/...     # Webhook endpoint management.
```

Each resource lives in its own Go package (directory). The package name is the plural lowercase resource name (e.g. `payins`, `receivers`, `wallets`). Multi-word resources use a single lowercase word (e.g. `bankaccounts`, `custodialwallets`, `webhookendpoints`, `virtualaccounts`, `partnerfees`).

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
- **Resource structs**: Singular PascalCase. `Payin`, `Payout`, `Receiver`, `BlockchainWallet`, `BankAccount`.

### Types and fields

- Required fields: bare type (`string`, `float64`, `types.Country`).
- Optional fields: pointer type (`*string`, `*float64`, `*types.Country`).
- ID fields that are path parameters (not sent in JSON body): tagged `json:"-"`. E.g. `ReceiverID string \`json:"-"\``.
- All JSON tags use `snake_case` matching the API.
- Optional JSON fields include `omitempty`.
- Shared enum types are defined in `internal/types/` and re-exported in root `types.go`.
- Package-specific enum types are defined directly in the package file (e.g. `receivers.ProofOfAddressDocType`).

### Go idioms

- Every method takes `context.Context` as the first parameter.
- All HTTP calls go through `request.Do[T](cfg, ctx, method, path, body)`.
- For GET with query params: build `url.Values{}`, append as `?` + `q.Encode()`.
- For DELETE returning 204: return `error` only (no response struct). Call `request.Do[struct{}]`.
- For PUT returning no body: same pattern, return `error` only.
- Validate required parameters at the top of each method. Return `fmt.Errorf(...)` for validation.
- No custom HTTP client usage -- always delegate to `request.Do`.

### Internal wiring

- `config.Config` holds `BaseURL`, `APIKey`, `InstanceID`, `HTTPClient`, `UserAgent`.
- Each sub-client stores `cfg *request.Config` and `instanceID string` (extracted via `cfg.ToRequestConfig()`).
- `request.Config` is a subset: `BaseURL`, `APIKey`, `HTTPClient`, `UserAgent` (no InstanceID).
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
)

// Invoice represents an invoice.
type Invoice struct {
    ID         string  `json:"id"`
    Amount     float64 `json:"amount"`
    Currency   string  `json:"currency"`
    ReceiverID string  `json:"receiver_id"`
}

// CreateParams represents parameters for creating an invoice.
type CreateParams struct {
    ReceiverID string  `json:"receiver_id"`
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
func (c *Client) List(ctx context.Context) ([]Invoice, error) {
    path := fmt.Sprintf("/instances/%s/invoices", c.instanceID)
    return request.Do[[]Invoice](c.cfg, ctx, "GET", path, nil)
}

// Create creates a new invoice.
func (c *Client) Create(ctx context.Context, params *CreateParams) (*CreateResponse, error) {
    if params == nil {
        return nil, fmt.Errorf("params cannot be nil")
    }

    path := fmt.Sprintf("/instances/%s/invoices", c.instanceID)
    return request.Do[*CreateResponse](c.cfg, ctx, "POST", path, params)
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
                Out:    json.RawMessage(`[{"id":"inv_001","amount":100,"currency":"USD","receiver_id":"re_001"}]`),
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

### Step 4: Bump version

Update `const Version` in `blindpay.go` (MINOR bump for new resource).

## 4. How to: Add a method to an existing resource

1. Define any new params/response structs in the resource's `client.go`.
2. Add the method to the `Client` receiver.
3. Follow the existing method patterns in the same file.
4. Add a test case in the corresponding `_test.go` file.

Example -- adding `GetByExternalID` to `receivers`:

```go
// In receivers/client.go

// GetByExternalID retrieves a receiver by external ID.
func (c *Client) GetByExternalID(ctx context.Context, externalID string) (*Receiver, error) {
    if externalID == "" {
        return nil, fmt.Errorf("external ID cannot be empty")
    }

    path := fmt.Sprintf("/instances/%s/receivers/external/%s", c.instanceID, externalID)
    return request.Do[*Receiver](c.cfg, ctx, "GET", path, nil)
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

### Renaming a field

Change the Go field name. If the API field name changed, also update the JSON tag. Bump version: MINOR (breaking).

### Changing a field type (e.g. string to *string)

Update the type. Bump version: MINOR (breaking).

### Fixing a JSON tag typo

Fix the tag. Bump version: PATCH.

## 6. How to: Remove a resource

1. Delete the resource directory (e.g. `rm -rf invoices/`).
2. In `blindpay.go`: remove the import, the struct field, and the initialization line.
3. If the resource's types were re-exported in root `types.go`, remove those aliases too.
4. Bump version: MINOR (breaking).

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

func (c *QuotesClient) Create(ctx context.Context, params *CreateQuoteParams) (*CreateQuoteResponse, error) {
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
    T:      t,                          // *testing.T
    In:     json.RawMessage(`{...}`),   // Expected request body (optional, nil to skip)
    Out:    json.RawMessage(`{...}`),   // Canned response body
    Method: http.MethodGet,             // Expected HTTP method
    Path:   "/instances/in_xxx/...",    // Expected URL path
}
```

## 9. Versioning

- The SDK version is stored in `blindpay.go` as `const Version = "X.Y.Z"` (no `v` prefix).
- Git tags use `v` prefix: `v1.4.0`.
- The SDK stays on **v1.x forever**. Never create v2+ tags (Go module path would need `/v2`).
- MINOR bump: new features, new fields, new enums, breaking changes to existing types.
- PATCH bump: bug fixes, typo corrections, small non-breaking fixes.
- Every PR that changes SDK behavior must bump the Version constant.
- After merging to main: `git tag vX.Y.Z && git push origin vX.Y.Z`, then create a GitHub release.
- Verify on `https://pkg.go.dev/github.com/blindpaylabs/blindpay-go@vX.Y.Z`.

## 10. OpenAPI to SDK mapping rules

### Path to package

| API Path Pattern | Package | Notes |
|---|---|---|
| `/instances/{id}/payins` | `payins` | |
| `/instances/{id}/payouts` | `payouts` | |
| `/instances/{id}/quotes` | `quotes` | Payout quotes |
| `/instances/{id}/payin-quotes` | `payins` (sub-client `Quotes`) | Payin quotes |
| `/instances/{id}/receivers` | `receivers` | |
| `/instances/{id}/receivers/{rid}/blockchain-wallets` | `wallets` | |
| `/instances/{id}/receivers/{rid}/bank-accounts` | `bankaccounts` | |
| `/instances/{id}/receivers/{rid}/bank-accounts/{bid}/offramp-wallets` | `wallets` (`OfframpClient`) | |
| `/instances/{id}/billing/fees` | `fees` | |
| `/available/*` | `available` | No instanceID in path |
| `/e/payins/{id}` | `payins` (method `GetTrack`) | External/public tracking |
| `/e/payouts/{id}` | `payouts` (method `GetTrack`) | External/public tracking |
| `/instances/{id}/export/payins` | `payins` (method `Export`) | |
| `/instances/{id}/export/payouts` | `payouts` (method `Export`) | |

### HTTP method to Go function

| HTTP Method | Go Method Name | Return Type |
|---|---|---|
| `GET` (single) | `Get` | `(*Resource, error)` |
| `GET` (list) | `List` | `(*ListResponse, error)` or `([]Resource, error)` |
| `POST` (create) | `Create` or `Create<Variant>` | `(*CreateResponse, error)` |
| `PUT` (update) | `Update` | `error` (no response body) |
| `DELETE` | `Delete` | `error` (no response body) |
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
| Path parameter (e.g. `{receiver_id}`) | Struct field with `json:"-"`, passed as method arg or in params struct |

### Pagination

For list endpoints that return paginated data, use:

```go
type ListResponse struct {
    Data       []Resource               `json:"data"`
    Pagination types.PaginationMetadata `json:"pagination"`
}
```

For list endpoints that return a simple array, use `[]Resource` directly.

### Query parameters

Build query params manually with `url.Values{}`. Do not use struct tags for query encoding:

```go
q := url.Values{}
if params.Status != "" {
    q.Set("status", string(params.Status))
}
if params.Limit > 0 {
    q.Set("limit", fmt.Sprintf("%d", params.Limit))
}
if len(q) > 0 {
    path += "?" + q.Encode()
}
```
