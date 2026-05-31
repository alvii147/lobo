package flux

// ListResponse represents the generic response body for paginated list endpoints.
type ListResponse[T any] struct {
	Data       []T        `json:"data"`
	Pagination Pagination `json:"pagination"`
}
