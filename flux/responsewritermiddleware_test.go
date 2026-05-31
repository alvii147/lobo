package flux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvii147/lobo/flux"
	"github.com/stretchr/testify/require"
)

func TestResponseWriterMiddleware(t *testing.T) {
	t.Parallel()

	nextCallCount := 0
	var next flux.HandlerFunc = func(w *flux.ResponseWriter, r *http.Request) {
		nextCallCount++
	}

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/request/url/path", http.NoBody)

	flux.ResponseWriterMiddleware(next).ServeHTTP(rec, r)
	require.Equal(t, 1, nextCallCount)
}
