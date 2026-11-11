package request

import (
	"fmt"
	"net/url"
)

// Pathf builds a request path from a format with one %s per path segment,
// percent-encoding every argument with url.PathEscape so an ID such as
// "a/b" or "../x" cannot change the route ("a/b" -> "a%2Fb").
func Pathf(format string, segments ...string) string {
	args := make([]any, len(segments))
	for i, s := range segments {
		args[i] = url.PathEscape(s)
	}
	return fmt.Sprintf(format, args...)
}
