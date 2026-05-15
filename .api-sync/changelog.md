# API Changes (SDK-relevant)

## New Endpoints

- **GET /v1/instances/{instance_id}/receivers/{receiver_id}/rfi** — Get Open RFI for Receiver
  Response (Rfi):
    - id (string)
    - receiver_id (string)
    - instance_id (string)
    - status (string, enum: pending | submitted | expired | cancelled)
    - request (array, items: RfiSection)
    - response (object)
    - expires_at (string, format: date-time)
    - submitted_at (string/null, format: date-time)
    - created_at (string, format: date-time)
    - receiver_type (string, enum: individual | business)
    - receiver_aiprise_session_id (string/null)
    - receiver_kyc_status (string)
- **POST /v1/instances/{instance_id}/receivers/{receiver_id}/rfi** — Submit RFI Response
  Response (inline):
    - success (boolean)

## Modified Endpoints

### /v1/instances/{instance_id}/receivers

  Request body (CreateReceiverIn) [POST]:
  - ENUM business_industry: added 1 values: 446120

### /v1/instances/{instance_id}/receivers/{id}

  Response (ReceiverOut) [GET]:
  - ENUM business_industry: added 1 values: 446120
  Request body (UpdateReceiverIn) [PUT]:
  - ENUM business_industry: added 1 values: 446120

### /v1/instances/{instance_id}/receivers/{receiver_id}/bank-accounts

  Request body (CreateBankAccountIn) [POST]:
  - ADDED field: swift_ifsc_branch_code (string/null, optional)
  - ENUM business_industry: added 1 values: 446120
  - ENUM swift_payment_code: added 4 values: hk_swift_charitabledonation, hk_swift_goods, hk_swift_personal, hk_swift_services
  Response (BankAccountOut) [POST]:
  - ADDED field: swift_ifsc_branch_code (string/null, optional)
  - ENUM business_industry: added 1 values: 446120
  - ENUM swift_payment_code: added 4 values: hk_swift_charitabledonation, hk_swift_goods, hk_swift_personal, hk_swift_services

### /v1/instances/{instance_id}/receivers/{receiver_id}/bank-accounts/{id}

  Response (BankAccountOut) [GET]:
  - ADDED field: swift_ifsc_branch_code (string/null, optional)
  - ENUM business_industry: added 1 values: 446120
  - ENUM swift_payment_code: added 4 values: hk_swift_charitabledonation, hk_swift_goods, hk_swift_personal, hk_swift_services

## Enum Value Changes

These enum fields gained or lost values across all schemas:

  - business_industry: ADDED 1 values: 446120
  - decision: ADDED 2 values: approved, rejected
  - receiver_type: ADDED 2 values: business, individual
  - swift_payment_code: ADDED 4 values: hk_swift_charitabledonation, hk_swift_goods, hk_swift_personal, hk_swift_services

## New Schemas

- **CreateComplaintIn** (5 fields)
- **CreateExternalReceiverTokenIn** (2 fields)
- **KycDecisionBody** (3 fields)
- **UpdateDynamicRateOtcIn** (1 fields)
