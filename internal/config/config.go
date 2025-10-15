package config

import "net/http"

type Config struct {
	BaseURL    string
	APIKey     string
	InstanceID string
	HTTPClient *http.Client
	UserAgent  string
}
