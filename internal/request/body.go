package request

import (
	"encoding/json"
	"fmt"
)

// BodyWith encodes params (its json tags decide what is sent, so a field
// added to the struct reaches the wire) and adds the extra members, e.g. a
// "type" discriminator the method injects. Values keep their exact JSON
// encoding.
func BodyWith(params any, extra map[string]any) (map[string]json.RawMessage, error) {
	b, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	body := map[string]json.RawMessage{}
	if err := json.Unmarshal(b, &body); err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}
	for k, v := range extra {
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		body[k] = raw
	}
	return body, nil
}
