package flux_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvii147/lobo/flux"
	"github.com/alvii147/lobo/flux/mocks"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestLoggerMiddleware(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		statusCode int
		method     string
		wantLevel  string
	}{
		"200 status code with GET request causes info level log": {
			statusCode: http.StatusOK,
			method:     http.MethodGet,
			wantLevel:  "Info",
		},
		"200 status code with POST request causes info level log": {
			statusCode: http.StatusOK,
			method:     http.MethodPost,
			wantLevel:  "Info",
		},
		"201 status code causes info level log": {
			statusCode: http.StatusCreated,
			method:     http.MethodGet,
			wantLevel:  "Info",
		},
		"302 status code causes info level log": {
			statusCode: http.StatusFound,
			method:     http.MethodGet,
			wantLevel:  "Info",
		},
		"400 status code causes warn level log": {
			statusCode: http.StatusBadRequest,
			method:     http.MethodGet,
			wantLevel:  "Warn",
		},
		"404 status code causes warn level log": {
			statusCode: http.StatusNotFound,
			method:     http.MethodGet,
			wantLevel:  "Warn",
		},
		"500 status code causes error level log": {
			statusCode: http.StatusInternalServerError,
			method:     http.MethodGet,
			wantLevel:  "Error",
		},
		"501 status code causes error level log": {
			statusCode: http.StatusNotImplemented,
			method:     http.MethodGet,
			wantLevel:  "Error",
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			nextCallCount := 0
			var next flux.HandlerFunc = func(w *flux.ResponseWriter, r *http.Request) {
				nextCallCount++
			}

			url := "/request/url/path"
			rec := httptest.NewRecorder()
			w := &flux.ResponseWriter{
				ResponseWriter: rec,
				StatusCode:     testcase.statusCode,
			}
			r := httptest.NewRequest(testcase.method, url, http.NoBody)

			ctrl := gomock.NewController(t)
			logger := mocks.NewMocklogger(ctrl)
			var mockLogFunc func(msg any, args ...any) *gomock.Call

			switch testcase.wantLevel {
			case "Info":
				mockLogFunc = logger.EXPECT().Info
			case "Warn":
				mockLogFunc = logger.EXPECT().Warn
			case "Error":
				mockLogFunc = logger.EXPECT().Error
			default:
				t.Fatal("Unknown log level")
			}

			mockLogFunc(
				"HTTP Traffic",
				[]any{
					"method",
					testcase.method,
					"url",
					url,
					"protocol",
					"HTTP/1.1",
					"status",
					testcase.statusCode,
					"status_text",
					http.StatusText(testcase.statusCode),
				},
			)

			flux.NewLoggerMiddleware(logger)(next)(w, r)
			require.Equal(t, 1, nextCallCount)
		})
	}
}
