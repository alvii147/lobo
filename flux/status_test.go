package flux_test

import (
	"net/http"
	"testing"

	"github.com/alvii147/lobo/flux"
	"github.com/stretchr/testify/require"
)

func TestIsHTTPSuccess(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		statusCode  int
		wantSuccess bool
	}{
		"200 OK is success": {
			statusCode:  http.StatusOK,
			wantSuccess: true,
		},
		"201 Created is success": {
			statusCode:  http.StatusCreated,
			wantSuccess: true,
		},
		"204 No content is success": {
			statusCode:  http.StatusNoContent,
			wantSuccess: true,
		},
		"302 Found is not success": {
			statusCode:  http.StatusFound,
			wantSuccess: false,
		},
		"400 Bad request is not success": {
			statusCode:  http.StatusBadRequest,
			wantSuccess: false,
		},
		"401 Unauthorized is not success": {
			statusCode:  http.StatusUnauthorized,
			wantSuccess: false,
		},
		"403 Forbidden is not success": {
			statusCode:  http.StatusForbidden,
			wantSuccess: false,
		},
		"404 Not found is not success": {
			statusCode:  http.StatusNotFound,
			wantSuccess: false,
		},
		"405 Method not allowed is not success": {
			statusCode:  http.StatusMethodNotAllowed,
			wantSuccess: false,
		},
		"500 Internal server error is not success": {
			statusCode:  http.StatusInternalServerError,
			wantSuccess: false,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, testcase.wantSuccess, flux.IsHTTPSuccess(testcase.statusCode))
		})
	}
}
