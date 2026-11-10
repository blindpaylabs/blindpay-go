package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
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
	cfg        *request.Config
	instanceID string
}

// NewClient creates a new upload client.
func NewClient(cfg *config.Config) *Client {
	return &Client{
		cfg:        cfg.ToRequestConfig(),
		instanceID: cfg.InstanceID,
	}
}

// Upload uploads a file.
//
// A non-2xx response is returned as *blindpay.APIError (same shape as every
// other method).
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

	path := "/upload?" + url.Values{"instance_id": {instanceID}}.Encode()
	return request.DoMultipart[*UploadResponse](c.cfg, ctx, path, body.Bytes(), writer.FormDataContentType())
}
