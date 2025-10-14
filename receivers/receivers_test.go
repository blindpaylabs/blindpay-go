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
	instanceID := "in_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`[
					{
						"id":"re_Euw7HN4OdxPn",
						"type":"individual",
						"kyc_type":"standard",
						"kyc_status":"verifying",
						"kyc_warnings":[
							{
								"code":null,
								"message":null,
								"resolution_status":null,
								"warning_id":null
							}
						],
						"email":"bernardo@gmail.com",
						"tax_id":"12345678900",
						"address_line_1":"Av. Paulista, 1000",
						"address_line_2":"Apto 101",
						"city":"São Paulo",
						"state_province_region":"SP",
						"country":"BR",
						"postal_code":"01310-100",
						"ip_address":"127.0.0.1",
						"image_url":"https://example.com/image.png",
						"phone_number":"+5511987654321",
						"proof_of_address_doc_type":"UTILITY_BILL",
						"proof_of_address_doc_file":"https://example.com/image.png",
						"first_name":"Bernardo",
						"last_name":"Simonassi",
						"date_of_birth":"1998-02-02T00:00:00.000Z",
						"id_doc_country":"BR",
						"id_doc_type":"PASSPORT",
						"id_doc_front_file":"https://example.com/image.png",
						"id_doc_back_file":"https://example.com/image.png",
						"aiprise_validation_key":"",
						"instance_id":"in_000000000000",
						"tos_id":"to_3ZZhllJkvo5Z",
						"created_at":"2021-01-01T00:00:00.000Z",
						"updated_at":"2021-01-01T00:00:00.000Z",
						"limit":{
							"per_transaction":100000,
							"daily":200000,
							"monthly":1000000
						}
					},
					{
						"id":"re_YuaMcI2B8zbQ",
						"type":"individual",
						"kyc_type":"enhanced",
						"kyc_status":"approved",
						"kyc_warnings":null,
						"email":"alice.johnson@example.com",
						"tax_id":"98765432100",
						"address_line_1":"123 Main St",
						"address_line_2":null,
						"city":"New York",
						"state_province_region":"NY",
						"country":"US",
						"postal_code":"10001",
						"ip_address":"192.168.1.1",
						"image_url":null,
						"phone_number":"+15555555555",
						"proof_of_address_doc_type":"BANK_STATEMENT",
						"proof_of_address_doc_file":"https://example.com/image.png",
						"first_name":"Alice",
						"last_name":"Johnson",
						"date_of_birth":"1990-05-10T00:00:00.000Z",
						"id_doc_country":"US",
						"id_doc_type":"PASSPORT",
						"id_doc_front_file":"https://example.com/image.png",
						"id_doc_back_file":null,
						"aiprise_validation_key":"enhanced-key",
						"instance_id":"in_000000000001",
						"source_of_funds_doc_type":"salary",
						"source_of_funds_doc_file":"https://example.com/image.png",
						"individual_holding_doc_front_file":"https://example.com/image.png",
						"purpose_of_transactions":"investment_purposes",
						"purpose_of_transactions_explanation":"Investing in stocks",
						"tos_id":"to_nppX66ntvtHs",
						"created_at":"2022-02-02T00:00:00.000Z",
						"updated_at":"2022-02-02T00:00:00.000Z",
						"limit":{
							"per_transaction":50000,
							"daily":100000,
							"monthly":500000
						}
					},
					{
						"id":"re_IOxAUL24LG7P",
						"type":"business",
						"kyc_type":"standard",
						"kyc_status":"pending",
						"kyc_warnings":null,
						"email":"business@example.com",
						"tax_id":"20096178000195",
						"address_line_1":"1 King St W",
						"address_line_2":"Suite 100",
						"city":"Toronto",
						"state_province_region":"ON",
						"country":"CA",
						"postal_code":"M5H 1A1",
						"ip_address":null,
						"image_url":null,
						"phone_number":"+14165555555",
						"proof_of_address_doc_type":"UTILITY_BILL",
						"proof_of_address_doc_file":"https://example.com/image.png",
						"legal_name":"Business Corp",
						"alternate_name":"BizCo",
						"formation_date":"2010-01-01T00:00:00.000Z",
						"website":"https://businesscorp.com",
						"owners":[
							{
								"role":"beneficial_owner",
								"first_name":"Carlos",
								"last_name":"Silva",
								"date_of_birth":"1995-05-15T00:00:00.000Z",
								"tax_id":"12345678901",
								"address_line_1":"Rua Augusta, 1500",
								"address_line_2":null,
								"city":"São Paulo",
								"state_province_region":"SP",
								"country":"BR",
								"postal_code":"01304-001",
								"id_doc_country":"BR",
								"id_doc_type":"PASSPORT",
								"id_doc_front_file":"https://example.com/image.png",
								"id_doc_back_file":"https://example.com/image.png",
								"proof_of_address_doc_type":"UTILITY_BILL",
								"proof_of_address_doc_file":"https://example.com/image.png",
								"id":"ub_000000000000",
								"instance_id":"in_000000000000",
								"receiver_id":"re_IOxAUL24LG7P"
							}
						],
						"incorporation_doc_file":"https://example.com/image.png",
						"proof_of_ownership_doc_file":"https://example.com/image.png",
						"external_id":null,
						"instance_id":"in_000000000002",
						"tos_id":"to_nppX66ntvtHs",
						"aiprise_validation_key":"",
						"created_at":"2015-03-15T00:00:00.000Z",
						"updated_at":"2015-03-15T00:00:00.000Z",
						"limit":{
							"per_transaction":200000,
							"daily":400000,
							"monthly":2000000
						}
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
	require.Len(t, receivers, 3)
	require.Equal(t, "re_Euw7HN4OdxPn", receivers[0].ID)
	require.Equal(t, "Bernardo", receivers[0].FirstName)
	require.Equal(t, "re_YuaMcI2B8zbQ", receivers[1].ID)
	require.Equal(t, "Alice", receivers[1].FirstName)
	require.Equal(t, "re_IOxAUL24LG7P", receivers[2].ID)
	require.Equal(t, "Business Corp", receivers[2].LegalName)
}

