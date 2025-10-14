package partnerfees

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/blindpaytest"
	"github.com/blindpaylabs/blindpay-go/internal/config"
	"github.com/stretchr/testify/require"
)

func TestPartnerFees_Create(t *testing.T) {
	instanceID := "inst_123"
	id := "pf_123"
	name := "Test Fee Config"
	evmWallet := "0x1234567890abcdef"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"evm_wallet_address":"0x1234567890abcdef",
					"name":"Test Fee Config",
					"payin_flat_fee":1.5,
					"payin_percentage_fee":0.02,
					"payout_flat_fee":2.0,
					"payout_percentage_fee":0.025
				}`),
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","instance_id":"%s","name":"%s","payout_percentage_fee":0.025,"payout_flat_fee":2.0,"payin_percentage_fee":0.02,"payin_flat_fee":1.5,"evm_wallet_address":"%s"}`, id, instanceID, name, evmWallet)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/partner-fees", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	fee, err := client.Create(context.Background(), &CreatePartnerFeeParams{
		Name:                name,
		EVMWalletAddress:    evmWallet,
		PayinFlatFee:        1.5,
		PayinPercentageFee:  0.02,
		PayoutFlatFee:       2.0,
		PayoutPercentageFee: 0.025,
	})
	require.NoError(t, err)
	require.Equal(t, id, fee.ID)
	require.Equal(t, name, fee.Name)
	require.Equal(t, 0.025, fee.PayoutPercentageFee)
}

func TestPartnerFees_List(t *testing.T) {
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{"id":"pf_123","instance_id":"inst_123","name":"Fee Config 1","payout_percentage_fee":0.025,"payout_flat_fee":2.0,"payin_percentage_fee":0.02,"payin_flat_fee":1.5},
					{"id":"pf_456","instance_id":"inst_123","name":"Fee Config 2","payout_percentage_fee":0.03,"payout_flat_fee":2.5,"payin_percentage_fee":0.025,"payin_flat_fee":2.0}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/partner-fees", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	fees, err := client.List(context.Background())
	require.NoError(t, err)
	require.Len(t, fees, 2)
	require.Equal(t, "pf_123", fees[0].ID)
	require.Equal(t, "Fee Config 1", fees[0].Name)
}

func TestPartnerFees_Get(t *testing.T) {
	instanceID := "inst_123"
	id := "pf_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s","instance_id":"%s","name":"Fee Config","payout_percentage_fee":0.025,"payout_flat_fee":2.0,"payin_percentage_fee":0.02,"payin_flat_fee":1.5}`, id, instanceID)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/partner-fees/%s", instanceID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	fee, err := client.Get(context.Background(), id)
	require.NoError(t, err)
	require.Equal(t, id, fee.ID)
	require.Equal(t, "Fee Config", fee.Name)
}

func TestPartnerFees_Delete(t *testing.T) {
	instanceID := "inst_123"
	id := "pf_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s/partner-fees/%s", instanceID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background(), id)
	require.NoError(t, err)
}
