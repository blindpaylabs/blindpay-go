package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/blindpaylabs/blindpay-go"
)

func main() {
	// Get credentials from environment variables (recommended)
	apiKey := os.Getenv("BLINDPAY_API_KEY")
	instanceID := os.Getenv("BLINDPAY_INSTANCE_ID")

	// Or use hardcoded values for testing (not recommended for production)
	if apiKey == "" {
		apiKey = "your-api-key-here"
	}
	if instanceID == "" {
		instanceID = "your-instance-id-here"
	}

	// Create a new BlindPay client
	client, err := blindpay.New(apiKey, instanceID)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	// Example 1: Get available rails
	fmt.Println("=== Example 1: Get Available Rails ===")
	rails, err := client.Available.GetRails(ctx)
	if err != nil {
		log.Fatalf("Failed to get rails: %v", err)
	}
	fmt.Printf("Found %d rails\n", len(rails))
	for _, rail := range rails {
		fmt.Printf("- %s (%s) in %s\n", rail.Label, rail.Value, rail.Country)
	}
	fmt.Println()

	// Example 2: Get bank details for a specific rail
	fmt.Println("=== Example 2: Get Bank Details ===")
	bankDetails, err := client.Available.GetBankDetails(ctx, "pix")
	if err != nil {
		log.Fatalf("Failed to get bank details: %v", err)
	}
	fmt.Printf("Bank details for PIX:\n")
	for _, detail := range bankDetails {
		fmt.Printf("- %s (key: %s, required: %v)\n", detail.Label, detail.Key, detail.Required)
	}
	fmt.Println()

	fmt.Println("✅ All examples completed successfully!")
}
