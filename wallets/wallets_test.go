package wallets

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

func TestWallets_CreateWithAddress(t *testing.T) {
	receiverID := "re_000000000000"
	instanceID := "in_000000000000"
	id := "bw_000000000000"
	address := "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
			"name":"Wallet Display Name",
			"network":"polygon",
			"address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
			"is_account_abstraction":true
		}`),
				Out: json.RawMessage(`{
					"id":"bw_000000000000",
					"name":"Wallet Display Name",
					"network":"polygon",
					"address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
					"signature_tx_hash":null,
					"is_account_abstraction":true,
					"receiver_id":"re_000000000000"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallet, err := client.CreateWithAddress(context.Background(), &CreateWithAddressParams{
		ReceiverID: receiverID,
		Name:       "Wallet Display Name",
		Network:    types.NetworkPolygon,
		Address:    address,
	})
	require.NoError(t, err)
	require.Equal(t, id, wallet.ID)
	require.Equal(t, address, wallet.Address)
	require.True(t, wallet.IsAccountAbstraction)
}

func TestWallets_List(t *testing.T) {
	receiverID := "re_000000000000"
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
						"id":"bw_000000000000",
						"name":"Wallet Display Name",
						"network":"polygon",
						"address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
						"signature_tx_hash":"0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
						"is_account_abstraction":false,
						"receiver_id":"re_000000000000"
					}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallets, err := client.List(context.Background(), receiverID)
	require.NoError(t, err)
	require.Len(t, wallets, 1)
	require.Equal(t, "bw_000000000000", wallets[0].ID)
	require.Equal(t, "Wallet Display Name", wallets[0].Name)
	require.Equal(t, types.NetworkPolygon, wallets[0].Network)
	require.Equal(t, "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C", wallets[0].Address)
	require.Equal(t, "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359", wallets[0].SignatureTxHash)
	require.False(t, wallets[0].IsAccountAbstraction)
	require.Equal(t, "re_000000000000", wallets[0].ReceiverID)
}

func TestWallets_Get(t *testing.T) {
	receiverID := "re_000000000000"
	id := "bw_000000000000"
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"bw_000000000000",
					"name":"Wallet Display Name",
					"network":"polygon",
					"address":"0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C",
					"signature_tx_hash":"0x3c499c542cef5e3811e1192ce70d8cc03d5c3359",
					"is_account_abstraction":false,
					"receiver_id":"re_000000000000"
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/%s", instanceID, receiverID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	wallet, err := client.Get(context.Background(), receiverID, id)
	require.NoError(t, err)
	require.Equal(t, id, wallet.ID)
	require.Equal(t, "Wallet Display Name", wallet.Name)
	require.Equal(t, types.NetworkPolygon, wallet.Network)
	require.Equal(t, "0xDD6a3aD0949396e57C7738ba8FC1A46A5a1C372C", wallet.Address)
	require.Equal(t, "0x3c499c542cef5e3811e1192ce70d8cc03d5c3359", wallet.SignatureTxHash)
	require.False(t, wallet.IsAccountAbstraction)
	require.Equal(t, "re_000000000000", wallet.ReceiverID)
}

func TestWallets_Delete(t *testing.T) {
	receiverID := "re_000000000000"
	id := "bw_000000000000"
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/%s", instanceID, receiverID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background(), receiverID, id)
	require.NoError(t, err)
}

func TestWallets_GetWalletMessage(t *testing.T) {
	receiverID := "re_000000000000"
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"message":"random"}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/blockchain-wallets/sign-message", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.GetWalletMessage(context.Background(), receiverID)
	require.NoError(t, err)
	require.Equal(t, "random", response.Message)
}

func TestWallets_CreateAssetTrustline(t *testing.T) {
	instanceID := "in_000000000000"
	address := "GCDNJUBQSX7AJWLJACMJ7I4BC3Z47BQUTMHEICZLE6MU4KQBRYG5JY6B"
	xdr := "AAAAAgAAAABqVFqpZzXx+KxRjXXFGO3sKwHCEYdHsWxDRrJTLGPDowAAAGQABVECAAAAAQAAAAEAAAAAAAAAAAAAAABmWFbUAAAAAAAAAAEAAAAAAAAABgAAAAFVU0RCAAAAAABbjPEfrLNLCLjNQyaWWgTeFn4tnbFnNd9FTJ3HgkLUCwAAAAAAAAAAAAAAAAAAAAE="

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"address":"GCDNJUBQSX7AJWLJACMJ7I4BC3Z47BQUTMHEICZLE6MU4KQBRYG5JY6B"
				}`),
				Out: json.RawMessage(`{
					"xdr":"AAAAAgAAAABqVFqpZzXx+KxRjXXFGO3sKwHCEYdHsWxDRrJTLGPDowAAAGQABVECAAAAAQAAAAEAAAAAAAAAAAAAAABmWFbUAAAAAAAAAAEAAAAAAAAABgAAAAFVU0RCAAAAAABbjPEfrLNLCLjNQyaWWgTeFn4tnbFnNd9FTJ3HgkLUCwAAAAAAAAAAAAAAAAAAAAE="
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/create-asset-trustline", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.CreateAssetTrustline(context.Background(), address)
	require.NoError(t, err)
	require.Equal(t, xdr, response.XDR)
}

