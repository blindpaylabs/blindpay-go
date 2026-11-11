package blindpay

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"testing"

	"github.com/blindpaylabs/blindpay-go/bankaccounts"
	"github.com/blindpaylabs/blindpay-go/custodialwallets"
	"github.com/blindpaylabs/blindpay-go/customers"
	"github.com/blindpaylabs/blindpay-go/instances"
	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/partnerfees"
	"github.com/blindpaylabs/blindpay-go/payables"
	"github.com/blindpaylabs/blindpay-go/payins"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/blindpaylabs/blindpay-go/quotes"
	"github.com/blindpaylabs/blindpay-go/termsofservice"
	"github.com/blindpaylabs/blindpay-go/transfers"
	"github.com/blindpaylabs/blindpay-go/upload"
	"github.com/blindpaylabs/blindpay-go/virtualaccounts"
	"github.com/blindpaylabs/blindpay-go/wallets"
	"github.com/blindpaylabs/blindpay-go/webhookendpoints"
	"github.com/stretchr/testify/require"
)

// Placeholder ids used both to call the SDK and to build the expected path.
const (
	tInstance    = "in_000000000000"
	tCustomer    = "cu_000000000000"
	tID          = "id_000000000000"
	tBankAccount = "ba_000000000000"
	tUser        = "us_000000000000"
	tSwift       = "AAAABBCCXXX"
)

// opCall invokes one SDK method and returns its error. The decoded result
// is discarded: decoding failures surface as errors.
type opCall func(ctx context.Context, c *Client) error

// opEntry maps one API operation to the SDK methods that call it, plus the
// metadata printed in the cross-SDK coverage table.
type opEntry struct {
	// Name is the primary SDK method (first entry of Calls), e.g.
	// "Payouts.List".
	Name string
	// Pager names the auto-pager, if the op has one.
	Pager string
	// Notes is free text for the coverage table.
	Notes string
	// Calls are every SDK method for this op; each must hit the op's method
	// and path and decode every fixture variant.
	Calls []opCall
}

func call[T any](f func(ctx context.Context, c *Client) (T, error)) opCall {
	return func(ctx context.Context, c *Client) error {
		_, err := f(ctx, c)
		return err
	}
}

