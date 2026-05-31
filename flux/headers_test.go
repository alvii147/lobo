package flux_test

import (
	"net/http"
	"testing"

	"github.com/alvii147/lobo/flux"
	"github.com/stretchr/testify/require"
)

func TestGetAuthorizationHeader(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		header    http.Header
		authType  string
		wantToken string
		wantOk    bool
	}{
		"Valid header with valid auth type": {
			header: map[string][]string{
				"Authorization": {"Bearer 0xdeadbeef"},
			},
			authType:  "Bearer",
			wantToken: "0xdeadbeef",
			wantOk:    true,
		},
		"No header": {
			header:    make(map[string][]string),
			authType:  "Bearer",
			wantToken: "0xdeadbeef",
			wantOk:    false,
		},
		"Invalid auth type": {
			header: map[string][]string{
				"Authorization": {"Bearer 0xdeadbeef"},
			},
			authType:  "Basic",
			wantToken: "0xdeadbeef",
			wantOk:    false,
		},
		"Valid header with spaces": {
			header: map[string][]string{
				"Authorization": {"  Bearer   0xdeadbeef    "},
			},
			authType:  "Bearer",
			wantToken: "0xdeadbeef",
			wantOk:    true,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			token, ok := flux.GetAuthorizationHeader(testcase.header, testcase.authType)
			require.Equal(t, testcase.wantOk, ok)
			if testcase.wantOk {
				require.Equal(t, testcase.wantToken, token)
			}
		})
	}
}
