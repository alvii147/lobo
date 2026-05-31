package flux

import (
	"encoding/json"
	"net/http"

	"github.com/alvii147/lobo/errfmt"
)

// ResponseWriter stores an http.ResponseWriter and the HTTP status code.
// This is used for retaining the status code after the handler is executed.
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// NewResponseWriter returns a new ResponseWriter.
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		StatusCode:     http.StatusOK,
	}
}

// Header returns the headers in ResponseWriter.
func (w *ResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

// Write writes bytes data to ResponseWriter.
func (w *ResponseWriter) Write(p []byte) (int, error) {
	n, err := w.ResponseWriter.Write(p)
	if err != nil {
		return 0, errfmt.Errorf(err, "ResponseWriter.Write failed")
	}

	return n, nil
}

// WriteHeader writes status code to ResponseWriter.
func (w *ResponseWriter) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// WriteText writes status code and JSON data to ResponseWriter.
func (w *ResponseWriter) WriteText(text string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, err := w.Write([]byte(text))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// WriteJSON writes status code and JSON data to ResponseWriter.
func (w *ResponseWriter) WriteJSON(data any, statusCode int) {
	w.WriteHeader(statusCode)
	if data == nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