// opTable has one entry per op in testdata/fixtures/index.json (the 94
// x-sdk operations of the reference spec).
var opTable = map[string]opEntry{
	"GET /v1/available/bank-details": {Name: "Available.GetBankDetails", Calls: []opCall{
		call(func(ctx context.Context, c *Client) ([]BankDetail, error) {
			return c.Available.GetBankDetails(ctx, "pix")
		}),
	}},
	"GET /v1/available/naics": {Name: "Available.GetNaicsCodes", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Available.GetNaicsCodes(ctx) }),
	}},
	"GET /v1/available/rails": {Name: "Available.GetRails", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Available.GetRails(ctx) }),
	}},
	"GET /v1/available/swift/{swift}": {Name: "Available.GetSwiftCodeBankDetails", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Available.GetSwiftCodeBankDetails(ctx, tSwift)
		}),
	}},

	// Bank accounts.
	"POST /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts": {Name: "BankAccounts.CreateAch", Notes: "one method per rail: CreatePix, CreateAch, CreateWire, CreateArgentinaTransfers, CreateSpei, CreateColombiaAch, CreateInternationalSwift, CreatePixSafe, CreateRtp, CreateTed, CreateSepa", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateAch(ctx, &bankaccounts.CreateAchParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreatePix(ctx, &bankaccounts.CreatePixParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateWire(ctx, &bankaccounts.CreateWireParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateArgentinaTransfers(ctx, &bankaccounts.CreateArgentinaTransfersParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateSpei(ctx, &bankaccounts.CreateSpeiParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateColombiaAch(ctx, &bankaccounts.CreateColombiaAchParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateInternationalSwift(ctx, &bankaccounts.CreateInternationalSwiftParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreatePixSafe(ctx, &bankaccounts.CreatePixSafeParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateRtp(ctx, &bankaccounts.CreateRtpParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateTed(ctx, &bankaccounts.CreateTedParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.CreateSepa(ctx, &bankaccounts.CreateSepaParams{CustomerID: tCustomer})
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts": {Name: "BankAccounts.List", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.BankAccounts.List(ctx, &bankaccounts.ListParams{CustomerID: tCustomer})
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts/{id}": {Name: "BankAccounts.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.BankAccounts.Get(ctx, tCustomer, tID) }),
	}},
	"DELETE /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts/{id}": {Name: "BankAccounts.Delete", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.BankAccounts.Delete(ctx, tCustomer, tID) },
	}},

	// Blockchain wallets.
	"POST /v1/instances/{instance_id}/customers/{customer_id}/blockchain-wallets": {Name: "Wallets.CreateWithAddress", Notes: "also Wallets.CreateWithHash", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Wallets.CreateWithAddress(ctx, &wallets.CreateWithAddressParams{CustomerID: tCustomer})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Wallets.CreateWithHash(ctx, &wallets.CreateWithHashParams{CustomerID: tCustomer})
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/blockchain-wallets": {Name: "Wallets.List", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Wallets.List(ctx, tCustomer) }),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/blockchain-wallets/sign-message": {Name: "Wallets.GetWalletMessage", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Wallets.GetWalletMessage(ctx, tCustomer) }),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/blockchain-wallets/{id}": {Name: "Wallets.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Wallets.Get(ctx, tCustomer, tID) }),
	}},
	"DELETE /v1/instances/{instance_id}/customers/{customer_id}/blockchain-wallets/{id}": {Name: "Wallets.Delete", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.Wallets.Delete(ctx, tCustomer, tID) },
	}},

	// Customers.
	"POST /v1/instances/{instance_id}/customers": {Name: "Customers.CreateIndividualWithStandardKYC", Notes: "also CreateIndividualWithEnhancedKYC, CreateBusinessWithStandardKYB", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Customers.CreateIndividualWithStandardKYC(ctx, &customers.CreateIndividualStandardParams{})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Customers.CreateIndividualWithEnhancedKYC(ctx, &customers.CreateIndividualEnhancedParams{})
		}),
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Customers.CreateBusinessWithStandardKYB(ctx, &customers.CreateBusinessStandardParams{})
		}),
	}},
	"GET /v1/instances/{instance_id}/customers": {Name: "Customers.ListWithParams", Pager: "Customers.ListAutoPaging", Notes: "also Customers.List, Customers.ListPage", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Customers.ListWithParams(ctx, &customers.ListParams{})
		}),
		call(func(ctx context.Context, c *Client) (any, error) { return c.Customers.List(ctx) }),
		call(func(ctx context.Context, c *Client) (any, error) { return c.Customers.ListPage(ctx, nil) }),
	}},
	"POST /v1/instances/{instance_id}/customers/{customer_id}/limit-increase": {Name: "Customers.RequestLimitIncrease", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Customers.RequestLimitIncrease(ctx, &customers.RequestLimitIncreaseParams{CustomerID: tCustomer})
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/limit-increase": {Name: "Customers.GetLimitIncreaseRequests", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Customers.GetLimitIncreaseRequests(ctx, tCustomer)
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{id}": {Name: "Customers.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Customers.Get(ctx, tID) }),
	}},
	"PUT /v1/instances/{instance_id}/customers/{id}": {Name: "Customers.Update", Calls: []opCall{
		func(ctx context.Context, c *Client) error {
			return c.Customers.Update(ctx, &customers.UpdateParams{CustomerID: tID})
		},
	}},
	"DELETE /v1/instances/{instance_id}/customers/{id}": {Name: "Customers.Delete", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.Customers.Delete(ctx, tID) },
	}},
	"GET /v1/instances/{instance_id}/limits/customers/{id}": {Name: "Customers.GetLimits", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Customers.GetLimits(ctx, tID) }),
	}},
	"GET /v1/instances/{instance_id}/billing/fees": {Name: "Fees.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Fees.Get(ctx) }),
	}},

	// Instances.
	"PUT /v1/instances/{id}": {Name: "Instances.Update", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.Instances.Update(ctx, &instances.UpdateParams{}) },
	}},
	"DELETE /v1/instances/{id}": {Name: "Instances.Delete", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.Instances.Delete(ctx) },
	}},
	"GET /v1/instances/{id}/members": {Name: "Instances.GetMembers", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Instances.GetMembers(ctx) }),
	}},
	"DELETE /v1/instances/{id}/members/{user_id}": {Name: "Instances.DeleteMember", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.Instances.DeleteMember(ctx, tUser) },
	}},
	"PUT /v1/instances/{id}/members/{user_id}": {Name: "Instances.UpdateMemberRole", Calls: []opCall{
		func(ctx context.Context, c *Client) error {
			return c.Instances.UpdateMemberRole(ctx, &instances.UpdateMemberRoleParams{MemberID: tUser})
		},
	}},
	"POST /v1/instances/{id}/ownership": {Name: "Instances.MigrateOwnership", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Instances.MigrateOwnership(ctx, &instances.MigrateOwnershipParams{UserID: tUser})
		}),
	}},

	// Offramp wallets.
	"POST /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts/{bank_account_id}/offramp-wallets": {Name: "OfframpWallets.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.OfframpWallets.Create(ctx, &wallets.CreateOfframpWalletParams{CustomerID: tCustomer, BankAccountID: tBankAccount})
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts/{bank_account_id}/offramp-wallets": {Name: "OfframpWallets.List", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.OfframpWallets.List(ctx, tCustomer, tBankAccount)
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts/{bank_account_id}/offramp-wallets/{id}": {Name: "OfframpWallets.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.OfframpWallets.Get(ctx, tCustomer, tBankAccount, tID)
		}),
	}},

	// Partner fees.
	"POST /v1/instances/{instance_id}/partner-fees": {Name: "PartnerFees.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.PartnerFees.Create(ctx, &partnerfees.CreatePartnerFeeParams{})
		}),
	}},
	"GET /v1/instances/{instance_id}/partner-fees": {Name: "PartnerFees.List", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.PartnerFees.List(ctx) }),
	}},
	"GET /v1/instances/{instance_id}/partner-fees/{id}": {Name: "PartnerFees.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.PartnerFees.Get(ctx, tID) }),
	}},
	"DELETE /v1/instances/{instance_id}/partner-fees/{id}": {Name: "PartnerFees.Delete", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.PartnerFees.Delete(ctx, tID) },
	}},

	// Payables.
	"POST /v1/instances/{instance_id}/payables": {Name: "Payables.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Payables.Create(ctx, &payables.CreateParams{})
		}),
	}},
	"GET /v1/instances/{instance_id}/payables": {Name: "Payables.List", Pager: "Payables.ListAutoPaging", Notes: "also Payables.ListPage", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payables.List(ctx, nil) }),
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payables.ListPage(ctx, nil) }),
	}},
	"GET /v1/instances/{instance_id}/payables/{id}": {Name: "Payables.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payables.Get(ctx, tID) }),
	}},
	"DELETE /v1/instances/{instance_id}/payables/{id}": {Name: "Payables.Delete", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payables.Delete(ctx, tID) }),
	}},

	// Payins and payin quotes.
	"POST /v1/instances/{instance_id}/payin-quotes": {Name: "Payins.Quotes.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Payins.Quotes.Create(ctx, &payins.CreateQuoteParams{})
		}),
	}},
	"POST /v1/instances/{instance_id}/payin-quotes/fx": {Name: "Payins.Quotes.GetFxRate", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Payins.Quotes.GetFxRate(ctx, &payins.GetFxRateParams{})
		}),
	}},
	"GET /v1/e/payins/{id}": {Name: "Payins.GetTrack", Notes: "public tracking endpoint; 401 from server-side calls in production (D19)", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payins.GetTrack(ctx, tID) }),
	}},
	"GET /v1/instances/{instance_id}/payins": {Name: "Payins.List", Pager: "Payins.ListAutoPaging", Notes: "also Payins.ListPage", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payins.List(ctx, nil) }),
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payins.ListPage(ctx, nil) }),
	}},
	"POST /v1/instances/{instance_id}/payins/evm": {Name: "Payins.CreateEvm", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payins.CreateEvm(ctx, "pq_000000000000") }),
	}},
	"GET /v1/instances/{instance_id}/payins/{id}": {Name: "Payins.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payins.Get(ctx, tID) }),
	}},

	// Payouts.
	"GET /v1/e/payouts/{id}": {Name: "Payouts.GetTrack", Notes: "public tracking endpoint; 401 from server-side calls in production (D19)", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payouts.GetTrack(ctx, tID) }),
	}},
	"GET /v1/instances/{instance_id}/payouts": {Name: "Payouts.List", Pager: "Payouts.ListAutoPaging", Notes: "also Payouts.ListPage", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payouts.List(ctx, nil) }),
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payouts.ListPage(ctx, nil) }),
	}},
	"POST /v1/instances/{instance_id}/payouts/evm": {Name: "Payouts.CreateEvm", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Payouts.CreateEvm(ctx, &payouts.CreateEvmParams{})
		}),
	}},
	"POST /v1/instances/{instance_id}/payouts/solana": {Name: "Payouts.CreateSolana", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Payouts.CreateSolana(ctx, &payouts.CreateSolanaParams{})
		}),
	}},
	"POST /v1/instances/{instance_id}/payouts/stellar": {Name: "Payouts.CreateStellar", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Payouts.CreateStellar(ctx, &payouts.CreateStellarParams{})
		}),
	}},
	"POST /v1/instances/{instance_id}/payouts/stellar/authorize": {Name: "Payouts.AuthorizeStellarToken", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Payouts.AuthorizeStellarToken(ctx, &payouts.AuthorizeStellarTokenParams{})
		}),
	}},
	"GET /v1/instances/{instance_id}/payouts/{id}": {Name: "Payouts.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Payouts.Get(ctx, tID) }),
	}},
	"POST /v1/instances/{instance_id}/payouts/{id}/documents": {Name: "Payouts.SubmitDocuments", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Payouts.SubmitDocuments(ctx, &payouts.SubmitDocumentsParams{PayoutID: tID})
		}),
	}},

	// Payout quotes.
	"POST /v1/instances/{instance_id}/quotes": {Name: "Quotes.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Quotes.Create(ctx, &quotes.CreateParams{}) }),
	}},
	"POST /v1/instances/{instance_id}/quotes/fx": {Name: "Quotes.GetFxRate", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Quotes.GetFxRate(ctx, &quotes.GetFxRateParams{})
		}),
	}},

	// Wallet helpers (Solana / Stellar).
	"POST /v1/instances/{instance_id}/mint-usdb-solana": {Name: "Wallets.MintUsdbSolana", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Wallets.MintUsdbSolana(ctx, &wallets.MintUsdbSolanaParams{})
		}),
	}},
	"POST /v1/instances/{instance_id}/prepare-delegate-solana": {Name: "Wallets.PrepareSolanaDelegationTransaction", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Wallets.PrepareSolanaDelegationTransaction(ctx, &wallets.PrepareSolanaDelegationTransactionParams{})
		}),
	}},
	"POST /v1/instances/{instance_id}/create-asset-trustline": {Name: "Wallets.CreateAssetTrustline", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Wallets.CreateAssetTrustline(ctx, "GABC") }),
	}},
	"POST /v1/instances/{instance_id}/mint-usdb-stellar": {Name: "Wallets.MintUsdbStellar", Calls: []opCall{
		func(ctx context.Context, c *Client) error {
			return c.Wallets.MintUsdbStellar(ctx, &wallets.MintUsdbStellarParams{})
		},
	}},

	// Terms of service.
	"POST /v1/e/instances/{instance_id}/tos": {Name: "Tos.Initiate", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Tos.Initiate(ctx, &termsofservice.InitiateParams{IdempotencyKey: "4c9a5b6e-0000-4000-8000-000000000000"})
		}),
	}},

	// Transfers.
	"POST /v1/instances/{instance_id}/transfer-quotes": {Name: "Transfers.Quotes.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Transfers.Quotes.Create(ctx, &transfers.CreateQuoteParams{})
		}),
	}},
	"GET /v1/e/transfers/{id}": {Name: "Transfers.GetTrack", Notes: "public tracking endpoint; 401 from server-side calls in production (D19)", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Transfers.GetTrack(ctx, tID) }),
	}},
	"GET /v1/instances/{instance_id}/transfers": {Name: "Transfers.List", Pager: "Transfers.ListAutoPaging", Notes: "also Transfers.ListPage", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Transfers.List(ctx, nil) }),
		call(func(ctx context.Context, c *Client) (any, error) { return c.Transfers.ListPage(ctx, nil) }),
	}},
	"POST /v1/instances/{instance_id}/transfers": {Name: "Transfers.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Transfers.Create(ctx, &transfers.CreateParams{})
		}),
	}},
	"GET /v1/instances/{instance_id}/transfers/{id}": {Name: "Transfers.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.Transfers.Get(ctx, tID) }),
	}},

	// Upload.
	"POST /v1/upload": {Name: "Upload.Upload", Notes: "multipart; never retried, no idempotency key", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.Upload.Upload(ctx, &upload.UploadParams{File: strings.NewReader("x"), FileName: "a.pdf", Bucket: upload.BucketOnboarding})
		}),
	}},

	// Virtual accounts.
	"POST /v1/instances/{instance_id}/customers/{customer_id}/virtual-accounts": {Name: "VirtualAccounts.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.VirtualAccounts.Create(ctx, &virtualaccounts.CreateParams{CustomerID: tCustomer})
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/virtual-accounts": {Name: "VirtualAccounts.List", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.VirtualAccounts.List(ctx, tCustomer) }),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/virtual-accounts/{id}": {Name: "VirtualAccounts.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.VirtualAccounts.Get(ctx, tCustomer, tID) }),
	}},
	"PUT /v1/instances/{instance_id}/customers/{customer_id}/virtual-accounts/{id}": {Name: "VirtualAccounts.Update", Calls: []opCall{
		func(ctx context.Context, c *Client) error {
			return c.VirtualAccounts.Update(ctx, &virtualaccounts.UpdateParams{CustomerID: tCustomer, VirtualAccountID: tID})
		},
	}},

	// Custodial (managed) wallets.
	"POST /v1/instances/{instance_id}/customers/{customer_id}/wallets": {Name: "CustodialWallets.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.CustodialWallets.Create(ctx, &custodialwallets.CreateParams{CustomerID: tCustomer})
		}),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/wallets": {Name: "CustodialWallets.List", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.CustodialWallets.List(ctx, tCustomer) }),
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/wallets/{id}": {Name: "CustodialWallets.Get", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.CustodialWallets.Get(ctx, tCustomer, tID) }),
	}},
	"DELETE /v1/instances/{instance_id}/customers/{customer_id}/wallets/{id}": {Name: "CustodialWallets.Delete", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.CustodialWallets.Delete(ctx, tCustomer, tID) },
	}},
	"GET /v1/instances/{instance_id}/customers/{customer_id}/wallets/{id}/balance": {Name: "CustodialWallets.GetBalance", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.CustodialWallets.GetBalance(ctx, tCustomer, tID)
		}),
	}},

	// Webhook endpoints.
	"POST /v1/instances/{instance_id}/webhook-endpoints": {Name: "WebhookEndpoints.Create", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) {
			return c.WebhookEndpoints.Create(ctx, &webhookendpoints.CreateParams{})
		}),
	}},
	"GET /v1/instances/{instance_id}/webhook-endpoints": {Name: "WebhookEndpoints.List", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.WebhookEndpoints.List(ctx) }),
	}},
	"GET /v1/instances/{instance_id}/webhook-endpoints/portal-access": {Name: "WebhookEndpoints.GetPortalAccessURL", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.WebhookEndpoints.GetPortalAccessURL(ctx) }),
	}},
	"DELETE /v1/instances/{instance_id}/webhook-endpoints/{id}": {Name: "WebhookEndpoints.Delete", Calls: []opCall{
		func(ctx context.Context, c *Client) error { return c.WebhookEndpoints.Delete(ctx, tID) },
	}},
	"GET /v1/instances/{instance_id}/webhook-endpoints/{id}/secret": {Name: "WebhookEndpoints.GetSecret", Calls: []opCall{
		call(func(ctx context.Context, c *Client) (any, error) { return c.WebhookEndpoints.GetSecret(ctx, tID) }),
	}},
}

