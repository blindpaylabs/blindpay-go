package pagination

import (
	"context"
	"errors"
	"testing"

	"github.com/blindpaylabs/blindpay-go/internal/types"
	"github.com/stretchr/testify/require"
)

func strPtr(s string) *string { return &s }

type fakePages struct {
	pages   map[string]*Page[string]
	cursors []string
	failAt  string
}

func (f *fakePages) fetch(_ context.Context, cursor string) (*Page[string], error) {
	f.cursors = append(f.cursors, cursor)
	if cursor == f.failAt && f.failAt != "" {
		return nil, errors.New("boom")
	}
	return f.pages[cursor], nil
}

func page(items []string, next string) *Page[string] {
	p := &Page[string]{Data: items}
	if next != "" {
		p.Pagination = types.PaginationMetadata{HasMore: true, NextPage: strPtr(next)}
	}
	return p
}

func collect(p *Pager[string]) []string {
	var out []string
	for p.Next() {
		out = append(out, p.Current())
	}
	return out
}

func TestPager(t *testing.T) {
	tests := []struct {
		name        string
		start       string
		pages       map[string]*Page[string]
		failAt      string
		want        []string
		wantCursors []string
		wantErr     string
	}{
		{
			name:        "follows next_page while has_more",
			pages:       map[string]*Page[string]{"": page([]string{"a", "b"}, "b"), "b": page([]string{"c"}, "c"), "c": page([]string{"d"}, "")},
			want:        []string{"a", "b", "c", "d"},
			wantCursors: []string{"", "b", "c"},
		},
		{
			name:        "single page",
			pages:       map[string]*Page[string]{"": page([]string{"a"}, "")},
			want:        []string{"a"},
			wantCursors: []string{""},
		},
		{
			name:        "empty list",
			pages:       map[string]*Page[string]{"": page(nil, "")},
			wantCursors: []string{""},
		},
		{
			name:        "empty page with has_more continues",
			pages:       map[string]*Page[string]{"": page(nil, "x"), "x": page([]string{"a"}, "")},
			want:        []string{"a"},
			wantCursors: []string{"", "x"},
		},
		{
			name:        "starts after the caller's cursor",
			start:       "z",
			pages:       map[string]*Page[string]{"z": page([]string{"a"}, "")},
			want:        []string{"a"},
			wantCursors: []string{"z"},
		},
		{
			name:        "error stops iteration after yielded items",
			pages:       map[string]*Page[string]{"": page([]string{"a"}, "b")},
			failAt:      "b",
			want:        []string{"a"},
			wantCursors: []string{"", "b"},
			wantErr:     "boom",
		},
		{
			name:        "cursor that does not advance is an error",
			pages:       map[string]*Page[string]{"": page([]string{"a"}, "a"), "a": page([]string{"b"}, "a")},
			want:        []string{"a", "b"},
			wantCursors: []string{"", "a"},
			wantErr:     "did not advance",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			f := &fakePages{pages: tc.pages, failAt: tc.failAt}
			p := NewPager[string](context.Background(), tc.start, f.fetch)

			require.Equal(t, tc.want, collect(p))
			require.Equal(t, tc.wantCursors, f.cursors)
			if tc.wantErr == "" {
				require.NoError(t, p.Err())
			} else {
				require.ErrorContains(t, p.Err(), tc.wantErr)
			}
			require.False(t, p.Next(), "Next stays false after the end")
		})
	}
}

func TestNewErrPager(t *testing.T) {
	p := NewErrPager[string](ErrEndingBeforeUnsupported)
	require.False(t, p.Next())
	require.ErrorIs(t, p.Err(), ErrEndingBeforeUnsupported)
}

func TestPage_NextCursor(t *testing.T) {
	require.Equal(t, "", (&Page[string]{}).NextCursor())
	require.Equal(t, "", (&Page[string]{Pagination: types.PaginationMetadata{HasMore: false, NextPage: strPtr("x")}}).NextCursor())
	p := page(nil, "x")
	require.True(t, p.HasMore())
	require.Equal(t, "x", p.NextCursor())
}
