# BlindPay Go SDK Examples

This directory contains example code showing how to use the BlindPay Go SDK.

## Prerequisites

1. Go 1.21 or higher installed
2. BlindPay API credentials (API key and Instance ID)

## Running Examples

### Quick Start - Get Available Rails

A simple example that shows how to get available rails:

```bash
cd examples/getrails
go run main.go
```

### Get Bank Details

A simple example that shows how to get bank details for a specific rail:

```bash
cd examples/getbankdetails
go run main.go
```

### Workflows

Each reads `BLINDPAY_API_KEY` and `BLINDPAY_INSTANCE_ID` from the environment
(use a development instance). They are compiled by `go build ./...` in CI but
never run there.

| Example | Shows |
|---|---|
| `webhooks` | Register a webhook endpoint, verify deliveries on the raw body, dedupe on `svix-id`, branch on typed events |
| `managedwalletpayout` | Quote to a bank account, pay out from a managed wallet (`WalletID`) with an idempotency key, handle a retryable 5xx |
| `fiatpayin` | Payin quote, then payin; print the PIX code to pay |
| `reconciliation` | Walk every completed payout of a customer with `ListAutoPaging` |

```bash
cd examples/reconciliation
BLINDPAY_API_KEY=... BLINDPAY_INSTANCE_ID=... go run .
```

- Check the main [README](../README.md) for SDK documentation
- Visit the [BlindPay API documentation](https://api.blindpay.com/reference)
- Open an issue on [GitHub](https://github.com/blindpaylabs/blindpay-go/issues)

