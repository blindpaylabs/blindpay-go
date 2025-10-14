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
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"pix",
					"name":"PIX Account",
					"pix_key":"14947677768"
				}`),
				Out: json.RawMessage(`{
					"id": "ba_000000000000",
					"type": "pix",
					"name": "PIX Account",
					"pix_key": "14947677768",
					"created_at": "2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreatePix(context.Background(), &CreatePixParams{
		ReceiverID: receiverID,
		Name:       "PIX Account",
		PixKey:     "14947677768",
	})
	require.NoError(t, err)
	require.Equal(t, "ba_000000000000", account.ID)
	require.Equal(t, "PIX Account", account.Name)
	require.Equal(t, "14947677768", account.PixKey)
}

func TestBankAccounts_CreateArgentinaTransfers(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"transfers_bitso",
					"name":"Argentina Transfers Account",
					"beneficiary_name":"Individual full name or business name",
					"transfers_type":"CVU",
					"transfers_account":"BM123123123123"
				}`),
				Out: json.RawMessage(`{
					"id":"ba_000000000000",
					"type":"transfers_bitso",
					"name":"Argentina Transfers Account",
					"beneficiary_name":"Individual full name or business name",
					"transfers_type":"CVU",
					"transfers_account":"BM123123123123",
					"created_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreateArgentinaTransfers(context.Background(), &CreateArgentinaTransfersParams{
		ReceiverID:       receiverID,
		Name:             "Argentina Transfers Account",
		BeneficiaryName:  "Individual full name or business name",
		TransfersType:    ArgentinaTransfersCVU,
		TransfersAccount: "BM123123123123",
	})
	require.NoError(t, err)
	require.Equal(t, "ba_000000000000", account.ID)
}

func TestCreateSpei(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"spei_bitso",
					"name":"SPEI Account",
					"beneficiary_name":"Individual full name or business name",
					"spei_protocol":"clabe",
					"spei_institution_code":"40002",
					"spei_clabe":"5482347403740546"
				}`),
				Out: json.RawMessage(`{
					"id":"ba_000000000000",
					"type":"spei_bitso",
					"name":"SPEI Account",
					"beneficiary_name":"Individual full name or business name",
					"spei_protocol":"clabe",
					"spei_institution_code":"40002",
					"spei_clabe":"5482347403740546",
					"created_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreateSpei(context.Background(), &CreateSpeiParams{
		ReceiverID:          receiverID,
		BeneficiaryName:     "Individual full name or business name",
		Name:                "SPEI Account",
		SpeiClabe:           "5482347403740546",
		SpeiInstitutionCode: "40002",
		SpeiProtocol:        SpeiProtocolClabe,
	})
	require.NoError(t, err)
	require.Equal(t, "ba_000000000000", account.ID)
}

func TestCreateColombiaAch(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"ach_cop_bitso",
					"name":"Colombia ACH Account",
					"account_type":"checking",
					"ach_cop_beneficiary_first_name":"Fernando",
					"ach_cop_beneficiary_last_name":"Guzman Alarcón",
					"ach_cop_document_id":"1661105408",
					"ach_cop_document_type":"CC",
					"ach_cop_email":"fernando.guzman@gmail.com",
					"ach_cop_bank_code":"051",
					"ach_cop_bank_account":"12345678"
				}`),
				Out: json.RawMessage(`{
					"id":"ba_000000000000",
					"type":"ach_cop_bitso",
					"name":"Colombia ACH Account",
					"account_type":"checking",
					"ach_cop_beneficiary_first_name":"Fernando",
					"ach_cop_beneficiary_last_name":"Guzman Alarcón",
					"ach_cop_document_id":"1661105408",
					"ach_cop_document_type":"CC",
					"ach_cop_email":"fernando.guzman@gmail.com",
					"ach_cop_bank_code":"051",
					"ach_cop_bank_account":"12345678",
					"created_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreateColombiaAch(context.Background(), &CreateColombiaAchParams{
		ReceiverID:                 receiverID,
		Name:                       "Colombia ACH Account",
		AccountType:                types.BankAccountTypeChecking,
		AchCopBeneficiaryFirstName: "Fernando",
		AchCopBeneficiaryLastName:  "Guzman Alarcón",
		AchCopDocumentID:           "1661105408",
		AchCopDocumentType:         AchCopDocumentCC,
		AchCopEmail:                "fernando.guzman@gmail.com",
		AchCopBankCode:             "051",
		AchCopBankAccount:          "12345678",
	})
	require.NoError(t, err)
	require.Equal(t, "ba_000000000000", account.ID)
}