func TestReceivers_CreateIndividualStandard(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_Euw7HN4OdxPn"

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
	addressLine2 := "Apto 101"
	phoneNumber := "+5511987654321"
	idDocBackFile := "https://example.com/image.png"
	response, err := client.CreateIndividualStandard(context.Background(), &CreateIndividualStandardParams{
		Email:                 "bernardo.simonassi@gmail.com",
		FirstName:             "Bernardo",
		LastName:              "Simonassi",
		DateOfBirth:           "1998-02-02T00:00:00.000Z",
		TaxID:                 "12345678900",
		AddressLine1:          "Av. Paulista, 1000",
		AddressLine2:          &addressLine2,
		City:                  "São Paulo",
		StateProvinceRegion:   "SP",
		Country:               types.CountryBR,
		PostalCode:            "01310-100",
		PhoneNumber:           &phoneNumber,
		IDDocCountry:          types.CountryBR,
		IDDocType:             IdentificationDocumentPassport,
		IDDocFrontFile:        "https://example.com/image.png",
		IDDocBackFile:         &idDocBackFile,
		ProofOfAddressDocType: ProofOfAddressDocTypeUtilityBill,
		ProofOfAddressDocFile: "https://example.com/image.png",
		TosID:                 "to_tPiz4bM2nh5K",
	})
	require.NoError(t, err)
	require.Equal(t, receiverID, response.ID)
}

func TestReceivers_CreateIndividualEnhanced(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_YuaMcI2B8zbQ"

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
	addressLine2 := "Apto 101"
	phoneNumber := "+5511987654321"
	idDocBackFile := "https://example.com/image.png"
	purposeExplanation := "I am receiving salary payments from my employer"
	response, err := client.CreateIndividualEnhanced(context.Background(), &CreateIndividualEnhancedParams{
		Email:                            "bernardo.simonassi@gmail.com",
		FirstName:                        "Bernardo",
		LastName:                         "Simonassi",
		DateOfBirth:                      "1998-02-02T00:00:00.000Z",
		TaxID:                            "12345678900",
		AddressLine1:                     "Av. Paulista, 1000",
		AddressLine2:                     &addressLine2,
		City:                             "São Paulo",
		StateProvinceRegion:              "SP",
		Country:                          types.CountryBR,
		PostalCode:                       "01310-100",
		PhoneNumber:                      &phoneNumber,
		IDDocCountry:                     types.CountryBR,
		IDDocType:                        IdentificationDocumentPassport,
		IDDocFrontFile:                   "https://example.com/image.png",
		IDDocBackFile:                    &idDocBackFile,
		ProofOfAddressDocType:            ProofOfAddressDocTypeUtilityBill,
		ProofOfAddressDocFile:            "https://example.com/image.png",
		IndividualHoldingDocFrontFile:    "https://example.com/image.png",
		PurposeOfTransactions:            PurposePersonalOrLivingExpenses,
		SourceOfFundsDocType:             SourceOfFundsSavings,
		PurposeOfTransactionsExplanation: &purposeExplanation,
		SourceOfFundsDocFile:             "https://example.com/image.png",
		TosID:                            "to_3ZZhllJkvo5Z",
	})
	require.NoError(t, err)
	require.Equal(t, receiverID, response.ID)
}

