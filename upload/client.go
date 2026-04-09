package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/blindpaylabs/blindpay-go/internal/config"
)

// Bucket represents the upload destination bucket.
type Bucket string

const (
	BucketAvatar        Bucket = "avatar"
	BucketOnboarding    Bucket = "onboarding"
	BucketLimitIncrease Bucket = "limit_increase"
)

// UploadParams represents parameters for uploading a file.
type UploadParams struct {
	File       io.Reader
	FileName   string
	Bucket     Bucket
	InstanceID string
}

// UploadResponse represents the response when uploading a file.
type UploadResponse struct {
	FileURL string `json:"file_url"`
}

// Client handles file upload operations.
type Client struct {
	baseURL    string
	apiKey     string
	instanceID string
	httpClient *http.Client
	userAgent  string
}

// NewClient creates a new upload client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		baseURL:    cfg.BaseURL,
		apiKey:     cfg.APIKey,
		instanceID: cfg.InstanceID,
		httpClient: cfg.HTTPClient,
		userAgent:  cfg.UserAgent,
	}
}

// Upload uploads a file.
func (c *Client) Upload(ctx context.Context, params *UploadParams) (*UploadResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.File == nil {
		return nil, fmt.Errorf("file cannot be nil")
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	part, err := writer.CreateFormFile("file", params.FileName)
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, params.File); err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	if err := writer.WriteField("bucket", string(params.Bucket)); err != nil {
		return nil, fmt.Errorf("failed to write bucket field: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close writer: %w", err)
	}

	instanceID := params.InstanceID
	if instanceID == "" {
		instanceID = c.instanceID
	}

	url := fmt.Sprintf("%s/upload?instance_id=%s", c.baseURL, instanceID)

	req, err := http.NewRequestWithContext(ctx, "POST", url, &body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("upload failed with status %d: %s", resp.StatusCode, string(respBody))
	}

	var result UploadResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
