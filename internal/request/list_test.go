package request

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

type listItem struct {
	ID string `json:"id"`
}

func TestList_UnmarshalJSON(t *testing.T) {
	next := "it_2"
	tests := []struct {
		name       string
		body       string
		wantIDs    []string
		pagination bool
		hasMore    bool
		nextPage   *string
		wantErr    bool
	}{
		{name: "bare array", body: `[{"id":"it_1"},{"id":"it_2"}]`, wantIDs: []string{"it_1", "it_2"}},
		{name: "bare array with leading whitespace", body: "\n  [{\"id\":\"it_1\"}]", wantIDs: []string{"it_1"}},
		{name: "empty array", body: `[]`, wantIDs: []string{}},
		{
			name:       "envelope",
			body:       `{"data":[{"id":"it_1"},{"id":"it_2"}],"pagination":{"has_more":true,"next_page":"it_2","prev_page":null}}`,
			wantIDs:    []string{"it_1", "it_2"},
			pagination: true,
			hasMore:    true,
			nextPage:   &next,
		},
		{
			name:       "envelope last page",
			body:       `{"data":[],"pagination":{"has_more":false,"next_page":null,"prev_page":"it_0"}}`,
			wantIDs:    []string{},
			pagination: true,
		},
		{name: "null", body: `null`},
		{name: "empty body", body: ``, wantErr: true},
		{name: "string", body: `"x"`, wantErr: true},
		{name: "array of wrong type", body: `[1]`, wantErr: true},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var got List[listItem]
			err := got.UnmarshalJSON([]byte(tc.body))
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			var ids []string
			if got.Data != nil {
				ids = []string{}
			}
			for _, it := range got.Data {
				ids = append(ids, it.ID)
			}
			require.Equal(t, tc.wantIDs, ids)

			if !tc.pagination {
				require.Nil(t, got.Pagination)
				return
			}
			require.NotNil(t, got.Pagination)
			require.Equal(t, tc.hasMore, got.Pagination.HasMore)
			require.Equal(t, tc.nextPage, got.Pagination.NextPage)
		})
	}
}

func TestList_ViaJSONUnmarshal(t *testing.T) {
	var got List[listItem]
	require.NoError(t, json.Unmarshal([]byte(`{"data":[{"id":"it_1"}],"pagination":{"has_more":false}}`), &got))
	require.Len(t, got.Data, 1)
	require.NotNil(t, got.Pagination)
}
