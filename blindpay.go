package blindpay

import (
	"fmt"
	"net/http"
	"time"

	"github.com/blindpaylabs/blindpay-go/apikeys"
	"github.com/blindpaylabs/blindpay-go/available"
	"github.com/blindpaylabs/blindpay-go/bankaccounts"
	"github.com/blindpaylabs/blindpay-go/instances"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/partnerfees"
	"github.com/blindpaylabs/blindpay-go/payins"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/blindpaylabs/blindpay-go/quotes"
	"github.com/blindpaylabs/blindpay-go/receivers"
	"github.com/blindpaylabs/blindpay-go/virtualaccounts"
	"github.com/blindpaylabs/blindpay-go/wallets"
	"github.com/blindpaylabs/blindpay-go/webhookendpoints"
)

// Version is the current version of the SDK.
const Version = "1.2.0"

// Client is the main BlindPay client.
type Client struct {
	baseURL    string
	apiKey     string
	instanceID string
	httpClient *http.Client

	Available        *available.Client
	APIKeys          *apikeys.Client
	BankAccounts     *bankaccounts.Client
	Instances        *instances.Client
	PartnerFees      *partnerfees.Client
	Payins           *payins.Client
	Payouts          *payouts.Client
	Quotes           *quotes.Client
	Receivers        *receivers.Client
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
		apiKey:     apiKey,
		instanceID: instanceID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	cfg := &config.Config{
		BaseURL:    c.baseURL,
		APIKey:     c.apiKey,
		InstanceID: c.instanceID,
		HTTPClient: c.httpClient,
		UserAgent:  c.userAgent(),
	}

	c.Available = available.NewClient(cfg)
	c.APIKeys = apikeys.NewClient(cfg)
	c.BankAccounts = bankaccounts.NewClient(cfg)
	c.Instances = instances.NewClient(cfg)
	c.PartnerFees = partnerfees.NewClient(cfg)
	c.Payins = payins.NewClient(cfg)
	c.Payouts = payouts.NewClient(cfg)
	c.Quotes = quotes.NewClient(cfg)
	c.Receivers = receivers.NewClient(cfg)
	c.VirtualAccounts = virtualaccounts.NewClient(cfg)
	c.Wallets = wallets.NewClient(cfg)
	c.OfframpWallets = wallets.NewOfframpClient(cfg)
	c.WebhookEndpoints = webhookendpoints.NewClient(cfg)

	return c, nil
}

// userAgent returns the User-Agent string for requests.
func (c *Client) userAgent() string {
	return fmt.Sprintf("blindpay-go/%s", Version)
}
