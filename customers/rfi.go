package customers

import (
	"context"
	"fmt"

	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/rfi"
)

// GetRfi retrieves a customer's open request for information. It returns
// nil, nil when no RFI is open.
func (c *Client) GetRfi(ctx context.Context, customerID string, opts ...option.RequestOption) (*rfi.CustomerRfi, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}

	path := request.Pathf("/instances/%s/customers/%s/rfi", c.instanceID, customerID)
	return request.Do[*rfi.CustomerRfi](c.cfg, ctx, "GET", path, nil, opts...)
}

// SubmitRfi answers a customer's open RFI. response maps each field key to
// its answer and is sent as the whole request body.
func (c *Client) SubmitRfi(ctx context.Context, customerID string, response rfi.Response, opts ...option.RequestOption) (*rfi.SubmitResponse, error) {
	if customerID == "" {
		return nil, fmt.Errorf("customer ID cannot be empty")
	}
	if response == nil {
		response = rfi.Response{}
	}

	path := request.Pathf("/instances/%s/customers/%s/rfi", c.instanceID, customerID)
	return request.Do[*rfi.SubmitResponse](c.cfg, ctx, "POST", path, response, opts...)
}
