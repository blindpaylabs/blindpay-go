package types

import "time"

type CurrencyType string

const (
	CurrencyTypeSender   CurrencyType = "sender"
	CurrencyTypeReceiver CurrencyType = "receiver"
)

type Network string

const (
	NetworkBase            Network = "base"
	NetworkSepolia         Network = "sepolia"
	NetworkArbitrumSepolia Network = "arbitrum_sepolia"
	NetworkBaseSepolia     Network = "base_sepolia"
	NetworkArbitrum        Network = "arbitrum"
	NetworkPolygon         Network = "polygon"
	NetworkPolygonAmoy     Network = "polygon_amoy"
	NetworkEthereum        Network = "ethereum"
	NetworkStellar         Network = "stellar"
	NetworkStellarTestnet  Network = "stellar_testnet"
	NetworkTron            Network = "tron"
)

type StablecoinToken string

const (
	StablecoinTokenUSDC StablecoinToken = "USDC"
	StablecoinTokenUSDT StablecoinToken = "USDT"
	StablecoinTokenUSDB StablecoinToken = "USDB"
)

type TransactionDocumentType string

const (
	TransactionDocumentTypeInvoice            TransactionDocumentType = "invoice"
	TransactionDocumentTypePurchaseOrder      TransactionDocumentType = "purchase_order"
	TransactionDocumentTypeDeliverySlip       TransactionDocumentType = "delivery_slip"
	TransactionDocumentTypeContract           TransactionDocumentType = "contract"
	TransactionDocumentTypeCustomsDeclaration TransactionDocumentType = "customs_declaration"
	TransactionDocumentTypeBillOfLading       TransactionDocumentType = "bill_of_lading"
	TransactionDocumentTypeOthers             TransactionDocumentType = "others"
)

type BankAccountType string

const (
	BankAccountTypeChecking BankAccountType = "checking"
	BankAccountTypeSavings  BankAccountType = "savings"
)

type Currency string

const (
	CurrencyUSDC Currency = "USDC"
	CurrencyUSDT Currency = "USDT"
	CurrencyUSDB Currency = "USDB"
	CurrencyBRL  Currency = "BRL"
	CurrencyUSD  Currency = "USD"
	CurrencyMXN  Currency = "MXN"
	CurrencyCOP  Currency = "COP"
	CurrencyARS  Currency = "ARS"
)

type AccountClass string

const (
	AccountClassIndividual AccountClass = "individual"
	AccountClassBusiness   AccountClass = "business"
)

type TransactionStatus string

const (
	TransactionStatusPending    TransactionStatus = "pending"
	TransactionStatusRefunded   TransactionStatus = "refunded"
	TransactionStatusProcessing TransactionStatus = "processing"
	TransactionStatusCompleted  TransactionStatus = "completed"
	TransactionStatusFailed     TransactionStatus = "failed"
	TransactionStatusOnHold     TransactionStatus = "on_hold"
)

type PaginationMetadata struct {
	HasMore  bool `json:"has_more"`
	NextPage int  `json:"next_page"`
	PrevPage int  `json:"prev_page"`
}

type TrackingTransaction struct {
	Step                 string    `json:"step"`
	Status               string    `json:"status"`
	TransactionHash      string    `json:"transaction_hash,omitempty"`
	ExternalID           string    `json:"external_id,omitempty"`
	CompletedAt          time.Time `json:"completed_at"`
	SenderName           string    `json:"sender_name,omitempty"`
	SenderTaxID          string    `json:"sender_tax_id,omitempty"`
	SenderBankCode       string    `json:"sender_bank_code,omitempty"`
	SenderAccountNumber  string    `json:"sender_account_number,omitempty"`
	TraceNumber          string    `json:"trace_number,omitempty"`
	TransactionReference string    `json:"transaction_reference,omitempty"`
	Description          string    `json:"description,omitempty"`
}

type TrackingPayment struct {
	Step                   string    `json:"step"`
	ProviderName           string    `json:"provider_name"`
	ProviderTransactionID  string    `json:"provider_transaction_id,omitempty"`
	ProviderStatus         string    `json:"provider_status,omitempty"`
	EstimatedTimeOfArrival string    `json:"estimated_time_of_arrival,omitempty"`
	CompletedAt            time.Time `json:"completed_at"`
}

type TrackingLiquidity struct {
	Step                   string    `json:"step"`
	ProviderTransactionID  string    `json:"provider_transaction_id"`
	ProviderStatus         string    `json:"provider_status"`
	EstimatedTimeOfArrival string    `json:"estimated_time_of_arrival"`
	CompletedAt            time.Time `json:"completed_at"`
}

