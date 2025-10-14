package receivers

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

func TestReceivers_List(t *testing.T) {
	instanceID := "inst_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"id":"recv_123",
						"type":"individual",
						"kyc_type":"standard",
						"kyc_status":"approved",
						"kyc_warnings":null,
						"email":"test@example.com",
						"tax_id":"123456789",
						"address_line_1":"123 Main St",
						"city":"New York",
						"state_province_region":"NY",
						"country":"US",
						"postal_code":"10001",
						"phone_number":"1234567890",
						"proof_of_address_doc_type":"UTILITY_BILL",
						"proof_of_address_doc_file":"file123",
						"first_name":"John",
						"last_name":"Doe",
						"date_of_birth":"1990-01-01",
						"id_doc_country":"US",
						"id_doc_type":"PASSPORT",
						"id_doc_front_file":"file456",
						"aiprise_validation_key":"key123",
						"instance_id":"inst_123",
						"created_at":"2024-01-01T00:00:00Z",
						"updated_at":"2024-01-01T00:00:00Z",
						"limit":{"per_transaction":10000,"daily":50000,"monthly":200000}
					}
				]`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	receivers, err := client.List(context.Background())
	require.NoError(t, err)
	require.Len(t, receivers, 1)
	require.Equal(t, "recv_123", receivers[0].ID)
	require.Equal(t, "John", receivers[0].FirstName)
}

func TestReceivers_CreateIndividualStandard(t *testing.T) {
	instanceID := "inst_123"
	receiverID := "recv_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, receiverID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.CreateIndividualStandard(context.Background(), &CreateIndividualStandardParams{
		Email:                 "test@example.com",
		FirstName:             "John",
		LastName:              "Doe",
		DateOfBirth:           "1990-01-01",
		TaxID:                 "123456789",
		AddressLine1:          "123 Main St",
		City:                  "New York",
		StateProvinceRegion:   "NY",
		Country:               types.CountryUS,
		PostalCode:            "10001",
		IDDocCountry:          types.CountryUS,
		IDDocType:             IdentificationDocumentPassport,
		IDDocFrontFile:        "file123",
		ProofOfAddressDocType: ProofOfAddressDocTypeUtilityBill,
		ProofOfAddressDocFile: "file456",
		TosID:                 "tos_123",
	})
	require.NoError(t, err)
	require.Equal(t, receiverID, response.ID)
}

func TestReceivers_CreateIndividualEnhanced(t *testing.T) {
	instanceID := "inst_123"
	receiverID := "recv_456"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, receiverID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.CreateIndividualEnhanced(context.Background(), &CreateIndividualEnhancedParams{
		Email:                         "test@example.com",
		FirstName:                     "Jane",
		LastName:                      "Smith",
		DateOfBirth:                   "1985-05-15",
		TaxID:                         "987654321",
		AddressLine1:                  "456 Oak Ave",
		City:                          "Los Angeles",
		StateProvinceRegion:           "CA",
		Country:                       types.CountryUS,
		PostalCode:                    "90001",
		IDDocCountry:                  types.CountryUS,
		IDDocType:                     IdentificationDocumentIDCard,
		IDDocFrontFile:                "file789",
		ProofOfAddressDocType:         ProofOfAddressDocTypeBankStatement,
		ProofOfAddressDocFile:         "file012",
		IndividualHoldingDocFrontFile: "file345",
		SourceOfFundsDocType:          SourceOfFundsSalary,
		SourceOfFundsDocFile:          "file678",
		PurposeOfTransactions:         PurposeBusinessTransactions,
		TosID:                         "tos_456",
	})
	require.NoError(t, err)
	require.Equal(t, receiverID, response.ID)
}

func TestReceivers_CreateBusinessStandard(t *testing.T) {
	instanceID := "inst_123"
	receiverID := "recv_789"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(fmt.Sprintf(`{"id":"%s"}`, receiverID)),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers", instanceID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	response, err := client.CreateBusinessStandard(context.Background(), &CreateBusinessStandardParams{
		Email:                   "business@example.com",
		LegalName:               "ACME Corp",
		AlternateName:           "ACME",
		TaxID:                   "12-3456789",
		FormationDate:           "2020-01-01",
		AddressLine1:            "789 Business Blvd",
		City:                    "Chicago",
		StateProvinceRegion:     "IL",
		Country:                 types.CountryUS,
		PostalCode:              "60601",
		ProofOfAddressDocType:   ProofOfAddressDocTypeUtilityBill,
		ProofOfAddressDocFile:   "file111",
		IncorporationDocFile:    "file222",
		ProofOfOwnershipDocFile: "file333",
		Owners:                  []Owner{},
		TosID:                   "tos_789",
	})
	require.NoError(t, err)
	require.Equal(t, receiverID, response.ID)
}

func TestReceivers_Get(t *testing.T) {
	instanceID := "inst_123"
	receiverID := "recv_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(fmt.Sprintf(`{
					"id":"%s",
					"type":"individual",
					"kyc_type":"standard",
					"kyc_status":"approved",
					"kyc_warnings":null,
					"email":"test@example.com",
					"tax_id":"123456789",
					"address_line_1":"123 Main St",
					"city":"New York",
					"state_province_region":"NY",
					"country":"US",
					"postal_code":"10001",
					"proof_of_address_doc_type":"UTILITY_BILL",
					"proof_of_address_doc_file":"file123",
					"first_name":"John",
					"last_name":"Doe",
					"date_of_birth":"1990-01-01",
					"id_doc_country":"US",
					"id_doc_type":"PASSPORT",
					"id_doc_front_file":"file456",
					"aiprise_validation_key":"key123",
					"instance_id":"inst_123",
					"created_at":"2024-01-01T00:00:00Z",
					"updated_at":"2024-01-01T00:00:00Z",
					"limit":{"per_transaction":10000,"daily":50000,"monthly":200000}
				}`, receiverID)),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	receiver, err := client.Get(context.Background(), receiverID)
	require.NoError(t, err)
	require.Equal(t, receiverID, receiver.ID)
	require.Equal(t, "John", receiver.FirstName)
	require.Equal(t, "Doe", receiver.LastName)
}

func TestReceivers_Update(t *testing.T) {
	instanceID := "inst_123"
	receiverID := "recv_123"
	newEmail := "newemail@example.com"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodPatch,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Update(context.Background(), &UpdateParams{
		ReceiverID: receiverID,
		Email:      &newEmail,
	})
	require.NoError(t, err)
}

func TestReceivers_Delete(t *testing.T) {
	instanceID := "inst_123"
	receiverID := "recv_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{}`),
				Method: http.MethodDelete,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	err := client.Delete(context.Background(), receiverID)
	require.NoError(t, err)
}

func TestReceivers_GetLimits(t *testing.T) {
	instanceID := "inst_123"
	receiverID := "recv_123"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"limits":{
						"payin":{"daily":50000,"monthly":200000},
						"payout":{"daily":100000,"monthly":500000}
					}
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/limits/receivers/%s", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	limits, err := client.GetLimits(context.Background(), receiverID)
	require.NoError(t, err)
	require.Equal(t, 50000.0, limits.Limits.Payin.Daily)
	require.Equal(t, 100000.0, limits.Limits.Payout.Daily)
}