func TestReceivers_CreateBusinessStandard(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_IOxAUL24LG7P"

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
	addressLine2 := "Sala 1201"
	website := "https://site.com/"
	ownerIDDocBackFile := "https://example.com/image.png"
	response, err := client.CreateBusinessStandard(context.Background(), &CreateBusinessStandardParams{
		Email:                   "contato@empresa.com.br",
		TaxID:                   "20096178000195",
		AddressLine1:            "Av. Brigadeiro Faria Lima, 400",
		AddressLine2:            &addressLine2,
		City:                    "São Paulo",
		StateProvinceRegion:     "SP",
		Country:                 types.CountryBR,
		PostalCode:              "04538-132",
		LegalName:               "Empresa Exemplo Ltda",
		AlternateName:           "Exemplo",
		FormationDate:           "2010-05-20T00:00:00.000Z",
		IncorporationDocFile:    "https://example.com/image.png",
		ProofOfAddressDocType:   ProofOfAddressDocTypeUtilityBill,
		ProofOfAddressDocFile:   "https://example.com/image.png",
		ProofOfOwnershipDocFile: "https://example.com/image.png",
		Website:                 &website,
		Owners: []Owner{
			{
				Role:                  "beneficial_owner",
				FirstName:             "Carlos",
				LastName:              "Silva",
				DateOfBirth:           "1995-05-15T00:00:00.000Z",
				TaxID:                 "12345678901",
				AddressLine1:          "Rua Augusta, 1500",
				AddressLine2:          nil,
				City:                  "São Paulo",
				StateProvinceRegion:   "SP",
				Country:               types.CountryBR,
				PostalCode:            "01304-001",
				IDDocCountry:          types.CountryBR,
				IDDocType:             IdentificationDocumentPassport,
				IDDocFrontFile:        "https://example.com/image.png",
				IDDocBackFile:         &ownerIDDocBackFile,
				ProofOfAddressDocType: ProofOfAddressDocTypeUtilityBill,
				ProofOfAddressDocFile: "https://example.com/image.png",
				ID:                    "ub_000000000000",
				InstanceID:            "in_000000000000",
				ReceiverID:            "re_IOxAUL24LG7P",
			},
		},
		TosID: "to_nppX66ntvtHs",
	})
	require.NoError(t, err)
	require.Equal(t, receiverID, response.ID)
}

func TestReceivers_Get(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_YuaMcI2B8zbQ"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"re_YuaMcI2B8zbQ",
					"type":"individual",
					"kyc_type":"enhanced",
					"kyc_status":"verifying",
					"kyc_warnings":[
						{
							"code":null,
							"message":null,
							"resolution_status":null,
							"warning_id":null
						}
					],
					"email":"bernardo.simonassi@gmail.com",
					"tax_id":"12345678900",
					"address_line_1":"Av. Paulista, 1000",
					"address_line_2":"Apto 101",
					"city":"São Paulo",
					"state_province_region":"SP",
					"country":"BR",
					"postal_code":"01310-100",
					"ip_address":"127.0.0.1",
					"image_url":"https://example.com/image.png",
					"phone_number":"+5511987654321",
					"proof_of_address_doc_type":"UTILITY_BILL",
					"proof_of_address_doc_file":"https://example.com/image.png",
					"first_name":"Bernardo",
					"last_name":"Simonassi",
					"date_of_birth":"1998-02-02T00:00:00.000Z",
					"id_doc_country":"BR",
					"id_doc_type":"PASSPORT",
					"id_doc_front_file":"https://example.com/image.png",
					"id_doc_back_file":"https://example.com/image.png",
					"aiprise_validation_key":"",
					"source_of_funds_doc_type":"savings",
					"source_of_funds_doc_file":"https://example.com/image.png",
					"individual_holding_doc_front_file":"https://example.com/image.png",
					"purpose_of_transactions":"personal_or_living_expenses",
					"purpose_of_transactions_explanation":"I am receiving salary payments from my employer",
					"instance_id":"in_000000000000",
					"tos_id":"to_3ZZhllJkvo5Z",
					"created_at":"2021-01-01T00:00:00.000Z",
					"updated_at":"2021-01-01T00:00:00.000Z",
					"limit":{
						"per_transaction":100000,
						"daily":200000,
						"monthly":1000000
					}
				}`),
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
	require.Equal(t, "Bernardo", receiver.FirstName)
	require.Equal(t, "Simonassi", receiver.LastName)
}

func TestReceivers_Update(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_YuaMcI2B8zbQ"
	newEmail := "bernardo.simonassi@gmail.com"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
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
	instanceID := "in_000000000000"
	receiverID := "re_YuaMcI2B8zbQ"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data":null}`),
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
	instanceID := "in_000000000000"
	receiverID := "re_YuaMcI2B8zbQ"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		APIKey:     "test_key",
		InstanceID: instanceID,
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"limits":{
						"payin":{
							"daily":10000,
							"monthly":50000
						},
						"payout":{
							"daily":20000,
							"monthly":100000
						}
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
	require.Equal(t, 10000.0, limits.Limits.Payin.Daily)
	require.Equal(t, 20000.0, limits.Limits.Payout.Daily)
}
