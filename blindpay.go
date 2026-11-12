package blindpay

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blindpaylabs/blindpay-go/available"
	"github.com/blindpaylabs/blindpay-go/bankaccounts"
	"github.com/blindpaylabs/blindpay-go/custodialwallets"
	"github.com/blindpaylabs/blindpay-go/customers"
	"github.com/blindpaylabs/blindpay-go/fees"
	"github.com/blindpaylabs/blindpay-go/instances"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/partnerfees"
	"github.com/blindpaylabs/blindpay-go/payables"
	"github.com/blindpaylabs/blindpay-go/payins"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/blindpaylabs/blindpay-go/quotes"
	"github.com/blindpaylabs/blindpay-go/termsofservice"
	"github.com/blindpaylabs/blindpay-go/transfers"
	"github.com/blindpaylabs/blindpay-go/upload"
	"github.com/blindpaylabs/blindpay-go/virtualaccounts"
	"github.com/blindpaylabs/blindpay-go/wallets"
	"github.com/blindpaylabs/blindpay-go/webhookendpoints"
)

// Version is the current version of the SDK.
const Version = "1.22.0"

// Client is the main BlindPay client. It is immutable after New and safe
// for concurrent use by multiple goroutines.
type Client struct {
	// The API key is deliberately not a field: it lives only behind the
	// sub-clients' shared request config, so no fmt verb applied to a
	// Client (or a struct embedding one) can print it.
	baseURL    string
	instanceID string
	httpClient *http.Client
	timeout    time.Duration
	maxRetries int

	maxResponseBodyBytes int64

	Available        *available.Client
	BankAccounts     *bankaccounts.Client
	CustodialWallets *custodialwallets.Client
	Customers        *customers.Client
	Fees             *fees.Client
	Instances        *instances.Client
	PartnerFees      *partnerfees.Client
	Payables         *payables.Client
	Payins           *payins.Client
	Payouts          *payouts.Client
	Quotes           *quotes.Client
	Tos              *termsofservice.Client
	Transfers        *transfers.Client
	Upload           *upload.Client
	VirtualAccounts  *virtualaccounts.Client
	Wallets          *wallets.Client
	OfframpWallets   *wallets.OfframpClient
	WebhookEndpoints *webhookendpoints.Client
}

// New creates a new BlindPay client with the given API key and instance ID.
func New(apiKey, instanceID string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("api key not provided, get your api key on blindpay dashboard")
	}

	if instanceID == "" {
		return nil, fmt.Errorf("instance id not provided, get your instance id on blindpay dashboard")
	}

	c := &Client{
		baseURL:    "https://api.blindpay.com/v1",
		instanceID: instanceID,
		// No http.Client.Timeout: it would cut off the per-attempt timeout,
		// which is applied through the request context instead.
		httpClient: &http.Client{},
		timeout:    DefaultTimeout,
		maxRetries: DefaultMaxRetries,
	}

	for _, opt := range opts {
		opt(c)
	}

	cfg := &config.Config{
		BaseURL:    c.baseURL,
		APIKey:     apiKey,
		InstanceID: c.instanceID,
		HTTPClient: c.httpClient,
		UserAgent:  c.userAgent(),
		Timeout:    c.timeout,
		MaxRetries: c.maxRetries,

		MaxResponseBodyBytes: c.maxResponseBodyBytes,
	}

	c.Available = available.NewClient(cfg)
	c.BankAccounts = bankaccounts.NewClient(cfg)
	c.CustodialWallets = custodialwallets.NewClient(cfg)
	c.Customers = customers.NewClient(cfg)
	c.Fees = fees.NewClient(cfg)
	c.Instances = instances.NewClient(cfg)
	c.PartnerFees = partnerfees.NewClient(cfg)
	c.Payables = payables.NewClient(cfg)
	c.Payins = payins.NewClient(cfg)
	c.Payouts = payouts.NewClient(cfg)
	c.Quotes = quotes.NewClient(cfg)
	c.Tos = termsofservice.NewClient(cfg)
	c.Transfers = transfers.NewClient(cfg)
	c.Upload = upload.NewClient(cfg)
	c.VirtualAccounts = virtualaccounts.NewClient(cfg)
	c.Wallets = wallets.NewClient(cfg)
	c.OfframpWallets = wallets.NewOfframpClient(cfg)
	c.WebhookEndpoints = webhookendpoints.NewClient(cfg)

	return c, nil
}

// String implements fmt.Stringer without exposing the API key. It has a
// value receiver so a Client printed by value (fmt.Sprintf("%+v", *c)) is
// redacted too, not only *Client.
func (c Client) String() string {
	return fmt.Sprintf("blindpay.Client{instanceID: %q, baseURL: %q, apiKey: [REDACTED]}", c.instanceID, c.baseURL)
}

// GoString implements fmt.GoStringer (the %#v verb) without exposing the
// API key.
func (c Client) GoString() string {
	return c.String()
}

// userAgent returns the User-Agent string for requests.
func (c *Client) userAgent() string {
	return fmt.Sprintf("blindpay-go/%s", Version)
}