type TrackingComplete struct {
	Step            string    `json:"step"`
	Status          string    `json:"status"`
	TransactionHash string    `json:"transaction_hash"`
	CompletedAt     time.Time `json:"completed_at"`
}

type TrackingPartnerFee struct {
	Step            string    `json:"step"`
	TransactionHash string    `json:"transaction_hash"`
	CompletedAt     time.Time `json:"completed_at"`
}

type Country string

const (
	CountryAF Country = "AF"
	CountryAL Country = "AL"
	CountryDZ Country = "DZ"
	CountryAS Country = "AS"
	CountryAD Country = "AD"
	CountryAO Country = "AO"
	CountryAI Country = "AI"
	CountryAQ Country = "AQ"
	CountryAG Country = "AG"
	CountryAR Country = "AR"
	CountryAM Country = "AM"
	CountryAW Country = "AW"
	CountryAU Country = "AU"
	CountryAT Country = "AT"
	CountryAZ Country = "AZ"
	CountryBS Country = "BS"
	CountryBH Country = "BH"
	CountryBD Country = "BD"
	CountryBB Country = "BB"
	CountryBY Country = "BY"
	CountryBE Country = "BE"
	CountryBZ Country = "BZ"
	CountryBJ Country = "BJ"
	CountryBM Country = "BM"
	CountryBT Country = "BT"
	CountryBO Country = "BO"
	CountryBQ Country = "BQ"
	CountryBA Country = "BA"
	CountryBW Country = "BW"
	CountryBV Country = "BV"
	CountryBR Country = "BR"
	CountryIO Country = "IO"
	CountryBN Country = "BN"
	CountryBG Country = "BG"
	CountryBF Country = "BF"
	CountryBI Country = "BI"
	CountryCV Country = "CV"
	CountryKH Country = "KH"
	CountryCM Country = "CM"
	CountryCA Country = "CA"
	CountryKY Country = "KY"
	CountryCF Country = "CF"
	CountryTD Country = "TD"
	CountryCL Country = "CL"
	CountryCN Country = "CN"
	CountryCX Country = "CX"
	CountryCC Country = "CC"
	CountryCO Country = "CO"
	CountryKM Country = "KM"
	CountryCD Country = "CD"
	CountryCG Country = "CG"
	CountryCK Country = "CK"
	CountryCR Country = "CR"
	CountryHR Country = "HR"
	CountryCU Country = "CU"
	CountryCW Country = "CW"
	CountryCY Country = "CY"
	CountryCZ Country = "CZ"
	CountryCI Country = "CI"
	CountryDK Country = "DK"
	CountryDJ Country = "DJ"
	CountryDM Country = "DM"
	CountryDO Country = "DO"
	CountryEC Country = "EC"
	CountryEG Country = "EG"
	CountrySV Country = "SV"
	CountryGQ Country = "GQ"
	CountryER Country = "ER"
	CountryEE Country = "EE"
	CountrySZ Country = "SZ"
	CountryET Country = "ET"
	CountryFK Country = "FK"
	CountryFO Country = "FO"
	CountryFJ Country = "FJ"
	CountryFI Country = "FI"
	CountryFR Country = "FR"
	CountryGF Country = "GF"
	CountryPF Country = "PF"
	CountryTF Country = "TF"
	CountryGA Country = "GA"
	CountryGM Country = "GM"
	CountryGE Country = "GE"
	CountryDE Country = "DE"
	CountryGH Country = "GH"
	CountryGI Country = "GI"
	CountryGR Country = "GR"
	CountryGL Country = "GL"
	CountryGD Country = "GD"
	CountryGP Country = "GP"
	CountryGU Country = "GU"
	CountryGT Country = "GT"
	CountryGG Country = "GG"
	CountryGN Country = "GN"
	CountryGW Country = "GW"
	CountryGY Country = "GY"
	CountryHT Country = "HT"
	CountryHM Country = "HM"
	CountryVA Country = "VA"
	CountryHN Country = "HN"
	CountryHK Country = "HK"
	CountryHU Country = "HU"
	CountryIS Country = "IS"
	CountryIN Country = "IN"
	CountryID Country = "ID"
	CountryIR Country = "IR"
	CountryIQ Country = "IQ"
	CountryIE Country = "IE"
	CountryIM Country = "IM"
	CountryIL Country = "IL"
	CountryIT Country = "IT"
	CountryJM Country = "JM"
	CountryJP Country = "JP"
	CountryJE Country = "JE"
	CountryJO Country = "JO"
	CountryKZ Country = "KZ"
	CountryKE Country = "KE"
	CountryKI Country = "KI"
	CountryKP Country = "KP"
	CountryKR Country = "KR"
	CountryKW Country = "KW"
	CountryKG Country = "KG"
	CountryLA Country = "LA"
	CountryLV Country = "LV"
	CountryLB Country = "LB"
	CountryLS Country = "LS"
	CountryLR Country = "LR"
	CountryLY Country = "LY"
	CountryLI Country = "LI"
	CountryLT Country = "LT"
	CountryLU Country = "LU"
	CountryMO Country = "MO"
	CountryMG Country = "MG"
	CountryMW Country = "MW"
	CountryMY Country = "MY"
	CountryMV Country = "MV"
	CountryML Country = "ML"
	CountryMT Country = "MT"
	CountryMH Country = "MH"
	CountryMQ Country = "MQ"
	CountryMR Country = "MR"
	CountryMU Country = "MU"
	CountryYT Country = "YT"
	CountryMX Country = "MX"
	CountryFM Country = "FM"
	CountryMD Country = "MD"
	CountryMC Country = "MC"
	CountryMN Country = "MN"
	CountryME Country = "ME"
	CountryMS Country = "MS"
	CountryMA Country = "MA"
	CountryMZ Country = "MZ"
	CountryMM Country = "MM"
	CountryNA Country = "NA"
	CountryNR Country = "NR"
	CountryNP Country = "NP"
	CountryNL Country = "NL"
	CountryNC Country = "NC"
	CountryNZ Country = "NZ"
	CountryNI Country = "NI"
	CountryNE Country = "NE"
	CountryNG Country = "NG"
	CountryNU Country = "NU"
	CountryNF Country = "NF"
	CountryMP Country = "MP"
	CountryNO Country = "NO"
	CountryOM Country = "OM"
	CountryPK Country = "PK"
	CountryPW Country = "PW"
	CountryPS Country = "PS"
	CountryPA Country = "PA"
	CountryPG Country = "PG"
	CountryPY Country = "PY"
	CountryPE Country = "PE"
	CountryPH Country = "PH"
	CountryPN Country = "PN"
	CountryPL Country = "PL"
	CountryPT Country = "PT"
	CountryPR Country = "PR"
	CountryQA Country = "QA"
	CountryMK Country = "MK"
	CountryRO Country = "RO"
	CountryRU Country = "RU"
	CountryRW Country = "RW"
	CountryRE Country = "RE"
	CountryBL Country = "BL"
	CountrySH Country = "SH"
	CountryKN Country = "KN"
	CountryLC Country = "LC"
	CountryMF Country = "MF"
	CountryPM Country = "PM"
	CountryVC Country = "VC"
	CountryWS Country = "WS"
	CountrySM Country = "SM"
	CountryST Country = "ST"
	CountrySA Country = "SA"
	CountrySN Country = "SN"
	CountryRS Country = "RS"
	CountrySC Country = "SC"
	CountrySL Country = "SL"
	CountrySG Country = "SG"
	CountrySX Country = "SX"
	CountrySK Country = "SK"
	CountrySI Country = "SI"
	CountrySB Country = "SB"
	CountrySO Country = "SO"
	CountryZA Country = "ZA"
	CountryGS Country = "GS"
	CountrySS Country = "SS"
	CountryES Country = "ES"
	CountryLK Country = "LK"
	CountrySD Country = "SD"
	CountrySR Country = "SR"
	CountrySJ Country = "SJ"
	CountrySE Country = "SE"
	CountryCH Country = "CH"
	CountrySY Country = "SY"
	CountryTW Country = "TW"
	CountryTJ Country = "TJ"
	CountryTZ Country = "TZ"
	CountryTH Country = "TH"
	CountryTL Country = "TL"
	CountryTG Country = "TG"
	CountryTK Country = "TK"
	CountryTO Country = "TO"
	CountryTT Country = "TT"
	CountryTN Country = "TN"
	CountryTR Country = "TR"
	CountryTM Country = "TM"
	CountryTC Country = "TC"
	CountryTV Country = "TV"
	CountryUG Country = "UG"
	CountryUA Country = "UA"
	CountryAE Country = "AE"
	CountryGB Country = "GB"
	CountryUM Country = "UM"
	CountryUS Country = "US"
	CountryUY Country = "UY"
	CountryUZ Country = "UZ"
	CountryVU Country = "VU"
	CountryVE Country = "VE"
	CountryVN Country = "VN"
	CountryVG Country = "VG"
	CountryVI Country = "VI"
	CountryWF Country = "WF"
	CountryEH Country = "EH"
	CountryYE Country = "YE"
	CountryZM Country = "ZM"
	CountryZW Country = "ZW"
	CountryAX Country = "AX"
)
