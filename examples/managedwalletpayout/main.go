// Command managedwalletpayout quotes a payout to a bank account and pays it
// from a BlindPay managed wallet (docs: migration guide, managed wallets).
package main

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/blindpaylabs/blindpay-go"
	"github.com/blindpaylabs/blindpay-go/payouts"
	"github.com/blindpaylabs/blindpay-go/quotes"
)

func main() {
	client, err := blindpay.New(os.Getenv("BLINDPAY_API_KEY"), os.Getenv("BLINDPAY_INSTANCE_ID"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	network, token := blindpay.NetworkBase, blindpay.StablecoinToken("USDC")
	quote, err := client.Quotes.Create(ctx, &quotes.CreateParams{
		BankAccountID: "ba_000000000000",
		CurrencyType:  blindpay.CurrencyTypeSender,
		Network:       &network,
		Token:         &token,
		RequestAmount: 10000, // integer cents: 100.00 USDC
	})
	if err != nil {
		log.Fatal(err)
	}

	// Pay from the managed wallet (WalletID), not a customer wallet
	// (SenderWalletAddress): set exactly one. The idempotency key makes a
	// resend safe.
	payout, err := client.Payouts.CreateEvm(ctx, &payouts.CreateEvmParams{
		QuoteID:  quote.ID,
		WalletID: "wa_000000000000",
	}, blindpay.WithIdempotencyKey("payout-"+quote.ID))

	var apiErr *blindpay.APIError
	switch {
	case errors.As(err, &apiErr) && apiErr.Retryable:
		// A 5xx is not retried automatically: check whether the payout exists
		// (e.g. list payouts) before resending with apiErr.IdempotencyKey.
		log.Fatalf("payout outcome unknown, key %s: %v", apiErr.IdempotencyKey, err)
	case err != nil:
		log.Fatal(err)
	}
	log.Printf("payout %s is %s; follow it with the payout.update and payout.complete webhooks", payout.ID, payout.Status)
}
