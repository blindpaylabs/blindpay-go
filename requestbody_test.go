package blindpay

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/customers"
	"github.com/blindpaylabs/blindpay-go/instances"
	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/payins"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/blindpaylabs/blindpay-go/quotes"
	"github.com/blindpaylabs/blindpay-go/transfers"
	"github.com/blindpaylabs/blindpay-go/virtualaccounts"
	"github.com/blindpaylabs/blindpay-go/wallets"
	"github.com/stretchr/testify/require"
)

func ptr[T any](v T) *T { return &v }

// TestRequestBodies pins the wire body of the request fields added for spec
// parity: optional fields are absent unless set (D23 "exactly one of"
// inputs included), new fields are sent when set.
func TestRequestBodies(t *testing.T) {
	ctx := context.Background()
	cases := []struct {
		name string
		body string
		call func(c *Client) error
	}{
		{"payout evm from managed wallet", `{"quote_id":"qu_1","wallet_id":"wa_1"}`, func(c *Client) error {
			_, err := c.Payouts.CreateEvm(ctx, &payouts.CreateEvmParams{QuoteID: "qu_1", WalletID: "wa_1"})
			return err
		}},
		{"payout evm from own wallet", `{"quote_id":"qu_1","sender_wallet_address":"0xabc"}`, func(c *Client) error {
			_, err := c.Payouts.CreateEvm(ctx, &payouts.CreateEvmParams{QuoteID: "qu_1", SenderWalletAddress: "0xabc"})
			return err
		}},
		{"payout stellar managed wallet", `{"quote_id":"qu_1","wallet_id":"wa_1"}`, func(c *Client) error {
			_, err := c.Payouts.CreateStellar(ctx, &payouts.CreateStellarParams{QuoteID: "qu_1", WalletID: "wa_1"})
			return err
		}},
		{"stellar authorize managed wallet", `{"quote_id":"qu_1","wallet_id":"wa_1"}`, func(c *Client) error {
			_, err := c.Payouts.AuthorizeStellarToken(ctx, &payouts.AuthorizeStellarTokenParams{QuoteID: "qu_1", WalletID: "wa_1"})
			return err
		}},
		{"payout solana managed wallet", `{"quote_id":"qu_1","wallet_id":"wa_1","signed_transaction":null}`, func(c *Client) error {
			_, err := c.Payouts.CreateSolana(ctx, &payouts.CreateSolanaParams{QuoteID: "qu_1", WalletID: "wa_1"})
			return err
		}},
		{"quote for a payable", `{"payable_id":"pa_1","network":"base","token":"USDC"}`, func(c *Client) error {
			n, tk := Network("base"), StablecoinToken("USDC")
			_, err := c.Quotes.Create(ctx, &quotes.CreateParams{PayableID: ptr("pa_1"), Network: &n, Token: &tk})
			return err
		}},
		{"payin quote ach pull", `{"wallet_id":"wa_1","currency_type":"sender","payment_method":"ach_pull","request_amount":1000,"token":"USDC","payer_rules":{"pix_allowed_tax_ids":null},"funding_bank_account_id":"ba_1"}`, func(c *Client) error {
			_, err := c.Payins.Quotes.Create(ctx, &payins.CreateQuoteParams{WalletID: ptr("wa_1"), CurrencyType: "sender", PaymentMethod: "ach_pull", RequestAmount: 1000, Token: "USDC", FundingBankAccountID: ptr("ba_1")})
			return err
		}},
		{"transfer quote partner fees", `{"wallet_id":"wa_1","sender_token":"USDC","customer_wallet_address":"0x1","customer_token":"USDC","customer_network":"base","request_amount":1000,"amount_reference":"sender","partner_percentage_fee":50,"partner_flat_fee":100}`, func(c *Client) error {
			_, err := c.Transfers.Quotes.Create(ctx, &transfers.CreateQuoteParams{WalletID: "wa_1", SenderToken: "USDC", CustomerWalletAddress: "0x1", CustomerToken: "USDC", CustomerNetwork: "base", RequestAmount: 1000, AmountReference: "sender", PartnerPercentageFee: ptr(int64(50)), PartnerFlatFee: ptr(int64(100))})
			return err
		}},
		{"solana delegation by quote", `{"quote_id":"qu_1","owner_address":"Own1"}`, func(c *Client) error {
			_, err := c.Wallets.PrepareSolanaDelegationTransaction(ctx, &wallets.PrepareSolanaDelegationTransactionParams{QuoteID: "qu_1", OwnerAddress: "Own1"})
			return err
		}},
		{"virtual account bre-b", `{"blockchain_wallet_id":"bw_1","token":"USDC","banking_partner":"bre_b","partner_fee_id":"fe_1","bre_b":{"gender":"female"}}`, func(c *Client) error {
			_, err := c.VirtualAccounts.Create(ctx, &virtualaccounts.CreateParams{CustomerID: tCustomer, BlockchainWalletID: "bw_1", Token: "USDC", BankingPartner: "bre_b", PartnerFeeID: ptr("fe_1"), BreB: &virtualaccounts.BreBParams{Gender: virtualaccounts.BreBGenderFemale}})
			return err
		}},
		{"offramp wallet partner fee", `{"external_id":"x","network":"tron","partner_fee_id":"fe_1"}`, func(c *Client) error {
			_, err := c.OfframpWallets.Create(ctx, &wallets.CreateOfframpWalletParams{CustomerID: tCustomer, BankAccountID: tBankAccount, ExternalID: "x", Network: "tron", PartnerFeeID: ptr("fe_1")})
			return err
		}},
		{"instance update branding", `{"name":"Acme","logo_url":"https://x/l.png","accent_color":"#000000","compliance_emails":["c@example.com"],"partner_tos_terms_accepted":true}`, func(c *Client) error {
			return c.Instances.Update(ctx, &instances.UpdateParams{Name: "Acme", LogoURL: ptr("https://x/l.png"), AccentColor: ptr("#000000"), ComplianceEmails: []string{"c@example.com"}, PartnerTosTermsAccepted: ptr(true)})
		}},
		{"limit increase several documents", `{"daily":100000,"monthly":500000,"per_transaction":50000,"supporting_documents":[{"type":"individual_bank_statement","file":"https://x/1.pdf"}]}`, func(c *Client) error {
			_, err := c.Customers.RequestLimitIncrease(ctx, &customers.RequestLimitIncreaseParams{CustomerID: tCustomer, Daily: 100000, Monthly: 500000, PerTransaction: 50000,
				SupportingDocuments: []customers.LimitIncreaseSupportingDocument{{Type: "individual_bank_statement", File: "https://x/1.pdf"}}})
			return err
		}},
		{"customer update location", `{"latitude":-23.5,"longitude":-46.6,"additional_info":[{"label":"BANK_STATEMENT","value":"https://x/1.pdf"}]}`, func(c *Client) error {
			return c.Customers.Update(ctx, &customers.UpdateParams{CustomerID: tCustomer, Latitude: ptr(-23.5), Longitude: ptr(-46.6),
				AdditionalInfo: []customers.AdditionalInfoItem{{Label: "BANK_STATEMENT", Value: "https://x/1.pdf"}}})
		}},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			c := rtClient(t, &blindpaytest.RoundTripper{T: t, In: json.RawMessage(tc.body), Out: json.RawMessage(`{}`), Status: http.StatusOK})
			require.NoError(t, tc.call(c))
		})
	}
}
