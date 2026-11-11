package blindpay

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/blindpaylabs/blindpay-go/bankaccounts"
	"github.com/blindpaylabs/blindpay-go/custodialwallets"
	"github.com/blindpaylabs/blindpay-go/customers"
	"github.com/blindpaylabs/blindpay-go/payables"
	"github.com/blindpaylabs/blindpay-go/payins"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/blindpaylabs/blindpay-go/transfers"
	"github.com/blindpaylabs/blindpay-go/virtualaccounts"
	"github.com/blindpaylabs/blindpay-go/wallets"
)

// WebhookPayload is a decoded webhook. Type-switch on the concrete types
// (*PayoutWebhook, *PayinWebhook, ..., *UnknownWebhook). EventType returns
// the webhook_event name, e.g. "payout.complete".
type WebhookPayload interface {
	EventType() WebhookEvent
}

// WebhookHeader carries the webhook_event discriminator shared by every
// typed event.
type WebhookHeader struct {
	Event WebhookEvent `json:"webhook_event"`
}

// EventType implements WebhookPayload.
func (h WebhookHeader) EventType() WebhookEvent { return h.Event }

// BankAccountWebhook is bankAccount.new.
type BankAccountWebhook struct {
	WebhookHeader
	bankaccounts.BankAccount
}

// BlockchainWalletWebhook is blockchainWallet.new and blockchainWallet.update.
type BlockchainWalletWebhook struct {
	WebhookHeader
	wallets.BlockchainWallet
}

// CustomerWebhook is customer.new and customer.update.
type CustomerWebhook struct {
	WebhookHeader
	customers.Customer
}

// CustomerDeleteWebhook is customer.delete.
type CustomerDeleteWebhook struct {
	WebhookHeader
	ID         string              `json:"id"`
	CustomerID string              `json:"customer_id"`
	InstanceID string              `json:"instance_id"`
	Type       AccountClass        `json:"type"`
	KycType    customers.KycType   `json:"kyc_type"`
	KycStatus  customers.KycStatus `json:"kyc_status"`
	DeletedAt  *time.Time          `json:"deleted_at,omitempty"`
}

// LimitIncreaseWebhook is limitIncrease.new and limitIncrease.update.
// Amounts are integer cents.
type LimitIncreaseWebhook struct {
	WebhookHeader
	customers.LimitIncreaseRequest
}

// PayinWebhook is payin.new, payin.update, payin.complete and
// payin.partnerFee.
type PayinWebhook struct {
	WebhookHeader
	payins.Payin
}

// PayableWebhook is payable.new, payable.update and payable.complete.
type PayableWebhook struct {
	WebhookHeader
	payables.Payable
}

// PayoutWebhook is payout.new, payout.update, payout.complete and
// payout.partnerFee. A failed payout can arrive on payout.update or
// payout.complete: branch on Status, not on the event name.
type PayoutWebhook struct {
	WebhookHeader
	payouts.Payout
	OfframpWalletID *string `json:"offramp_wallet_id,omitempty"`
}

// TosAcceptWebhook is tos.accept.
type TosAcceptWebhook struct {
	WebhookHeader
	TosID          string  `json:"tos_id"`
	IdempotencyKey string  `json:"idempotency_key"`
	CustomerID     *string `json:"customer_id,omitempty"`
	Version        *string `json:"version,omitempty"`
}

// TransferWebhook is transfer.new, transfer.update and transfer.complete.
type TransferWebhook struct {
	WebhookHeader
	transfers.Transfer
}

// VirtualAccountWebhook is virtualAccount.new and virtualAccount.complete.
type VirtualAccountWebhook struct {
	WebhookHeader
	virtualaccounts.VirtualAccount
}

// WalletWebhook is wallet.new and wallet.update.
type WalletWebhook struct {
	WebhookHeader
	custodialwallets.CustodialWallet
}

// WalletInboundWebhook is wallet.inbound (a deposit to a managed wallet).
type WalletInboundWebhook struct {
	WebhookHeader
	ID              string             `json:"id"`
	Address         string             `json:"address"`
	Network         Network            `json:"network"`
	Token           WalletInboundToken `json:"token"`
	DepositID       string             `json:"deposit_id"`
	TransactionHash string             `json:"transaction_hash"`
}

