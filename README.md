# BlindPay Go SDK

The official Go SDK for [BlindPay](https://blindpay.com) - Global payments infrastructure made simple.

## 🚀 Installation

```bash
go get github.com/blindpaylabs/blindpay-go
```

## 📋 Requirements

- Go 1.21 or higher
- BlindPay API credentials (API Key and Instance ID)

## 🔧 Quick Start

### Basic Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/blindpaylabs/blindpay-go"
)

func main() {
	// Create a new BlindPay client
	client, err := blindpay.New(
		"your-api-key-here",
		"your-instance-id-here",
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get available payment rails
	rails, err := client.Available.GetRails(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d payment rails\n", len(rails))
	for _, rail := range rails {
		fmt.Printf("- %s (%s) in %s\n", rail.Label, rail.Value, rail.Country)
	}
}
```

### Running the Example

```bash
go run main.go
```

For detailed API documentation, visit:
- [Blindpay API documentation](https://blindpay.com/docs/getting-started/overview)
- [API Reference](https://api.blindpay.com/reference)

## Support

- 📧 Email: [gabriel@blindpay.com](mailto:gabriel@blindpay.com)
- 🐛 Issues: [GitHub Issues](https://github.com/blindpaylabs/blindpay-go/issues)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Made with ❤️ by the [Blindpay](https://blindpay.com) team