func TestCreateAch(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"ach",
					"name":"ACH Account",
					"account_class":"individual",
					"account_number":"1001001234",
					"account_type":"checking",
					"beneficiary_name":"Individual full name or business name",
					"routing_number":"012345678"
				}`),
				Out: json.RawMessage(`{
					"id":"ba_000000000000",
					"type":"ach",
					"name":"ACH Account",
					"beneficiary_name":"Individual full name or business name",
					"routing_number":"012345678",
					"account_number":"1001001234",
					"account_type":"checking",
					"account_class":"individual",
					"address_line_1":null,
					"address_line_2":null,
					"city":null,
					"state_province_region":null,
					"country":null,
					"postal_code":null,
					"ach_cop_beneficiary_first_name":null,
					"ach_cop_beneficiary_last_name":null,
					"ach_cop_document_id":null,
					"ach_cop_document_type":null,
					"ach_cop_email":null,
					"ach_cop_bank_code":null,
					"ach_cop_bank_account":null,
					"created_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreateAch(context.Background(), &CreateAchParams{
		ReceiverID:      receiverID,
		Name:            "ACH Account",
		AccountClass:    types.AccountClassIndividual,
		AccountNumber:   "1001001234",
		AccountType:     types.BankAccountTypeChecking,
		BeneficiaryName: "Individual full name or business name",
		RoutingNumber:   "012345678",
	})
	require.NoError(t, err)
	require.Equal(t, "ba_000000000000", account.ID)
}

func TestCreateWire(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"wire",
					"name":"Wire Account",
					"account_number":"1001001234",
					"beneficiary_name":"Individual full name or business name",
					"routing_number":"012345678",
					"address_line_1":"Address line 1",
					"address_line_2":"Address line 2",
					"city":"City",
					"state_province_region":"State/Province/Region",
					"country":"US",
					"postal_code":"Postal code"
				}`),
				Out: json.RawMessage(`{
					"id":"ba_000000000000",
					"type":"wire",
					"name":"Wire Account",
					"beneficiary_name":"Individual full name or business name",
					"routing_number":"012345678",
					"account_number":"1001001234",
					"address_line_1":"Address line 1",
					"address_line_2":"Address line 2",
					"city":"City",
					"state_province_region":"State/Province/Region",
					"country":"US",
					"postal_code":"Postal code",
					"created_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreateWire(context.Background(), &CreateWireParams{
		ReceiverID:          receiverID,
		Name:                "Wire Account",
		AccountNumber:       "1001001234",
		BeneficiaryName:     "Individual full name or business name",
		RoutingNumber:       "012345678",
		AddressLine1:        "Address line 1",
		AddressLine2:        "Address line 2",
		City:                "City",
		StateProvinceRegion: "State/Province/Region",
		Country:             types.CountryUS,
		PostalCode:          "Postal code",
	})
	require.NoError(t, err)
	require.Equal(t, "ba_000000000000", account.ID)
}