// WalletInboundToken is the token received by a wallet.inbound deposit.
// Amount is in whole tokens (10.5 = 10.5 USDC), not cents.
type WalletInboundToken struct {
	ID      string          `json:"id"`
	Address string          `json:"address"`
	Symbol  StablecoinToken `json:"symbol"`
	Amount  float64         `json:"amount"`
}

// UnknownWebhook is any event this SDK version does not model. Raw holds
// the full payload.
type UnknownWebhook struct {
	WebhookHeader
	Raw json.RawMessage
}

// webhookTypes maps every webhook_event in the spec to its payload type.
var webhookTypes = map[WebhookEvent]func() WebhookPayload{
	WebhookEventBankAccountNew:         func() WebhookPayload { return &BankAccountWebhook{} },
	WebhookEventBlockchainWalletNew:    func() WebhookPayload { return &BlockchainWalletWebhook{} },
	WebhookEventBlockchainWalletUpdate: func() WebhookPayload { return &BlockchainWalletWebhook{} },
	WebhookEventCustomerNew:            func() WebhookPayload { return &CustomerWebhook{} },
	WebhookEventCustomerUpdate:         func() WebhookPayload { return &CustomerWebhook{} },
	WebhookEventCustomerDelete:         func() WebhookPayload { return &CustomerDeleteWebhook{} },
	WebhookEventLimitIncreaseNew:       func() WebhookPayload { return &LimitIncreaseWebhook{} },
	WebhookEventLimitIncreaseUpdate:    func() WebhookPayload { return &LimitIncreaseWebhook{} },
	WebhookEventPayinNew:               func() WebhookPayload { return &PayinWebhook{} },
	WebhookEventPayinUpdate:            func() WebhookPayload { return &PayinWebhook{} },
	WebhookEventPayinComplete:          func() WebhookPayload { return &PayinWebhook{} },
	WebhookEventPayinPartnerFee:        func() WebhookPayload { return &PayinWebhook{} },
	WebhookEventPayableNew:             func() WebhookPayload { return &PayableWebhook{} },
	WebhookEventPayableUpdate:          func() WebhookPayload { return &PayableWebhook{} },
	WebhookEventPayableComplete:        func() WebhookPayload { return &PayableWebhook{} },
	WebhookEventPayoutNew:              func() WebhookPayload { return &PayoutWebhook{} },
	WebhookEventPayoutUpdate:           func() WebhookPayload { return &PayoutWebhook{} },
	WebhookEventPayoutComplete:         func() WebhookPayload { return &PayoutWebhook{} },
	WebhookEventPayoutPartnerFee:       func() WebhookPayload { return &PayoutWebhook{} },
	WebhookEventTosAccept:              func() WebhookPayload { return &TosAcceptWebhook{} },
	WebhookEventTransferNew:            func() WebhookPayload { return &TransferWebhook{} },
	WebhookEventTransferUpdate:         func() WebhookPayload { return &TransferWebhook{} },
	WebhookEventTransferComplete:       func() WebhookPayload { return &TransferWebhook{} },
	WebhookEventVirtualAccountNew:      func() WebhookPayload { return &VirtualAccountWebhook{} },
	WebhookEventVirtualAccountComplete: func() WebhookPayload { return &VirtualAccountWebhook{} },
	WebhookEventWalletNew:              func() WebhookPayload { return &WalletWebhook{} },
	WebhookEventWalletUpdate:           func() WebhookPayload { return &WalletWebhook{} },
	WebhookEventWalletInbound:          func() WebhookPayload { return &WalletInboundWebhook{} },
}

// ParseWebhookEvent decodes a webhook body WITHOUT verifying it. Use
// ConstructWebhookEvent for untrusted input.
func ParseWebhookEvent(payload []byte) (WebhookPayload, error) {
	var head WebhookHeader
	if err := json.Unmarshal(payload, &head); err != nil {
		return nil, fmt.Errorf("blindpay: decoding webhook: %w", err)
	}

	newPayload, ok := webhookTypes[head.Event]
	if !ok {
		return &UnknownWebhook{WebhookHeader: head, Raw: append(json.RawMessage(nil), payload...)}, nil
	}
	event := newPayload()
	if err := json.Unmarshal(payload, event); err != nil {
		return nil, fmt.Errorf("blindpay: decoding %s webhook: %w", head.Event, err)
	}
	return event, nil
}
