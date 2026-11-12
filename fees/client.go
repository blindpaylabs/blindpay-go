// Package fees reads the instance's BlindPay fee schedule (flat fees in
// integer cents, percentage fees in basis points). Use it through
// blindpay.Client.Fees.
package fees

import (
	"context"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/option"
)

// FeeOptions represents fee configuration for a rail.
type FeeOptions struct {
	// PayinFlat is in integer cents (1000 = 10.00).
	PayinFlat float64 `json:"payin_flat"`
	// PayinPercentage is in basis points (100 = 1%).
	PayinPercentage float64 `json:"payin_percentage"`
	// PayoutFlat is in integer cents (1000 = 10.00).
	PayoutFlat float64 `json:"payout_flat"`
	// PayoutPercentage is in basis points (100 = 1%).
	PayoutPercentage float64 `json:"payout_percentage"`
}

// YieldFee is how managed-wallet yield is split, in basis points
// (100 = 1%). PartnerWallets maps a network to the partner's payout
// address on it.
type YieldFee struct {
	BlindpayShareBps int64             `json:"blindpay_share_bps"`
	PartnerShareBps  int64             `json:"partner_share_bps"`
	PartnerWallets   map[string]string `json:"partner_wallets,omitempty"`
}

// GetResponse represents the response when retrieving fees.
type GetResponse struct {
	ID                 string      `json:"id"`
	InstanceID         string      `json:"instance_id"`
	ACH                FeeOptions  `json:"ach"`
	DomesticWire       FeeOptions  `json:"domestic_wire"`
	RTP                FeeOptions  `json:"rtp"`
	InternationalSwift FeeOptions  `json:"international_swift"`
	Pix                FeeOptions  `json:"pix"`
	PixSafe            FeeOptions  `json:"pix_safe"`
	AchColombia        FeeOptions  `json:"ach_colombia"`
	Transfers3         FeeOptions  `json:"transfers_3"`
	Spei               FeeOptions  `json:"spei"`
	Tron               FeeOptions  `json:"tron"`
	Ethereum           FeeOptions  `json:"ethereum"`
	Polygon            FeeOptions  `json:"polygon"`
	Base               FeeOptions  `json:"base"`
	Arbitrum           FeeOptions  `json:"arbitrum"`
	Stellar            FeeOptions  `json:"stellar"`
	Solana             FeeOptions  `json:"solana"`
	CreatedAt          time.Time   `json:"created_at"`
	UpdatedAt          time.Time   `json:"updated_at"`
	Sepa               FeeOptions  `json:"sepa"`
	FasterPayments     FeeOptions  `json:"faster_payments"`
	Ted                *FeeOptions `json:"ted,omitempty"`
	// Yield is the split of managed-wallet yield, if enabled.
	Yield *YieldFee `json:"yield,omitempty"`
}

// Client handles fee-related operations.
type Client struct {
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new fees client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// Get retrieves the fee configuration for the instance.
func (c *Client) Get(ctx context.Context, opts ...option.RequestOption) (*GetResponse, error) {
	path := request.Pathf("/instances/%s/billing/fees", c.instanceID)
	return request.Do[*GetResponse](c.cfg, ctx, "GET", path, nil, opts...)
}