func TestCreateInternationalSwift(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"international_swift",
					"name":"International Swift Account",
					"swift_account_holder_name":"John Doe",
					"swift_account_number_iban":"123456789",
					"swift_bank_address_line_1":"123 Main Street, Suite 100, Downtown District, City Center CP 12345",
					"swift_bank_city":"City",
					"swift_bank_country":"MX",
					"swift_bank_name":"Banco Regional SA",
					"swift_bank_postal_code":"11530",
					"swift_bank_state_province_region":"District",
					"swift_beneficiary_address_line_1":"123 Main Street, Suite 100, Downtown District, City Center CP 12345",
					"swift_beneficiary_city":"City",
					"swift_beneficiary_country":"MX",
					"swift_beneficiary_postal_code":"11530",
					"swift_beneficiary_state_province_region":"District",
					"swift_code_bic":"123456789",
					"swift_intermediary_bank_account_number_iban":null,
					"swift_intermediary_bank_country":null,
					"swift_intermediary_bank_name":null,
					"swift_intermediary_bank_swift_code_bic":null
				}`),
				Out: json.RawMessage(`{
					"id":"ba_000000000000",
					"type":"international_swift",
					"name":"International Swift Account",
					"beneficiary_name":null,
					"address_line_1":null,
					"address_line_2":null,
					"city":null,
					"state_province_region":null,
					"country":null,
					"postal_code":null,
					"swift_code_bic":"123456789",
					"swift_account_holder_name":"John Doe",
					"swift_account_number_iban":"123456789",
					"swift_beneficiary_address_line_1":"123 Main Street, Suite 100, Downtown District, City Center CP 12345",
					"swift_beneficiary_address_line_2":null,
					"swift_beneficiary_country":"MX",
					"swift_beneficiary_city":"City",
					"swift_beneficiary_state_province_region":"District",
					"swift_beneficiary_postal_code":"11530",
					"swift_bank_name":"Banco Regional SA",
					"swift_bank_address_line_1":"123 Main Street, Suite 100, Downtown District, City Center CP 12345",
					"swift_bank_address_line_2":null,
					"swift_bank_country":"MX",
					"swift_bank_city":"City",
					"swift_bank_state_province_region":"District",
					"swift_bank_postal_code":"11530",
					"swift_intermediary_bank_swift_code_bic":null,
					"swift_intermediary_bank_account_number_iban":null,
					"swift_intermediary_bank_name":null,
					"swift_intermediary_bank_country":null,
					"created_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreateInternationalSwift(context.Background(), &CreateInternationalSwiftParams{
		ReceiverID:                             receiverID,
		Name:                                   "International Swift Account",
		SwiftAccountHolderName:                 "John Doe",
		SwiftAccountNumberIban:                 "123456789",
		SwiftBankAddressLine1:                  "123 Main Street, Suite 100, Downtown District, City Center CP 12345",
		SwiftBankCity:                          "City",
		SwiftBankCountry:                       types.CountryMX,
		SwiftBankName:                          "Banco Regional SA",
		SwiftBankPostalCode:                    "11530",
		SwiftBankStateProvinceRegion:           "District",
		SwiftBeneficiaryAddressLine1:           "123 Main Street, Suite 100, Downtown District, City Center CP 12345",
		SwiftBeneficiaryCity:                   "City",
		SwiftBeneficiaryCountry:                types.CountryMX,
		SwiftBeneficiaryPostalCode:             "11530",
		SwiftBeneficiaryStateProvinceRegion:    "District",
		SwiftCodeBic:                           "123456789",
		SwiftIntermediaryBankAccountNumberIban: nil,
		SwiftIntermediaryBankCountry:           nil,
		SwiftIntermediaryBankName:              nil,
		SwiftIntermediaryBankSwiftCodeBic:      nil,
	})
	require.NoError(t, err)
	require.Equal(t, "ba_000000000000", account.ID)
}

func TestBankAccounts_CreateRtp(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				In: json.RawMessage(`{
					"type":"rtp",
					"name":"John Doe RTP",
					"beneficiary_name":"John Doe",
					"routing_number":"121000358",
					"account_number":"325203027578",
					"address_line_1":"Street of the fools",
					"city":"Fools City",
					"state_province_region":"FL",
					"country":"US",
					"postal_code":"22599"
				}`),
				Out: json.RawMessage(`{
					"id":"ba_JW5ZtlKMlgS1",
					"type":"rtp",
					"name":"John Doe RTP",
					"beneficiary_name":"John Doe",
					"routing_number":"121000358",
					"account_number":"325203027578",
					"address_line_1":"Street of the fools",
					"address_line_2":null,
					"city":"Fools City",
					"state_province_region":"FL",
					"country":"US",
					"postal_code":"22599",
					"created_at":"2025-09-30T04:23:30.823Z"
				}`),
				Method: http.MethodPost,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.CreateRtp(context.Background(), &CreateRtpParams{
		ReceiverID:          receiverID,
		Name:                "John Doe RTP",
		BeneficiaryName:     "John Doe",
		RoutingNumber:       "121000358",
		AccountNumber:       "325203027578",
		AddressLine1:        "Street of the fools",
		City:                "Fools City",
		StateProvinceRegion: "FL",
		Country:             types.CountryUS,
		PostalCode:          "22599",
	})
	require.NoError(t, err)
	require.Equal(t, "ba_JW5ZtlKMlgS1", account.ID)
}

