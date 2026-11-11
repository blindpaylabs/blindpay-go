package upload

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

type rtFunc func(*http.Request) (*http.Response, error)

func (f rtFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type capturedPart struct {
	fileName, contentType, value string
}

// multipartServer records the multipart parts of the request and replies
// with body.
func multipartServer(t *testing.T, wantPath string, parts map[string]capturedPart, body []byte) *Client {
	return NewClient(&config.Config{
		BaseURL:    "https://api.blindpay.com/v1",
		APIKey:     "test_key",
		InstanceID: "in_000000000000",
		UserAgent:  "blindpay-go/test",
		HTTPClient: &http.Client{Transport: rtFunc(func(r *http.Request) (*http.Response, error) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, wantPath, r.URL.Path)
			require.Equal(t, url.Values{"instance_id": {"in_000000000000"}}, r.URL.Query())
			require.Empty(t, r.Header.Get("Idempotency-Key"), "multipart requests carry no idempotency key")

			mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
			require.NoError(t, err)
			require.Equal(t, "multipart/form-data", mediaType)
			mr := multipart.NewReader(r.Body, params["boundary"])
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					break
				}
				require.NoError(t, err)
				b, _ := io.ReadAll(p)
				parts[p.FormName()] = capturedPart{fileName: p.FileName(), contentType: p.Header.Get("Content-Type"), value: string(b)}
			}
			return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}}, Body: io.NopCloser(bytes.NewReader(body)), Request: r}, nil
		})},
	})
}

func TestUpload_Analyze(t *testing.T) {
	v := blindpaytest.FindVariant(t, "POST /v1/upload/analyze", "full")
	parts := map[string]capturedPart{}
	client := multipartServer(t, "/v1/upload/analyze", parts, v.Body)

	res, err := client.Analyze(context.Background(), &AnalyzeParams{
		File:        strings.NewReader("%PDF"),
		FileName:    "statement.pdf",
		ContentType: "application/pdf",
		Type:        AnalyzeTypeBankStatement,
		Metadata:    `{"name":"Jane"}`,
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.ApprovalRate)

	require.Equal(t, capturedPart{fileName: "statement.pdf", contentType: "application/pdf", value: "%PDF"}, parts["file"])
	require.Equal(t, "bank_statement", parts["type"].value)
	require.Equal(t, `{"name":"Jane"}`, parts["metadata"].value)
}

func TestUpload_UploadMultipartParts(t *testing.T) {
	parts := map[string]capturedPart{}
	client := multipartServer(t, "/v1/upload", parts, []byte(`{"file_url":"https://x/1.png"}`))

	_, err := client.Upload(context.Background(), &UploadParams{File: strings.NewReader("png"), FileName: `a"b.png`, Bucket: BucketDocuments})
	require.NoError(t, err)
	require.Equal(t, capturedPart{fileName: `a"b.png`, contentType: "application/octet-stream", value: "png"}, parts["file"])
	require.Equal(t, "documents", parts["bucket"].value)
	_, hasMetadata := parts["metadata"]
	require.False(t, hasMetadata)
}

func TestUpload_Presign(t *testing.T) {
	v := blindpaytest.FindVariant(t, "POST /v1/presign", "full")
	client := newTestClient(&blindpaytest.RoundTripper{
		T:      t,
		Method: http.MethodPost,
		Path:   "/v1/presign",
		Query:  url.Values{"instance_id": {"in_000000000000"}},
		In:     json.RawMessage(`{"file_url":"https://x/1.pdf"}`),
		Out:    v.Body,
	})

	res, err := client.Presign(context.Background(), &PresignParams{FileURL: "https://x/1.pdf"})
	require.NoError(t, err)
	require.NotEmpty(t, res.FileURL)
	require.NotEmpty(t, res.ExpiresAt)

	_, err = client.Presign(context.Background(), &PresignParams{})
	require.Error(t, err)
}
