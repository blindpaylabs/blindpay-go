// Command fiatpayin creates a payin quote and a payin, then prints the
// instructions the customer pays (docs: payins).
package main

import (
	"context"
	"log"
	"os"

	"github.com/blindpaylabs/blindpay-go"
	"github.com/blindpaylabs/blindpay-go/payins"
)

func main() {
	client, err := blindpay.New(os.Getenv("BLINDPAY_API_KEY"), os.Getenv("BLINDPAY_INSTANCE_ID"))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	quote, err := client.Payins.Quotes.Create(ctx, &payins.CreateQuoteParams{
		BlockchainWalletID: "bw_000000000000",
		CurrencyType:       blindpay.CurrencyTypeSender,
		PaymentMethod:      "pix",
		RequestAmount:      100000, // integer cents: 1,000.00 BRL
		Token:              "USDC",
	})
	if err != nil {
		log.Fatal(err)
	}

	payin, err := client.Payins.CreateEvm(ctx, quote.ID, blindpay.WithIdempotencyKey("payin-"+quote.ID))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("payin %s: ask the customer to pay PIX code %q", payin.ID, payin.PixCode)
}
