// Command reconciliation walks every completed payout of a customer with the
// auto-pager (docs: pagination).
package main

import (
	"context"
	"log"
	"os"

	"github.com/blindpaylabs/blindpay-go"
	"github.com/blindpaylabs/blindpay-go/payouts"
)

func main() {
	client, err := blindpay.New(os.Getenv("BLINDPAY_API_KEY"), os.Getenv("BLINDPAY_INSTANCE_ID"))
	if err != nil {
		log.Fatal(err)
	}

	pager := client.Payouts.ListAutoPaging(context.Background(), &payouts.ListParams{
		CustomerID: "cu_000000000000",
		Status:     blindpay.TransactionStatusCompleted,
	})
	var total float64
	for pager.Next() {
		p := pager.Current()
		total += p.SenderAmount
	}
	if err := pager.Err(); err != nil {
		log.Fatal(err)
	}
	log.Printf("completed payouts total %.2f (cents)", total)
}
