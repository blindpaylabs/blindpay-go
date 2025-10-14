package bankaccounts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/stretchr/testify/require"
)

func TestBankAccounts_CreatePix(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"
	id := "ba_123"
	name := "My PIX Account"
	pixKey := "test@example.com"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				In:     json.RawMessage(fmt.Sprintf(`{"type":"pix","name":"%s","pix_key":"%s"}`, name, pixKey)),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","type":"pix","name":"%s","pix_key":"%s","created_at":"2024-01-01T00:00:00Z"}`, id, name, pixKey)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreatePix(context.Background(), &CreatePixParams{
		ReceiverID: receiverID,
		Name:       name,
		PixKey:     pixKey,
	})
	require.NoError(t, err)
	require.Equal(t, id, account.ID)
	require.Equal(t, name, account.Name)
	require.Equal(t, pixKey, account.PixKey)
}

func TestBankAccounts_CreateAch(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"
	id := "ba_456"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"ach",
					"name":"My ACH Account",
					"account_class":"individual",
					"account_number":"123456789",
					"account_type":"checking",
					"beneficiary_name":"John Doe",
					"routing_number":"021000021"
				}`),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","type":"ach","name":"My ACH Account","created_at":"2024-01-01T00:00:00Z"}`, id)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreateAch(context.Background(), &CreateAchParams{
		ReceiverID:      receiverID,
		Name:            "My ACH Account",
		AccountClass:    types.AccountClassIndividual,
		AccountNumber:   "123456789",
		AccountType:     types.BankAccountTypeChecking,
		BeneficiaryName: "John Doe",
		RoutingNumber:   "021000021",
	})
	require.NoError(t, err)
	require.Equal(t, id, account.ID)
}

func TestBankAccounts_List(t *testing.T) {
	receiverID := "recv_123"
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"data":[
						{"id":"ba_123","type":"pix","name":"PIX Account","created_at":"2024-01-01T00:00:00Z"},
						{"id":"ba_456","type":"ach","name":"ACH Account","created_at":"2024-01-01T00:00:00Z"}
					]
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.List(context.Background(), receiverID)
	require.NoError(t, err)
	require.Len(t, response.Data, 2)
	require.Equal(t, "ba_123", response.Data[0].ID)
}

func TestBankAccounts_Get(t *testing.T) {
	receiverID := "recv_123"
	id := "ba_123"
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","receiver_id":"%s","account_holder_name":"John Doe","created_at":"2024-01-01T00:00:00Z","updated_at":"2024-01-01T00:00:00Z"}`, id, receiverID)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s", instanceID, receiverID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.Get(context.Background(), receiverID, id)
	require.NoError(t, err)
	require.Equal(t, id, account.ID)
	require.Equal(t, receiverID, account.ReceiverID)
}

func TestBankAccounts_Delete(t *testing.T) {
	receiverID := "recv_123"
	id := "ba_123"
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s", instanceID, receiverID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background(), receiverID, id)
	require.NoError(t, err)
}
