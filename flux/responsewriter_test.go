package flux_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alvii147/lobo/flux"
	"github.com/stretchr/testify/require"
)

type mockResponseWriter struct {
	headers      map[string][]string
	writtenBytes []byte
	writeErr     error
	statusCode   int
}

func (w *mockResponseWriter) Header() http.Header {
	return w.headers
}

func (w *mockResponseWriter) Write(p []byte) (int, error) {
	w.writtenBytes = append(w.writtenBytes, p...)

	return len(p), w.writeErr
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

func TestResponseWriterHeader(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	w := flux.NewResponseWriter(rec)

	w.Header().Set("Content-Type", "application/json")
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestResponseWriterWriteSuccess(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	w := flux.NewResponseWriter(rec)

	data := "DEADBEEF"

	_, err := w.Write([]byte(data))
	require.NoError(t, err)
	require.Equal(t, data, rec.Body.String())
}

func TestResponseWriterWriteError(t *testing.T) {
	t.Parallel()

	writeErr := errors.New("Write failed")
	w := flux.NewResponseWriter(&mockResponseWriter{
		headers:  make(map[string][]string),
		writeErr: writeErr,
	})

	_, err := w.Write([]byte("DEADBEEF"))
	require.ErrorIs(t, err, writeErr)
}

func TestResponseWriterWriteHeader(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		statusCode int
	}{
		"Status code OK": {
			statusCode: http.StatusOK,
		},
		"Status code created": {
			statusCode: http.StatusCreated,
		},
		"Status code moved permanently": {
			statusCode: http.StatusMovedPermanently,
		},
		"Status code found": {
			statusCode: http.StatusFound,
		},
		"Status code not found": {
			statusCode: http.StatusNotFound,
		},
		"Status code internal server error": {
			statusCode: http.StatusInternalServerError,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			w := flux.NewResponseWriter(rec)

			w.WriteHeader(testcase.statusCode)
			require.Equal(t, testcase.statusCode, w.StatusCode)
			require.Equal(t, testcase.statusCode, rec.Code)
		})
	}
}

func TestResponseWriterWriteTextSuccess(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	w := flux.NewResponseWriter(rec)

	text := "hello"

	w.WriteText(text, http.StatusOK)
	require.Equal(t, http.StatusOK, rec.Code)

	b, err := io.ReadAll(rec.Body)
	require.NoError(t, err)
	require.Equal(t, text, string(b))
	require.Equal(t, "text/plain; charset=utf-8", w.Header().Get("Content-Type"))
}

func TestResponseWriterWriteTextError(t *testing.T) {
	t.Parallel()

	w := flux.NewResponseWriter(&mockResponseWriter{
		headers:  make(map[string][]string),
		writeErr: errors.New("Write failed"),
	})

	w.WriteText("hello", http.StatusOK)
	require.Equal(t, http.StatusInternalServerError, w.StatusCode)
}

func TestResponseWriterWriteJSONWithoutData(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	w := flux.NewResponseWriter(rec)

	w.WriteJSON(nil, http.StatusOK)
	require.Equal(t, http.StatusOK, rec.Code)
}

func TestResponseWriterWriteJSONWithData(t *testing.T) {
	t.Parallel()

	rec := httptest.NewRecorder()
	w := flux.NewResponseWriter(rec)

	data := map[string]any{
		"number": float64(42),
		"string": "Hello",
		"null":   nil,
		"listOfNumbers": []any{
			float64(3),
			float64(1),
			float64(4),
			float64(1),
			float64(6),
		},
	}

	w.WriteJSON(data, http.StatusOK)
	require.Equal(t, http.StatusOK, rec.Code)

	writtenData := make(map[string]any)
	err := json.NewDecoder(rec.Body).Decode(&writtenData)
	require.NoError(t, err)
	require.Equal(t, data["number"], writtenData["number"])
	require.Equal(t, data["string"], writtenData["string"])
	require.Equal(t, data["null"], writtenData["null"])
	require.Equal(t, data["listOfNumbers"], writtenData["listOfNumbers"])
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

func TestResponseWriterWriteJSONError(t *testing.T) {
	t.Parallel()

	w := flux.NewResponseWriter(&mockResponseWriter{
		headers:  make(map[string][]string),
		writeErr: errors.New("Write failed"),
	})

	w.WriteJSON(map[string]any{}, http.StatusOK)
	require.Equal(t, http.StatusInternalServerError, w.StatusCode)
}
