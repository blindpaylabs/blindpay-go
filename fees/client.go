package fees

import (
	"context"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
)

// FeeOptions represents fee configuration for a rail.
type FeeOptions struct {
	PayinFlat        float64 `json:"payin_flat"`
	PayinPercentage  float64 `json:"payin_percentage"`
	PayoutFlat       float64 `json:"payout_flat"`
	PayoutPercentage float64 `json:"payout_percentage"`
}

// GetResponse represents the response when retrieving fees.
type GetResponse struct {
	ID                 string     `json:"id"`
	InstanceID         string     `json:"instance_id"`
	ACH                FeeOptions `json:"ach"`
	DomesticWire       FeeOptions `json:"domestic_wire"`
	RTP                FeeOptions `json:"rtp"`
	InternationalSwift FeeOptions `json:"international_swift"`
	Pix                FeeOptions `json:"pix"`
	PixSafe            FeeOptions `json:"pix_safe"`
	AchColombia        FeeOptions `json:"ach_colombia"`
	Transfers3         FeeOptions `json:"transfers_3"`
	Spei               FeeOptions `json:"spei"`
	Tron               FeeOptions `json:"tron"`
	Ethereum           FeeOptions `json:"ethereum"`
	Polygon            FeeOptions `json:"polygon"`
	Base               FeeOptions `json:"base"`
	Arbitrum           FeeOptions `json:"arbitrum"`
	Stellar            FeeOptions `json:"stellar"`
	Solana             FeeOptions `json:"solana"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
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
func (c *Client) Get(ctx context.Context) (*GetResponse, error) {
	path := fmt.Sprintf("/instances/%s/billing/fees", c.instanceID)
	return request.Do[*GetResponse](c.cfg, ctx, "GET", path, nil)
}