func TestBankAccounts_Get(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"
	id := "ba_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T: t,
				Out: json.RawMessage(`{
					"id":"ba_000000000000",
					"receiver_id":"rcv_123",
					"account_holder_name":"Individual full name or business name",
					"account_number":"1001001234",
					"routing_number":"012345678",
					"account_type":"checking",
					"bank_name":"Bank Name",
					"swift_code":"123456789",
					"iban":null,
					"is_primary":false,
					"created_at":"2021-01-01T00:00:00Z",
					"updated_at":"2021-01-01T00:00:00Z"
				}`),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts/%s", instanceID, receiverID, id),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	account, err := client.Get(context.Background(), receiverID, id)
	require.NoError(t, err)
	require.Equal(t, "ba_000000000000", account.ID)
	require.Equal(t, "rcv_123", account.ReceiverID)
}

func TestBankAccounts_List(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"

	mockOut := `[
		{
			"id": "ba_000000000000",
			"type": "wire",
			"name": "Bank Account Name",
			"pix_key": "14947677768",
			"beneficiary_name": "Individual full name or business name",
			"routing_number": "012345678",
			"account_number": "1001001234",
			"account_type": "checking",
			"account_class": "individual",
			"address_line_1": "Address line 1",
			"address_line_2": "Address line 2",
			"city": "City",
			"state_province_region": "State/Province/Region",
			"country": "US",
			"postal_code": "Postal code",
			"spei_protocol": "clabe",
			"spei_institution_code": "40002",
			"spei_clabe": "5482347403740546",
			"transfers_type": "CVU",
			"transfers_account": "BM123123123123",
			"ach_cop_beneficiary_first_name": "Fernando",
			"ach_cop_beneficiary_last_name": "Guzman Alarcón",
			"ach_cop_document_id": "1661105408",
			"ach_cop_document_type": "CC",
			"ach_cop_email": "fernando.guzman@gmail.com",
			"ach_cop_bank_code": "051",
			"ach_cop_bank_account": "12345678",
			"swift_code_bic": "123456789",
			"swift_account_holder_name": "John Doe",
			"swift_account_number_iban": "123456789",
			"swift_beneficiary_address_line_1": "123 Main Street, Suite 100, Downtown District, City Center CP 12345",
			"swift_beneficiary_address_line_2": "456 Oak Avenue, Building 7, Financial District, Business Center CP 54321",
			"swift_beneficiary_country": "MX",
			"swift_beneficiary_city": "City",
			"swift_beneficiary_state_province_region": "District",
			"swift_beneficiary_postal_code": "11530",
			"swift_bank_name": "Banco Regional SA",
			"swift_bank_address_line_1": "123 Main Street, Suite 100, Downtown District, City Center CP 12345",
			"swift_bank_address_line_2": "456 Oak Avenue, Building 7, Financial District, Business Center CP 54321",
			"swift_bank_country": "MX",
			"swift_bank_city": "City",
			"swift_bank_state_province_region": "District",
			"swift_bank_postal_code": "11530",
			"swift_intermediary_bank_swift_code_bic": "AEIBARB1",
			"swift_intermediary_bank_account_number_iban": "123456789",
			"swift_intermediary_bank_name": "Banco Regional SA",
			"swift_intermediary_bank_country": "US",
			"tron_wallet_hash": "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
			"offramp_wallets": [
				{
					"address": "TALJN9zTTEL9TVBb4WuTt6wLvPqJZr3hvb",
					"id": "ow_000000000000",
					"network": "tron",
					"external_id": "your_external_id"
				}
			],
			"created_at": "2021-01-01T00:00:00Z"
		}
	]`

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(mockOut),
				Method: http.MethodGet,
				Path:   fmt.Sprintf("/instances/%s/receivers/%s/bank-accounts", instanceID, receiverID),
			},
		},
		UserAgent: "test",
	}

	client := NewClient(cfg)
	bankAccounts, err := client.List(context.Background(), receiverID)
	require.NoError(t, err)
	require.Len(t, bankAccounts, 1)
	require.Equal(t, "ba_000000000000", bankAccounts[0].ID)
}

func TestBankAccounts_Delete(t *testing.T) {
	instanceID := "in_000000000000"
	receiverID := "re_000000000000"
	id := "ba_000000000000"

	cfg := &config.Config{
		BaseURL:    "https://api.blindpay.com",
		InstanceID: instanceID,
		APIKey:     "test-key",
		HTTPClient: &http.Client{
			Transport: &blindpaytest.RoundTripper{
				T:      t,
				Out:    json.RawMessage(`{"data": null}`),
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
