package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blindpaylabs/blindpay-go"
)

func getAvailableRails() error {
	// Create a new BlindPay client
	client, err := blindpay.New(
		"your-api-key-here",
		"your-instance-id-here",
	)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Get available rails
	rails, err := client.Available.GetRails(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get rails: %w", err)
	}

	fmt.Println("Rails:", rails)
	return nil
}

func main() {
	if err := getAvailableRails(); err != nil {
		log.Fatal(err)
	}
}
