package flux

import (
	"net/http"
	"strconv"
)

const (
	// PaginationPageQueryParamKey is the query param key with page data.
	PaginationPageQueryParamKey = "page"
	// PaginationPageSizeQueryParamKey is the query param key with page size data.
	PaginationPageSizeQueryParamKey = "page_size"
)

// Pagination represents pagination metadata.
type Pagination struct {
	Page     int64 `json:"page"`
	PageSize int64 `json:"page_size"`
	Total    int64 `json:"total"`
}

// GetPaginationPageParam extracts the page query parameter of a request.
// If no query parameter is found, the default page is returned.
func GetPaginationPageParam(r *http.Request, defaultPage int64) int64 {
	param := r.URL.Query().Get(PaginationPageQueryParamKey)
	if param == "" {
		return defaultPage
	}

	page, err := strconv.ParseInt(param, 10, 0)
	if err != nil {
		return defaultPage
	}

	return page
}

// GetPaginationPageSizeParam extracts the page size query parameter of a request.
// If no query parameter is found, the default page size is returned.
func GetPaginationPageSizeParam(r *http.Request, defaultPageSize int64) int64 {
	param := r.URL.Query().Get(PaginationPageSizeQueryParamKey)
	if param == "" {
		return defaultPageSize
	}

	pageSize, err := strconv.ParseInt(param, 10, 0)
	if err != nil {
		return defaultPageSize
	}

	return pageSize
}

// GetPaginationParams extracts the page and page size query parameters of a request.
// If no query parameters are found, the default page and page size are returned.
func GetPaginationParams(r *http.Request, defaultPage int64, defaultPageSize int64) (int64, int64) {
	page := GetPaginationPageParam(r, defaultPage)
	pageSize := GetPaginationPageSizeParam(r, defaultPageSize)

	return page, pageSize
}
