// Command managedwalletpayout funds a BlindPay managed wallet with a payin,
// then pays out from it to a bank account using the wallet's ID (WalletID).
// Run it against a development instance: settlement uses the test token
// USDB, a payin completes after about 30 seconds, and a payout quote of
// 66600 cents fails while 77700 is refunded.
package main

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"github.com/blindpaylabs/blindpay-go"
	"github.com/blindpaylabs/blindpay-go/payins"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/blindpaylabs/blindpay-go/quotes"
)

const (
	walletID      = "wa_000000000000" // a managed wallet (client.CustodialWallets)
	bankAccountID = "ba_000000000000" // the customer's bank account
)

func main() {
	client, err := blindpay.New(os.Getenv("BLINDPAY_API_KEY"), os.Getenv("BLINDPAY_INSTANCE_ID"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// 1. Payin into the managed wallet: quote with WalletID, then execute.
	wallet := walletID
	payinQuote, err := client.Payins.Quotes.Create(ctx, &payins.CreateQuoteParams{
		WalletID:      &wallet,
		CurrencyType:  blindpay.CurrencyTypeSender,
		PaymentMethod: "pix",
		RequestAmount: 100000, // integer cents: 1,000.00 BRL
		Token:         blindpay.StablecoinUSDB,
	})
	if err != nil {
		log.Fatal(err)
	}
	payin, err := client.Payins.CreateEvm(ctx, payinQuote.ID, blindpay.WithIdempotencyKey("payin-"+payinQuote.ID))
	if err != nil {
		log.Fatal(err)
	}

	// 2. Wait for the payin to complete. In production, react to the
	// payin.complete (or wallet.inbound) webhook instead of polling.
	for payin.Status != blindpay.TransactionStatusCompleted {
		if payin.Status == blindpay.TransactionStatusFailed || payin.Status == blindpay.TransactionStatusRefunded {
			log.Fatalf("payin %s ended %s", payin.ID, payin.Status)
		}
		time.Sleep(10 * time.Second) // development instances are rate limited
		p, err := client.Payins.Get(ctx, payin.ID)
		if err != nil {
			log.Fatal(err)
		}
		payin.Status = p.Status
	}

	// 3. Quote the payout to the bank account.
	network, token := blindpay.NetworkBase, blindpay.StablecoinUSDB
	quote, err := client.Quotes.Create(ctx, &quotes.CreateParams{
		BankAccountID: bankAccountID,
		CurrencyType:  blindpay.CurrencyTypeSender,
		Network:       &network, // the managed wallet's network
		Token:         &token,
		RequestAmount: 10000, // integer cents: 100.00 USDB
	})
	if err != nil {
		log.Fatal(err)
	}

	// 4. Pay from the managed wallet (WalletID), not a customer wallet
	// (SenderWalletAddress): set exactly one. The idempotency key makes a
	// resend safe.
	payout, err := client.Payouts.CreateEvm(ctx, &payouts.CreateEvmParams{
		QuoteID:  quote.ID,
		WalletID: walletID,
	}, blindpay.WithIdempotencyKey("payout-"+quote.ID))

	var apiErr *blindpay.APIError
	switch {
	case errors.As(err, &apiErr) && apiErr.Retryable:
		// A 5xx is not retried automatically: check whether the payout exists
		// (e.g. list payouts) before resending with apiErr.IdempotencyKey.
		log.Fatalf("payout outcome unknown, key %s: %v", apiErr.IdempotencyKey, err)
	case errors.As(err, &apiErr) && apiErr.RetryClass() == blindpay.RetryClassActionable:
		log.Fatalf("fix the request first (%s): %s", apiErr.Code, apiErr.Description)
	case err != nil:
		log.Fatal(err)
	}
	log.Printf("payout %s is %s; follow it with the payout.update and payout.complete webhooks", payout.ID, payout.Status)
}