func TestWallets_MintUsdbStellar(t *testing.T) {
	instanceID := "in_000000000000"
	address := "GCDNJUBQSX7AJWLJACMJ7I4BC3Z47BQUTMHEICZLE6MU4KQBRYG5JY6B"
	amount := "1000000"
	signedXDR := "AAAAAgAAAABqVFqpZzXx+KxRjXXFGO3sKwHCEYdHsWxDRrJTLGPDowAAAGQABVECAAAAAQAAAAEAAAAAAAAAAAAAAABmWFbUAAAAAAAAAAEAAAAAAAAABgAAAAFVU0RCAAAAAABbjPEfrLNLCLjNQyaWWgTeFn4tnbFnNd9FTJ3HgkLUCwAAAAAAAAAAAAAAAAAAAAE="

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"address":"GCDNJUBQSX7AJWLJACMJ7I4BC3Z47BQUTMHEICZLE6MU4KQBRYG5JY6B",
					"amount":"1000000",
					"signedXdr":"AAAAAgAAAABqVFqpZzXx+KxRjXXFGO3sKwHCEYdHsWxDRrJTLGPDowAAAGQABVECAAAAAQAAAAEAAAAAAAAAAAAAAABmWFbUAAAAAAAAAAEAAAAAAAAABgAAAAFVU0RCAAAAAABbjPEfrLNLCLjNQyaWWgTeFn4tnbFnNd9FTJ3HgkLUCwAAAAAAAAAAAAAAAAAAAAE="
				}`),
				Out:    json.RawMessage(`{}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/mint-usdb-stellar", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.MintUsdbStellar(context.Background(), &MintUsdbStellarParams{
		Address:   address,
		Amount:    amount,
		SignedXDR: signedXDR,
	})
	require.NoError(t, err)
}

func TestWallets_MintUsdbSolana(t *testing.T) {
	instanceID := "in_000000000000"
	address := "7YttLkHDoNj9wyDur5pM1ejNaAvT9X4eqaYcHQqtj2G5"
	amount := "1000000"
	signature := "4wceVEQeJG4vpS4k2o1dHU5cFWeWTQU8iaCEpRaV5KkqSxPfbdAc8hzXa7nNYG6rvqgAmDkzBycbcXkKKAeK8Jtu"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"address":"7YttLkHDoNj9wyDur5pM1ejNaAvT9X4eqaYcHQqtj2G5",
					"amount":"1000000"
				}`),
				Out: json.RawMessage(`{
					"success":true,
					"signature":"4wceVEQeJG4vpS4k2o1dHU5cFWeWTQU8iaCEpRaV5KkqSxPfbdAc8hzXa7nNYG6rvqgAmDkzBycbcXkKKAeK8Jtu",
					"error":""
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/mint-usdb-solana", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.MintUsdbSolana(context.Background(), &MintUsdbSolanaParams{
		Address: address,
		Amount:  amount,
	})
	require.NoError(t, err)
	require.True(t, response.Success)
	require.Equal(t, signature, response.Signature)
	require.Empty(t, response.Error)
}

func TestWallets_PrepareSolanaDelegationTransaction(t *testing.T) {
	instanceID := "in_000000000000"
	tokenAddress := "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA"
	amount := "1000000"
	ownerAddress := "7YttLkHDoNj9wyDur5pM1ejNaAvT9X4eqaYcHQqtj2G5"
	transaction := "AAGBf4K95Gp5i6f0BAEYAgABAgMEBQYHCAkKCwwNDg8QERITFBUWFxgZGhscHR4fICEiIw=="

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"token_address":"TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA",
					"amount":"1000000",
					"owner_address":"7YttLkHDoNj9wyDur5pM1ejNaAvT9X4eqaYcHQqtj2G5"
				}`),
				Out: json.RawMessage(`{
					"success":true,
					"transaction":"AAGBf4K95Gp5i6f0BAEYAgABAgMEBQYHCAkKCwwNDg8QERITFBUWFxgZGhscHR4fICEiIw=="
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/prepare-delegate-solana", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.PrepareSolanaDelegationTransaction(context.Background(), &PrepareSolanaDelegationTransactionParams{
		TokenAddress: tokenAddress,
		Amount:       amount,
		OwnerAddress: ownerAddress,
	})
	require.NoError(t, err)
	require.True(t, response.Success)
	require.Equal(t, transaction, response.Transaction)
}

