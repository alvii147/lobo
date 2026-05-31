package flux_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/alvii147/lobo/flux"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPClient(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		modifier    func(c *http.Client)
		wantTimeout time.Duration
	}{
		"No modifier": {
			modifier:    nil,
			wantTimeout: flux.HTTPClientDefaultTimeout,
		},
		"Empty modifier": {
			modifier:    func(c *http.Client) {},
			wantTimeout: flux.HTTPClientDefaultTimeout,
		},
		"Timeout modifier": {
			modifier: func(c *http.Client) {
				c.Timeout = 5 * time.Second
			},
			wantTimeout: 5 * time.Second,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			httpClient := flux.NewHTTPClient(testcase.modifier)
			require.Equal(t, testcase.wantTimeout, httpClient.Timeout)
		})
	}
}
