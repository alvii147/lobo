package flux

// ErrorResponse represents the general error response body.
type ErrorResponse struct {
	Code               string              `json:"code"`
	Detail             string              `json:"detail"`
	ValidationFailures map[string][]string `json:"failures,omitempty"`
}