func TestOfframpWallets_List(t *testing.T) {
	receiverID := "re_000000000000"
	bankAccountID := "ba_000000000000"
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
						"id":"ow_000000000000",
						"external_id":"your_external_id",
						"instance_id":"in_000000000000",
						"receiver_id":"re_000000000000",
						"bank_account_id":"ba_000000000000",
						"network":"tron",
						"address":"TALJN9zTTEL9TVBb4WuTt6wLvPqJZr3hvb",
						"created_at":"2021-01-01T00:00:00Z",
						"updated_at":"2021-01-01T00:00:00Z"
					}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s/offramp-wallets", instanceID, receiverID, bankAccountID),
			},
		},
		UserAgent: "test",
	}

	client := NewOfframpClient(cfg)
	wallets, err := client.List(context.Background(), receiverID, bankAccountID)
	require.NoError(t, err)
	require.Len(t, wallets, 1)
	require.Equal(t, "ow_000000000000", wallets[0].ID)
	require.Equal(t, "your_external_id", wallets[0].ExternalID)
	require.Equal(t, "in_000000000000", wallets[0].InstanceID)
	require.Equal(t, "re_000000000000", wallets[0].ReceiverID)
	require.Equal(t, "ba_000000000000", wallets[0].BankAccountID)
	require.Equal(t, "tron", wallets[0].Network)
	require.Equal(t, "TALJN9zTTEL9TVBb4WuTt6wLvPqJZr3hvb", wallets[0].Address)
}

func TestOfframpWallets_Create(t *testing.T) {
	receiverID := "re_000000000000"
	bankAccountID := "ba_000000000000"
	instanceID := "in_000000000000"
	externalID := "your_external_id"
	address := "TALJN9zTTEL9TVBb4WuTt6wLvPqJZr3hvb"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"external_id":"your_external_id",
					"network":"tron"
				}`),
				Out: json.RawMessage(`{
					"id":"ow_000000000000",
					"external_id":"your_external_id",
					"network":"tron",
					"address":"TALJN9zTTEL9TVBb4WuTt6wLvPqJZr3hvb"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s/offramp-wallets", instanceID, receiverID, bankAccountID),
			},
		},
		UserAgent: "test",
	}

	client := NewOfframpClient(cfg)
	wallet, err := client.Create(context.Background(), &CreateOfframpWalletParams{
		ReceiverID:    receiverID,
		BankAccountID: bankAccountID,
		ExternalID:    externalID,
		Network:       "tron",
	})
	require.NoError(t, err)
	require.Equal(t, "ow_000000000000", wallet.ID)
	require.Equal(t, externalID, wallet.ExternalID)
	require.Equal(t, "tron", wallet.Network)
	require.Equal(t, address, wallet.Address)
}

func TestOfframpWallets_Get(t *testing.T) {
	receiverID := "re_000000000000"
	bankAccountID := "ba_000000000000"
	id := "ow_000000000000"
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test-key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"ow_000000000000",
					"external_id":"your_external_id",
					"instance_id":"in_000000000000",
					"receiver_id":"re_000000000000",
					"bank_account_id":"ba_000000000000",
					"network":"tron",
					"address":"TALJN9zTTEL9TVBb4WuTt6wLvPqJZr3hvb",
					"created_at":"2021-01-01T00:00:00Z",
					"updated_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s/offramp-wallets/%s", instanceID, receiverID, bankAccountID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewOfframpClient(cfg)
	wallet, err := client.Get(context.Background(), receiverID, bankAccountID, id)
	require.NoError(t, err)
	require.Equal(t, id, wallet.ID)
	require.Equal(t, "your_external_id", wallet.ExternalID)
	require.Equal(t, "in_000000000000", wallet.InstanceID)
	require.Equal(t, receiverID, wallet.ReceiverID)
	require.Equal(t, bankAccountID, wallet.BankAccountID)
	require.Equal(t, "tron", wallet.Network)
	require.Equal(t, "TALJN9zTTEL9TVBb4WuTt6wLvPqJZr3hvb", wallet.Address)
}
