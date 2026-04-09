package types

type Rail string

const (
	RailWire               Rail = "wire"
	RailACH                Rail = "ach"
	RailPix                Rail = "pix"
	RailPixSafe            Rail = "pix_safe"
	RailSpeiBitso          Rail = "spei_bitso"
	RailTransfersBitso     Rail = "transfers_bitso"
	RailACHCopBitso        Rail = "ach_cop_bitso"
	RailInternationalSwift Rail = "international_swift"
	RailRTP                Rail = "rtp"
)

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
	WebhookEventReceiverNew            WebhookEvent = "receiver.new"
	WebhookEventReceiverUpdate         WebhookEvent = "receiver.update"
	WebhookEventBankAccountNew         WebhookEvent = "bankAccount.new"
	WebhookEventPayoutNew              WebhookEvent = "payout.new"
	WebhookEventPayoutUpdate           WebhookEvent = "payout.update"
	WebhookEventPayoutComplete         WebhookEvent = "payout.complete"
	WebhookEventPayoutPartnerFee       WebhookEvent = "payout.partnerFee"
	WebhookEventBlockchainWalletNew    WebhookEvent = "blockchainWallet.new"
	WebhookEventPayinNew               WebhookEvent = "payin.new"
	WebhookEventPayinUpdate            WebhookEvent = "payin.update"
	WebhookEventPayinComplete          WebhookEvent = "payin.complete"
	WebhookEventPayinPartnerFee        WebhookEvent = "payin.partnerFee"
	WebhookEventTosAccept              WebhookEvent = "tos.accept"
	WebhookEventLimitIncreaseNew       WebhookEvent = "limitIncrease.new"
	WebhookEventLimitIncreaseUpdate    WebhookEvent = "limitIncrease.update"
	WebhookEventVirtualAccountNew      WebhookEvent = "virtualAccount.new"
	WebhookEventVirtualAccountComplete WebhookEvent = "virtualAccount.complete"
	WebhookEventTransferNew            WebhookEvent = "transfer.new"
	WebhookEventTransferUpdate         WebhookEvent = "transfer.update"
	WebhookEventTransferComplete       WebhookEvent = "transfer.complete"
	WebhookEventWalletNew              WebhookEvent = "wallet.new"
	WebhookEventWalletInbound          WebhookEvent = "wallet.inbound"
)

type RecipientRelationship string

const (
	RecipientRelationshipFirstParty            RecipientRelationship = "first_party"
	RecipientRelationshipEmployee              RecipientRelationship = "employee"
	RecipientRelationshipIndependentContractor RecipientRelationship = "independent_contractor"
	RecipientRelationshipVendorOrSupplier      RecipientRelationship = "vendor_or_supplier"
	RecipientRelationshipSubsidiaryOrAffiliate RecipientRelationship = "subsidiary_or_affiliate"
	RecipientRelationshipMerchantOrPartner     RecipientRelationship = "merchant_or_partner"
	RecipientRelationshipCustomer              RecipientRelationship = "customer"
	RecipientRelationshipLandlord              RecipientRelationship = "landlord"
	RecipientRelationshipFamily                RecipientRelationship = "family"
	RecipientRelationshipOther                 RecipientRelationship = "other"
)

type PayinPaymentMethod string

const (
	PayinPaymentMethodACH                PayinPaymentMethod = "ach"
	PayinPaymentMethodWire               PayinPaymentMethod = "wire"
	PayinPaymentMethodPix                PayinPaymentMethod = "pix"
	PayinPaymentMethodSpei               PayinPaymentMethod = "spei"
	PayinPaymentMethodTransfers          PayinPaymentMethod = "transfers"
	PayinPaymentMethodPSE                PayinPaymentMethod = "pse"
	PayinPaymentMethodInternationalSwift PayinPaymentMethod = "international_swift"
)

type TrackingStatus string

const (
	TrackingStatusProcessing    TrackingStatus = "processing"
	TrackingStatusOnHold        TrackingStatus = "on_hold"
	TrackingStatusCompleted     TrackingStatus = "completed"
	TrackingStatusPendingReview TrackingStatus = "pending_review"
)

type BusinessIndustry string

const (
	BusinessIndustry541511 BusinessIndustry = "541511"
	BusinessIndustry541512 BusinessIndustry = "541512"
	BusinessIndustry541519 BusinessIndustry = "541519"
	BusinessIndustry518210 BusinessIndustry = "518210"
	BusinessIndustry511210 BusinessIndustry = "511210"
	BusinessIndustry541611 BusinessIndustry = "541611"
	BusinessIndustry541618 BusinessIndustry = "541618"
	BusinessIndustry541330 BusinessIndustry = "541330"
	BusinessIndustry541990 BusinessIndustry = "541990"
	BusinessIndustry522110 BusinessIndustry = "522110"
	BusinessIndustry523110 BusinessIndustry = "523110"
	BusinessIndustry523920 BusinessIndustry = "523920"
	BusinessIndustry423430 BusinessIndustry = "423430"
	BusinessIndustry423690 BusinessIndustry = "423690"
	BusinessIndustry423110 BusinessIndustry = "423110"
	BusinessIndustry423830 BusinessIndustry = "423830"
	BusinessIndustry423840 BusinessIndustry = "423840"
	BusinessIndustry423510 BusinessIndustry = "423510"
	BusinessIndustry424210 BusinessIndustry = "424210"
	BusinessIndustry424690 BusinessIndustry = "424690"
	BusinessIndustry424990 BusinessIndustry = "424990"
	BusinessIndustry454110 BusinessIndustry = "454110"
	BusinessIndustry334111 BusinessIndustry = "334111"
	BusinessIndustry334118 BusinessIndustry = "334118"
	BusinessIndustry325412 BusinessIndustry = "325412"
	BusinessIndustry339112 BusinessIndustry = "339112"
	BusinessIndustry336111 BusinessIndustry = "336111"
	BusinessIndustry336390 BusinessIndustry = "336390"
	BusinessIndustry551112 BusinessIndustry = "551112"
	BusinessIndustry561499 BusinessIndustry = "561499"
	BusinessIndustry488510 BusinessIndustry = "488510"
	BusinessIndustry484121 BusinessIndustry = "484121"
	BusinessIndustry493110 BusinessIndustry = "493110"
	BusinessIndustry424410 BusinessIndustry = "424410"
	BusinessIndustry424480 BusinessIndustry = "424480"
	BusinessIndustry315990 BusinessIndustry = "315990"
	BusinessIndustry313110 BusinessIndustry = "313110"
	BusinessIndustry213112 BusinessIndustry = "213112"
	BusinessIndustry517110 BusinessIndustry = "517110"
	BusinessIndustry541214 BusinessIndustry = "541214"
)
