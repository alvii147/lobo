package flux_test

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/alvii147/lobo/flux"
	"github.com/stretchr/testify/require"
)

func TestGetPaginationPageParam(t *testing.T) {
	t.Parallel()

	var defaultPage int64 = 5
	testcases := map[string]struct {
		rawQuery string
		wantPage int64
	}{
		"Valid page": {
			rawQuery: "page=42",
			wantPage: 42,
		},
		"No page": {
			rawQuery: "key=value",
			wantPage: defaultPage,
		},
		"Invalid page": {
			rawQuery: "page=deadbeef",
			wantPage: defaultPage,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "localhost:8080",
					Path:     "/path",
					RawQuery: testcase.rawQuery,
				},
			}

			page := flux.GetPaginationPageParam(req, defaultPage)
			require.Equal(t, testcase.wantPage, page)
		})
	}
}

func TestGetPaginationPageSizeParam(t *testing.T) {
	t.Parallel()

	var defaultPageSize int64 = 50
	testcases := map[string]struct {
		rawQuery     string
		wantPageSize int64
	}{
		"Valid page size": {
			rawQuery:     "page_size=42",
			wantPageSize: 42,
		},
		"No page size": {
			rawQuery:     "key=value",
			wantPageSize: defaultPageSize,
		},
		"Invalid page size": {
			rawQuery:     "page_size=deadbeef",
			wantPageSize: defaultPageSize,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "localhost:8080",
					Path:     "/auth/api-keys",
					RawQuery: testcase.rawQuery,
				},
			}

			pageSize := flux.GetPaginationPageSizeParam(req, defaultPageSize)
			require.Equal(t, testcase.wantPageSize, pageSize)
		})
	}
}

func TestGetPaginationParams(t *testing.T) {
	t.Parallel()

	var defaultPage int64 = 5
	var defaultPageSize int64 = 50

	testcases := map[string]struct {
		rawQuery     string
		wantPage     int64
		wantPageSize int64
	}{
		"Valid page & page size": {
			rawQuery:     "page=314&page_size=42",
			wantPage:     314,
			wantPageSize: 42,
		},
		"No page": {
			rawQuery:     "page_size=42",
			wantPage:     defaultPage,
			wantPageSize: 42,
		},
		"No page size": {
			rawQuery:     "page=314",
			wantPage:     314,
			wantPageSize: defaultPageSize,
		},
		"Invalid page": {
			rawQuery:     "page=deadbeef&page_size=42",
			wantPage:     defaultPage,
			wantPageSize: 42,
		},
		"Invalid page size": {
			rawQuery:     "page=314&page_size=deadbeef",
			wantPage:     314,
			wantPageSize: defaultPageSize,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			req := &http.Request{
				URL: &url.URL{
					Scheme:   "https",
					Host:     "localhost:8080",
					Path:     "/path",
					RawQuery: testcase.rawQuery,
				},
			}

			page, pageSize := flux.GetPaginationParams(req, defaultPage, defaultPageSize)
			require.Equal(t, testcase.wantPage, page)
			require.Equal(t, testcase.wantPageSize, pageSize)
		})
	}
}
