package request

import (
	"net/url"
	"strconv"
)

// Query builds a query string, skipping unset values so nothing like
// "limit=0" or "status=" is ever sent.
type Query struct {
	values url.Values
}

// String adds key=value when value is non-empty.
func (q *Query) String(key, value string) {
	if value == "" {
		return
	}
	if q.values == nil {
		q.values = url.Values{}
	}
	q.values.Set(key, value)
}

// Int adds key=n (as a decimal string) when n > 0.
func (q *Query) Int(key string, n int) {
	if n > 0 {
		q.String(key, strconv.Itoa(n))
	}
}

// Apply appends the encoded query to path, if any.
func (q *Query) Apply(path string) string {
	if len(q.values) == 0 {
		return path
	}
	return path + "?" + q.values.Encode()
}
