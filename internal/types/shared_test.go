package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestPaginationMetadata_Unmarshal pins the wire shape reported by the
// spec (packages/api-contract/src/generic.schema.ts: next_page/prev_page
// are `z.string().nullish()`, example "pi_123"). NextPage/PrevPage were
// bare int before this fix, which could never decode a real cursor value:
// unmarshalling {"has_more":true,"next_page":"pi_123","prev_page":null}
// failed with "json: cannot unmarshal string into Go struct field
// PaginationMetadata.next_page of type int" for every paginated list
// response that actually had a next page.
func TestPaginationMetadata_Unmarshal(t *testing.T) {
	tests := []struct {
		name         string
		wire         string
		wantHasMore  bool
		wantNextPage *string
		wantPrevPage *string
	}{
		{
			name:         "a real cursor on next_page, null prev_page (the exact reported failure payload)",
			wire:         `{"has_more":true,"next_page":"pi_123","prev_page":null}`,
			wantHasMore:  true,
			wantNextPage: strPtr("pi_123"),
			wantPrevPage: nil,
		},
		{
			name:         "both cursors null (first page, no more results)",
			wire:         `{"has_more":false,"next_page":null,"prev_page":null}`,
			wantHasMore:  false,
			wantNextPage: nil,
			wantPrevPage: nil,
		},
		{
			name:         "a cursor on both next_page and prev_page (a middle page)",
			wire:         `{"has_more":true,"next_page":"pi_456","prev_page":"pi_123"}`,
			wantHasMore:  true,
			wantNextPage: strPtr("pi_456"),
			wantPrevPage: strPtr("pi_123"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got PaginationMetadata
			err := json.Unmarshal([]byte(tt.wire), &got)
			require.NoError(t, err)
			require.Equal(t, tt.wantHasMore, got.HasMore)
			if tt.wantNextPage == nil {
				require.Nil(t, got.NextPage)
			} else {
				require.NotNil(t, got.NextPage)
				require.Equal(t, *tt.wantNextPage, *got.NextPage)
			}
			if tt.wantPrevPage == nil {
				require.Nil(t, got.PrevPage)
			} else {
				require.NotNil(t, got.PrevPage)
				require.Equal(t, *tt.wantPrevPage, *got.PrevPage)
			}
		})
	}
}

func strPtr(s string) *string { return &s }
