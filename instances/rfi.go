package instances

import (
	"context"

	"github.com/blindpaylabs/blindpay-go/internal/request"
	"github.com/blindpaylabs/blindpay-go/option"
	"github.com/blindpaylabs/blindpay-go/rfi"
)

// GetRfi retrieves the instance's open request for information. It returns
// nil, nil when no RFI is open.
func (c *Client) GetRfi(ctx context.Context, opts ...option.RequestOption) (*rfi.InstanceRfi, error) {
	path := request.Pathf("/instances/%s/rfi", c.instanceID)
	return request.Do[*rfi.InstanceRfi](c.cfg, ctx, "GET", path, nil, opts...)
}

// SubmitRfi answers the instance's open RFI. response maps each field key to
// its answer and is sent as the whole request body.
func (c *Client) SubmitRfi(ctx context.Context, response rfi.Response, opts ...option.RequestOption) (*rfi.SubmitResponse, error) {
	if response == nil {
		response = rfi.Response{}
	}
	path := request.Pathf("/instances/%s/rfi", c.instanceID)
	return request.Do[*rfi.SubmitResponse](c.cfg, ctx, "POST", path, response, opts...)
}
