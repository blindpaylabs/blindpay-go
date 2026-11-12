// Package upload uploads documents (multipart, never retried), analyzes a
// document before you submit it, and presigns a temporary URL for an
// uploaded file. Use it through blindpay.Client.Upload.
package upload

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"strings"

	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/option"
)

// Bucket represents the upload destination bucket.
type Bucket string

const (
	BucketAvatar        Bucket = "avatar"
	BucketOnboarding    Bucket = "onboarding"
	BucketLimitIncrease Bucket = "limit_increase"
	BucketDocuments     Bucket = "documents"
)

// UploadParams represents parameters for uploading a file.
type UploadParams struct {
	File     io.Reader
	FileName string
	// ContentType is the file's MIME type (e.g. "application/pdf"). Empty
	// sends application/octet-stream.
	ContentType string
	Bucket      Bucket
	// InstanceID overrides the client's instance. The API authenticates the
	// key against it, so leave it empty.
	InstanceID string
}

// AnalyzeType is the kind of document Analyze checks.
type AnalyzeType string

const (
	AnalyzeTypeProofOfAddress        AnalyzeType = "proof_of_address"
	AnalyzeTypeIdentityDocument      AnalyzeType = "identity_document"
	AnalyzeTypeIncorporationDocument AnalyzeType = "incorporation_document"
	AnalyzeTypeProofOfOwnership      AnalyzeType = "proof_of_ownership"
	AnalyzeTypeSourceOfFunds         AnalyzeType = "source_of_funds"
	AnalyzeTypeSelfie                AnalyzeType = "selfie"
	AnalyzeTypeBankStatement         AnalyzeType = "bank_statement"
	AnalyzeTypeTaxReturn             AnalyzeType = "tax_return"
	AnalyzeTypeProofOfIncome         AnalyzeType = "proof_of_income"
	AnalyzeTypeFinancialStatement    AnalyzeType = "financial_statement"
	AnalyzeTypeInvoice               AnalyzeType = "invoice"
	AnalyzeTypeTransactionDocument   AnalyzeType = "transaction_document"
	AnalyzeTypePassport              AnalyzeType = "passport"
	AnalyzeTypePayStub               AnalyzeType = "pay_stub"
	AnalyzeTypeEmploymentLetter      AnalyzeType = "employment_letter"
	AnalyzeTypeInvestment            AnalyzeType = "investment"
	AnalyzeTypeCryptoExchange        AnalyzeType = "crypto_exchange"
	AnalyzeTypeBlockchainWallet      AnalyzeType = "blockchain_wallet"
	AnalyzeTypeContract              AnalyzeType = "contract"
	AnalyzeTypeAccountsReceivable    AnalyzeType = "accounts_receivable"
	AnalyzeTypeMerchantProcessor     AnalyzeType = "merchant_processor"
	AnalyzeTypeShareholderLoan       AnalyzeType = "shareholder_loan"
)

// AnalyzeParams are the inputs of Analyze.
type AnalyzeParams struct {
	File        io.Reader
	FileName    string
	ContentType string
	Type        AnalyzeType
	// Metadata is optional context for the analysis, as a JSON string
	// (e.g. `{"name":"Jane Doe"}`).
	Metadata string
}

// ApprovalRate is how likely the analyzed document is to be accepted.
type ApprovalRate string

const (
	ApprovalRateLow    ApprovalRate = "low"
	ApprovalRateMedium ApprovalRate = "medium"
	ApprovalRateHigh   ApprovalRate = "high"
)

// AnalyzeResponse is the result of Analyze.
type AnalyzeResponse struct {
	ApprovalRate ApprovalRate `json:"approval_rate"`
	Description  string       `json:"description"`
}

// PresignParams are the inputs of Presign.
type PresignParams struct {
	// FileURL is a file URL returned by Upload.
	FileURL string `json:"file_url"`
}

// PresignResponse is a temporary, signed URL for a private uploaded file.
type PresignResponse struct {
	FileURL   string `json:"file_url"`
	ExpiresAt string `json:"expires_at"`
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

// Upload uploads a file (multipart/form-data) and returns its URL.
//
// Uploads are never retried and carry no idempotency key: the API does not
// protect multipart requests. A non-2xx response is returned as
// *blindpay.APIError, like every other method.
func (c *Client) Upload(ctx context.Context, params *UploadParams, opts ...option.RequestOption) (*UploadResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.File == nil {
		return nil, fmt.Errorf("file cannot be nil")
	}

	body, contentType, err := multipartBody(params.File, params.FileName, params.ContentType, [][2]string{{"bucket", string(params.Bucket)}})
	if err != nil {
		return nil, err
	}

	instanceID := params.InstanceID
	if instanceID == "" {
		instanceID = c.instanceID
	}

	var q request.Query
	q.String("instance_id", instanceID)
	path := q.Apply("/upload")
	return request.DoMultipart[*UploadResponse](c.cfg, ctx, path, body, contentType, opts...)
}

// Analyze checks a document before you submit it (e.g. for KYC) and returns
// how likely it is to be approved. Like Upload, it is multipart and never
// retried.
func (c *Client) Analyze(ctx context.Context, params *AnalyzeParams, opts ...option.RequestOption) (*AnalyzeResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.File == nil {
		return nil, fmt.Errorf("file cannot be nil")
	}

	fields := [][2]string{{"type", string(params.Type)}}
	if params.Metadata != "" {
		fields = append(fields, [2]string{"metadata", params.Metadata})
	}
	body, contentType, err := multipartBody(params.File, params.FileName, params.ContentType, fields)
	if err != nil {
		return nil, err
	}

	var q request.Query
	q.String("instance_id", c.instanceID)
	path := q.Apply("/upload/analyze")
	return request.DoMultipart[*AnalyzeResponse](c.cfg, ctx, path, body, contentType, opts...)
}

// Presign returns a temporary signed URL for a file uploaded with Upload.
func (c *Client) Presign(ctx context.Context, params *PresignParams, opts ...option.RequestOption) (*PresignResponse, error) {
	if params == nil {
		return nil, fmt.Errorf("params cannot be nil")
	}
	if params.FileURL == "" {
		return nil, fmt.Errorf("file URL cannot be empty")
	}

	var q request.Query
	q.String("instance_id", c.instanceID)
	path := q.Apply("/presign")
	return request.Do[*PresignResponse](c.cfg, ctx, "POST", path, params, opts...)
}

// multipartBody encodes a "file" part followed by the given text fields.
func multipartBody(file io.Reader, fileName, fileType string, fields [][2]string) ([]byte, string, error) {
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	if fileType == "" {
		fileType = "application/octet-stream"
	}
	header := textproto.MIMEHeader{}
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="file"; filename="%s"`, escapeQuotes(fileName)))
	header.Set("Content-Type", fileType)
	part, err := writer.CreatePart(header)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create form file: %w", err)
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, "", fmt.Errorf("failed to copy file: %w", err)
	}

	for _, f := range fields {
		if err := writer.WriteField(f[0], f[1]); err != nil {
			return nil, "", fmt.Errorf("failed to write %s field: %w", f[0], err)
		}
	}
	if err := writer.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close writer: %w", err)
	}
	return body.Bytes(), writer.FormDataContentType(), nil
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

// escapeQuotes matches mime/multipart's (unexported) filename escaping.
func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
