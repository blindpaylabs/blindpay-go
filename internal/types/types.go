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
	RailTed                Rail = "ted"
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
	PayinPaymentMethodRTP                PayinPaymentMethod = "rtp"
	PayinPaymentMethodTed                PayinPaymentMethod = "ted"
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
	BusinessIndustry541511      BusinessIndustry = "541511"
	BusinessIndustry541512      BusinessIndustry = "541512"
	BusinessIndustry541519      BusinessIndustry = "541519"
	BusinessIndustry518210      BusinessIndustry = "518210"
	BusinessIndustry511210      BusinessIndustry = "511210"
	BusinessIndustry541611      BusinessIndustry = "541611"
	BusinessIndustry541618      BusinessIndustry = "541618"
	BusinessIndustry541330      BusinessIndustry = "541330"
	BusinessIndustry541990      BusinessIndustry = "541990"
	BusinessIndustry522110      BusinessIndustry = "522110"
	BusinessIndustry423430      BusinessIndustry = "423430"
	BusinessIndustry423690      BusinessIndustry = "423690"
	BusinessIndustry423110      BusinessIndustry = "423110"
	BusinessIndustry423830      BusinessIndustry = "423830"
	BusinessIndustry423840      BusinessIndustry = "423840"
	BusinessIndustry423510      BusinessIndustry = "423510"
	BusinessIndustry424210      BusinessIndustry = "424210"
	BusinessIndustry424690      BusinessIndustry = "424690"
	BusinessIndustry424990      BusinessIndustry = "424990"
	BusinessIndustry454110      BusinessIndustry = "454110"
	BusinessIndustry334111      BusinessIndustry = "334111"
	BusinessIndustry334118      BusinessIndustry = "334118"
	BusinessIndustry325412      BusinessIndustry = "325412"
	BusinessIndustry339112      BusinessIndustry = "339112"
	BusinessIndustry336390      BusinessIndustry = "336390"
	BusinessIndustry551112      BusinessIndustry = "551112"
	BusinessIndustry561499      BusinessIndustry = "561499"
	BusinessIndustry488510      BusinessIndustry = "488510"
	BusinessIndustry484121      BusinessIndustry = "484121"
	BusinessIndustry493110      BusinessIndustry = "493110"
	BusinessIndustry424410      BusinessIndustry = "424410"
	BusinessIndustry424480      BusinessIndustry = "424480"
	BusinessIndustry315990      BusinessIndustry = "315990"
	BusinessIndustry313110      BusinessIndustry = "313110"
	BusinessIndustry213112      BusinessIndustry = "213112"
	BusinessIndustry541214      BusinessIndustry = "541214"
	BusinessIndustry111998      BusinessIndustry = "111998"
	BusinessIndustry112120      BusinessIndustry = "112120"
	BusinessIndustry113310      BusinessIndustry = "113310"
	BusinessIndustry115114      BusinessIndustry = "115114"
	BusinessIndustry541211      BusinessIndustry = "541211"
	BusinessIndustry541810      BusinessIndustry = "541810"
	BusinessIndustry541430      BusinessIndustry = "541430"
	BusinessIndustry541715      BusinessIndustry = "541715"
	BusinessIndustry541930      BusinessIndustry = "541930"
	BusinessIndustry561422      BusinessIndustry = "561422"
	BusinessIndustry561311      BusinessIndustry = "561311"
	BusinessIndustry561612      BusinessIndustry = "561612"
	BusinessIndustry561740      BusinessIndustry = "561740"
	BusinessIndustry561730      BusinessIndustry = "561730"
	BusinessIndustry236115      BusinessIndustry = "236115"
	BusinessIndustry236220      BusinessIndustry = "236220"
	BusinessIndustry237310      BusinessIndustry = "237310"
	BusinessIndustry238210      BusinessIndustry = "238210"
	BusinessIndustry811111      BusinessIndustry = "811111"
	BusinessIndustry812111      BusinessIndustry = "812111"
	BusinessIndustry812112      BusinessIndustry = "812112"
	BusinessIndustry532111      BusinessIndustry = "532111"
	BusinessIndustry624410      BusinessIndustry = "624410"
	BusinessIndustry541922      BusinessIndustry = "541922"
	BusinessIndustry811210      BusinessIndustry = "811210"
	BusinessIndustry812199      BusinessIndustry = "812199"
	BusinessIndustry611110      BusinessIndustry = "611110"
	BusinessIndustry611310      BusinessIndustry = "611310"
	BusinessIndustry611410      BusinessIndustry = "611410"
	BusinessIndustry611710      BusinessIndustry = "611710"
	BusinessIndustry211120      BusinessIndustry = "211120"
	BusinessIndustry212114      BusinessIndustry = "212114"
	BusinessIndustry221310      BusinessIndustry = "221310"
	BusinessIndustry562111      BusinessIndustry = "562111"
	BusinessIndustry562920      BusinessIndustry = "562920"
	BusinessIndustry522210      BusinessIndustry = "522210"
	BusinessIndustry522320      BusinessIndustry = "522320"
	BusinessIndustry523150      BusinessIndustry = "523150"
	BusinessIndustry523940      BusinessIndustry = "523940"
	BusinessIndustry523999      BusinessIndustry = "523999"
	BusinessIndustry524113      BusinessIndustry = "524113"
	BusinessIndustry813110      BusinessIndustry = "813110"
	BusinessIndustry813211      BusinessIndustry = "813211"
	BusinessIndustry813219      BusinessIndustry = "813219"
	BusinessIndustry721110      BusinessIndustry = "721110"
	BusinessIndustry722511      BusinessIndustry = "722511"
	BusinessIndustry722513      BusinessIndustry = "722513"
	BusinessIndustry561510      BusinessIndustry = "561510"
	BusinessIndustry713110      BusinessIndustry = "713110"
	BusinessIndustry713210      BusinessIndustry = "713210"
	BusinessIndustry712110      BusinessIndustry = "712110"
	BusinessIndustry711110      BusinessIndustry = "711110"
	BusinessIndustry711211      BusinessIndustry = "711211"
	BusinessIndustry621111      BusinessIndustry = "621111"
	BusinessIndustry621210      BusinessIndustry = "621210"
	BusinessIndustry622110      BusinessIndustry = "622110"
	BusinessIndustry623110      BusinessIndustry = "623110"
	BusinessIndustry621511      BusinessIndustry = "621511"
	BusinessIndustry623220      BusinessIndustry = "623220"
	BusinessIndustry541940      BusinessIndustry = "541940"
	BusinessIndustry621399      BusinessIndustry = "621399"
	BusinessIndustry621910      BusinessIndustry = "621910"
	BusinessIndustry541110      BusinessIndustry = "541110"
	BusinessIndustry311421      BusinessIndustry = "311421"
	BusinessIndustry337121      BusinessIndustry = "337121"
	BusinessIndustry322220      BusinessIndustry = "322220"
	BusinessIndustry339920      BusinessIndustry = "339920"
	BusinessIndustry334210      BusinessIndustry = "334210"
	BusinessIndustry339930      BusinessIndustry = "339930"
	BusinessIndustry312130      BusinessIndustry = "312130"
	BusinessIndustry336110      BusinessIndustry = "336110"
	BusinessIndustry339910      BusinessIndustry = "339910"
	BusinessIndustry516120      BusinessIndustry = "516120"
	BusinessIndustry513130      BusinessIndustry = "513130"
	BusinessIndustry512250      BusinessIndustry = "512250"
	BusinessIndustry519130      BusinessIndustry = "519130"
	BusinessIndustry711410      BusinessIndustry = "711410"
	BusinessIndustry711510      BusinessIndustry = "711510"
	BusinessIndustry531110      BusinessIndustry = "531110"
	BusinessIndustry531120      BusinessIndustry = "531120"
	BusinessIndustry531130      BusinessIndustry = "531130"
	BusinessIndustry531190      BusinessIndustry = "531190"
	BusinessIndustry531210      BusinessIndustry = "531210"
	BusinessIndustry531311      BusinessIndustry = "531311"
	BusinessIndustry531312      BusinessIndustry = "531312"
	BusinessIndustry531320      BusinessIndustry = "531320"
	BusinessIndustry531390      BusinessIndustry = "531390"
	BusinessIndustry445110      BusinessIndustry = "445110"
	BusinessIndustry455110      BusinessIndustry = "455110"
	BusinessIndustry457110      BusinessIndustry = "457110"
	BusinessIndustry449210      BusinessIndustry = "449210"
	BusinessIndustry444110      BusinessIndustry = "444110"
	BusinessIndustry459210      BusinessIndustry = "459210"
	BusinessIndustry459120      BusinessIndustry = "459120"
	BusinessIndustry445320      BusinessIndustry = "445320"
	BusinessIndustry458110      BusinessIndustry = "458110"
	BusinessIndustry458210      BusinessIndustry = "458210"
	BusinessIndustry458310      BusinessIndustry = "458310"
	BusinessIndustry455219      BusinessIndustry = "455219"
	BusinessIndustry456110      BusinessIndustry = "456110"
	BusinessIndustry517111      BusinessIndustry = "517111"
	BusinessIndustry517112      BusinessIndustry = "517112"
	BusinessIndustry517410      BusinessIndustry = "517410"
	BusinessIndustry481111      BusinessIndustry = "481111"
	BusinessIndustry483111      BusinessIndustry = "483111"
	BusinessIndustry485210      BusinessIndustry = "485210"
	BusinessIndustry423940      BusinessIndustry = "423940"
	BusinessIndustryDapp        BusinessIndustry = "dapp"
	BusinessIndustryExchange    BusinessIndustry = "exchange"
	BusinessIndustryGambling    BusinessIndustry = "gambling"
	BusinessIndustryGaming      BusinessIndustry = "gaming"
	BusinessIndustryInfra       BusinessIndustry = "infra"
	BusinessIndustryMarketplace BusinessIndustry = "marketplace"
	BusinessIndustryNeoBank     BusinessIndustry = "neo_bank"
	BusinessIndustryOther       BusinessIndustry = "other"
	BusinessIndustrySaas        BusinessIndustry = "saas"
	BusinessIndustrySocial      BusinessIndustry = "social"
	BusinessIndustryWallet      BusinessIndustry = "wallet"
	BusinessIndustry446120      BusinessIndustry = "446120"
)

type Decision string

const (
	DecisionApproved Decision = "approved"
	DecisionRejected Decision = "rejected"
)

type ReceiverType string

const (
	ReceiverTypeBusiness   ReceiverType = "business"
	ReceiverTypeIndividual ReceiverType = "individual"
)

type SwiftPaymentCode string

const (
	SwiftPaymentCodeHkSwiftCharitableDonation SwiftPaymentCode = "hk_swift_charitabledonation"
	SwiftPaymentCodeHkSwiftGoods              SwiftPaymentCode = "hk_swift_goods"
	SwiftPaymentCodeHkSwiftPersonal           SwiftPaymentCode = "hk_swift_personal"
	SwiftPaymentCodeHkSwiftServices           SwiftPaymentCode = "hk_swift_services"
)