// pendingOps are surface ops the SDK does not implement yet. The test
// fails if one of them gains a table entry without leaving this set, so it
// can only shrink.
var pendingOps = map[string]bool{
	"POST /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts/plaid":      true,
	"PATCH /v1/instances/{instance_id}/customers/{customer_id}/blockchain-wallets/{id}": true,
	"GET /v1/instances/{instance_id}/customers/{customer_id}/rfi":                       true,
	"POST /v1/instances/{instance_id}/customers/{customer_id}/rfi":                      true,
	"GET /v1/instances/{id}": true,
	"GET /v1/instances/{instance_id}/onboarding/business_details":                                                  true,
	"GET /v1/instances/{instance_id}/onboarding/business_profile":                                                  true,
	"GET /v1/instances/{instance_id}/onboarding/ownership_documents":                                               true,
	"GET /v1/instances/{instance_id}/rfi":                                                                          true,
	"POST /v1/instances/{instance_id}/rfi":                                                                         true,
	"PUT /v1/instances/{instance_id}/customers/{customer_id}/bank-accounts/{bank_account_id}/offramp-wallets/{id}": true,
	"POST /v1/presign":        true,
	"POST /v1/upload/analyze": true,
	"GET /v1/instances/{instance_id}/customers/{customer_id}/wallets/{id}/deposits": true,
}

