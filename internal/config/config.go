package config

import (
	"net/http"

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
}

/*
Converts the top-level Config into the transport/request layer config used by internal HTTP helpers.
This reduces coupling between the transport layer and higher-level SDK components.
*/
func (c *Config) ToRequestConfig() *request.Config {
	if c == nil {
		return nil
	}

	return &request.Config{
		BaseURL:    c.BaseURL,
		APIKey:     c.APIKey,
		HTTPClient: c.HTTPClient,
		UserAgent:  c.UserAgent,
	}
}
