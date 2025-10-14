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
	instanceID := "in_000000000000"
	id := "fe_000000000000"
	name := "Display Name"
	evmWallet := "0x1234567890123456789012345678901234567890"
	stellarWallet := "GAB22222222222222222222222222222222222222222222222222222222222222"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"name":"Display Name",
					"payout_percentage_fee":0,
					"payout_flat_fee":0,
					"payin_percentage_fee":0,
					"payin_flat_fee":0,
					"evm_wallet_address":"0x1234567890123456789012345678901234567890",
					"stellar_wallet_address":"GAB22222222222222222222222222222222222222222222222222222222222222"
				}`),
				Out: json.RawMessage(fmt.Sprintf(
					`{
						"id":"%s",
						"instance_id":"%s",
						"name":"%s",
						"payout_percentage_fee":0,
						"payout_flat_fee":0,
						"payin_percentage_fee":0,
						"payin_flat_fee":0,
						"evm_wallet_address":"%s",
						"stellar_wallet_address":"%s"
					}`, id, instanceID, name, evmWallet, stellarWallet,
				)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/partner-fees", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	fee, err := client.Create(context.Background(), &CreatePartnerFeeParams{
		Name:                 name,
		PayoutPercentageFee:  0,
		PayoutFlatFee:        0,
		PayinPercentageFee:   0,
		PayinFlatFee:         0,
		EVMWalletAddress:     evmWallet,
		StellarWalletAddress: &stellarWallet,
	})
	require.NoError(t, err)
	require.Equal(t, id, fee.ID)
	require.Equal(t, name, fee.Name)
	require.Equal(t, 0.0, fee.PayoutPercentageFee)
}

func TestPartnerFees_List(t *testing.T) {
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"id":"fe_000000000000",
						"instance_id":"in_000000000000",
						"name":"Display Name",
						"payout_percentage_fee":0,
						"payout_flat_fee":0,
						"payin_percentage_fee":0,
						"payin_flat_fee":0,
						"evm_wallet_address":"0x1234567890123456789012345678901234567890",
						"stellar_wallet_address":"GAB22222222222222222222222222222222222222222222222222222222222222"
					}
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
	require.Len(t, fees, 1)
	require.Equal(t, "fe_000000000000", fees[0].ID)
	require.Equal(t, "Display Name", fees[0].Name)
}

func TestPartnerFees_Get(t *testing.T) {
	instanceID := "in_000000000000"
	id := "fe_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(
					`{
						"id":"fe_000000000000",
						"instance_id":"in_000000000000",
						"name":"Display Name",
						"payout_percentage_fee":0,
						"payout_flat_fee":0,
						"payin_percentage_fee":0,
						"payin_flat_fee":0,
						"evm_wallet_address":"0x1234567890123456789012345678901234567890",
						"stellar_wallet_address":"GAB22222222222222222222222222222222222222222222222222222222222222"
					}`),
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
	require.Equal(t, "Display Name", fee.Name)
}

func TestPartnerFees_Delete(t *testing.T) {
	instanceID := "in_000000000000"
	id := "fe_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
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
