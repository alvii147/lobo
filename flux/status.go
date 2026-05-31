package flux

import "net/http"

// IsHTTPSuccess determines whether or not a given status code is 2xx.
func IsHTTPSuccess(statusCode int) bool {
	return statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices
}
