package config

import (
	"net/http"
	"time"

	"github.com/blindpaylabs/blindpay-go/internal/request"
)

/*
Config holds the top-level configuration required to initialize Blindpay SDK clients.
It may include both transport-level fields and the instance identifier used to compose API paths.
*/
type Config struct {
	BaseURL    string
	APIKey     string
	InstanceID string
	HTTPClient *http.Client
	UserAgent  string
	// Timeout bounds each attempt; zero means none.
	Timeout time.Duration
	// MaxRetries is the number of retries after the first attempt; zero
	// disables retries and automatic idempotency keys.
	MaxRetries int
	// MaxResponseBodyBytes caps how much of a response body is read; zero
	// means the default (32 MiB).
	MaxResponseBodyBytes int64
}

/*
Converts the top-level Config into the transport/request layer config used by internal HTTP helpers.
This reduces coupling between the transport layer and higher-level SDK components.
*/
func (c *Config) ToRequestConfig() *request.Config {
	if c == nil {
		return nil
	}
	apiKey := c.APIKey

	return &request.Config{
		BaseURL:    c.BaseURL,
		APIKey:     &apiKey,
		HTTPClient: c.HTTPClient,
		UserAgent:  c.UserAgent,
		Timeout:    c.Timeout,
		MaxRetries: c.MaxRetries,

		MaxResponseBodyBytes: c.MaxResponseBodyBytes,
	}
}
