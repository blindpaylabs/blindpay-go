package wallets

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/option"
)

// SetSelfCustody records whether the customer holds the blockchain wallet's
// private keys (required for Brazilian customers). It can be set only once,
// while the wallet's IsSelfCustody is still null; a second call fails with
// 409 WALLETS_SELF_CUSTODY_ALREADY_SET.
func (c *Client) SetSelfCustody(ctx context.Context, customerID, id string, isSelfCustody bool, opts ...option.RequestOption) (*BlockchainWallet, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/blockchain-wallets/%s", c.instanceID, customerID, id)
	return request.Do[*BlockchainWallet](c.cfg, ctx, "PATCH", path, &selfCustodyBody{IsSelfCustody: isSelfCustody}, opts...)
}

// selfCustodyBody is the body of PATCH .../blockchain-wallets/{id}.
type selfCustodyBody struct {
	IsSelfCustody bool `json:"is_self_custody"`
}

// UpdateOfframpWalletParams updates an offramp wallet. partner_fee_id is
// tri-state: leave PartnerFeeID nil to keep the current fee, set it to
// attach a partner fee, or set ClearPartnerFee to remove it (sends null).
type UpdateOfframpWalletParams struct {
	CustomerID      string  `json:"-"`
	BankAccountID   string  `json:"-"`
	ID              string  `json:"-"`
	PartnerFeeID    *string `json:"partner_fee_id,omitempty"`
	ClearPartnerFee bool    `json:"-"`
}

// MarshalJSON sends partner_fee_id as the value, as null (ClearPartnerFee)
// or not at all.
func (p UpdateOfframpWalletParams) MarshalJSON() ([]byte, error) {
	var body updateOfframpWalletBody
	switch {
	case p.ClearPartnerFee:
		body.PartnerFeeID = &nullableString{}
	case p.PartnerFeeID != nil:
		body.PartnerFeeID = &nullableString{value: p.PartnerFeeID}
	}
	return json.Marshal(body)
}

// updateOfframpWalletBody is the wire body of an offramp wallet update; a
// nil PartnerFeeID omits the member.
type updateOfframpWalletBody struct {
	PartnerFeeID *nullableString `json:"partner_fee_id,omitempty"`
}

// nullableString encodes as its value, or as null when value is nil.
type nullableString struct {
	value *string
}

// MarshalJSON implements json.Marshaler.
func (n nullableString) MarshalJSON() ([]byte, error) {
	if n.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(*n.value)
}

// Update changes an offramp wallet's partner fee.
func (c *OfframpClient) Update(ctx context.Context, params *UpdateOfframpWalletParams, opts ...option.RequestOption) error {
	if params == nil {
		return fmt.Errorf("params cannot be nil")
	}
	if params.CustomerID == "" {
		return fmt.Errorf("customer ID cannot be empty")
	}
	if params.BankAccountID == "" {
		return fmt.Errorf("bank account ID cannot be empty")
	}
	if params.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if params.ClearPartnerFee && params.PartnerFeeID != nil {
		return fmt.Errorf("set either PartnerFeeID or ClearPartnerFee, not both")
	}

	path := request.Pathf("/instances/%s/customers/%s/bank-accounts/%s/offramp-wallets/%s",
		c.instanceID, params.CustomerID, params.BankAccountID, params.ID)
	_, err := request.Do[struct{}](c.cfg, ctx, "PUT", path, params, opts...)
	return err
}