// expectedPath fills an op's path template with the placeholder ids.
func expectedPath(template string) string {
	p := strings.Replace(template, "/instances/{id}", "/instances/"+tInstance, 1)
	return strings.NewReplacer(
		"{instance_id}", tInstance,
		"{customer_id}", tCustomer,
		"{bank_account_id}", tBankAccount,
		"{user_id}", tUser,
		"{swift}", tSwift,
		"{id}", tID,
	).Replace(p)
}

func readOpIndex(t *testing.T) []string {
	t.Helper()
	var index struct {
		Operations map[string]string `json:"operations"`
	}
	require.NoError(t, json.Unmarshal(blindpaytest.ReadFixture(t, "index.json"), &index))
	ops := make([]string, 0, len(index.Operations))
	for op := range index.Operations {
		ops = append(ops, op)
	}
	sort.Strings(ops)
	return ops
}

// TestOperationCoverage is the A8 op-coverage test: every x-sdk operation
// maps to SDK methods that send the right method and path and decode every
// variant of the operation's shared fixture.
func TestOperationCoverage(t *testing.T) {
	ops := readOpIndex(t)
	require.Len(t, ops, 94)

	inIndex := map[string]bool{}
	for _, op := range ops {
		inIndex[op] = true
		_, implemented := opTable[op]
		require.NotEqual(t, implemented, pendingOps[op], "%s: must be in exactly one of opTable and pendingOps", op)
	}
	for op := range opTable {
		require.True(t, inIndex[op], "opTable has %s, which is not an x-sdk operation", op)
	}

	for _, op := range ops {
		entry, ok := opTable[op]
		if !ok {
			continue
		}
		method, template, _ := strings.Cut(op, " ")
		path := expectedPath(template)
		variants := blindpaytest.ResponseVariants(t, op)

		for i, fn := range entry.Calls {
			i, fn := i, fn
			for _, v := range variants {
				v := v
				t.Run(fmt.Sprintf("%s/call%d/%s", op, i, v.Name), func(t *testing.T) {
					rt := &blindpaytest.RoundTripper{T: t, Method: method, Path: path, Out: v.Body}
					c, err := New("test_key", tInstance,
						WithBaseURL("https://api.blindpay.com/v1"),
						WithHTTPClient(&http.Client{Transport: rt}),
						WithMaxRetries(0),
					)
					require.NoError(t, err)
					require.NoError(t, fn(context.Background(), c))
				})
			}
		}
	}

	if out := os.Getenv("BLINDPAY_COVERAGE_OUT"); out != "" {
		writeCoverageTable(t, out)
	}
}

// writeCoverageTable writes the cross-SDK coverage table (plan §9 item 8),
// one line per op in sdk-surface.txt order, to path.
func writeCoverageTable(t *testing.T, path string) {
	t.Helper()
	surface := os.Getenv("BLINDPAY_SURFACE")
	require.NotEmpty(t, surface, "set BLINDPAY_SURFACE to sdk-surface.txt")
	f, err := os.Open(surface)
	require.NoError(t, err)
	defer f.Close()

	var b strings.Builder
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		head, _, _ := strings.Cut(sc.Text(), "|")
		fields := strings.Fields(head)
		op := fields[0] + " " + fields[1]
		entry, ok := opTable[op]
		if !ok {
			fmt.Fprintf(&b, "%s | -- | options:no | autopager:-- | MISSING\n", op)
			continue
		}
		pager, notes := "--", "--"
		if entry.Pager != "" {
			pager = entry.Pager
		}
		if entry.Notes != "" {
			notes = entry.Notes
		}
		fmt.Fprintf(&b, "%s | %s | options:yes | autopager:%s | %s\n", op, entry.Name, pager, notes)
	}
	require.NoError(t, sc.Err())
	require.NoError(t, os.WriteFile(path, []byte(b.String()), 0o644))
}
