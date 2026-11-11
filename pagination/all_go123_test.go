//go:build go1.23

package pagination

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPager_All(t *testing.T) {
	f := &fakePages{pages: map[string]*Page[string]{"": page([]string{"a"}, "b")}, failAt: "b"}
	var got []string
	var gotErr error
	for item, err := range NewPager[string](context.Background(), "", f.fetch).All() {
		if err != nil {
			gotErr = err
			break
		}
		got = append(got, item)
	}
	require.Equal(t, []string{"a"}, got)
	require.ErrorContains(t, gotErr, "boom")

	// Breaking early stops fetching.
	f = &fakePages{pages: map[string]*Page[string]{"": page([]string{"a", "b"}, "b"), "b": page([]string{"c"}, "")}}
	for range NewPager[string](context.Background(), "", f.fetch).All() {
		break
	}
	require.Equal(t, []string{""}, f.cursors)
}
