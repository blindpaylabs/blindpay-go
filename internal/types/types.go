package types

type Rail string

type BankDetail struct {
	Label    string `json:"label"`
	Regex    string `json:"regex"`
	Key      string `json:"key"`
	Required bool   `json:"required"`
}

type RailEntry struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Country string `json:"country"`
}

type WebhookEvent string

const (
	WebhookEventReceiverNew         WebhookEvent = "receiver.new"
	WebhookEventReceiverUpdate      WebhookEvent = "receiver.update"
	WebhookEventBankAccountNew      WebhookEvent = "bankAccount.new"
	WebhookEventPayoutNew           WebhookEvent = "payout.new"
	WebhookEventPayoutUpdate        WebhookEvent = "payout.update"
	WebhookEventPayoutComplete      WebhookEvent = "payout.complete"
	WebhookEventPayoutPartnerFee    WebhookEvent = "payout.partnerFee"
	WebhookEventBlockchainWalletNew WebhookEvent = "blockchainWallet.new"
	WebhookEventPayinNew            WebhookEvent = "payin.new"
	WebhookEventPayinUpdate         WebhookEvent = "payin.update"
	WebhookEventPayinComplete       WebhookEvent = "payin.complete"
	WebhookEventPayinPartnerFee     WebhookEvent = "payin.partnerFee"
)
