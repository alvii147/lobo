package flux

import "net/http"

// ResponseWriterMiddleware converts a HandlerFunc to an http.Handler.
// This should be the top-level middleware when setting up routes.
func ResponseWriterMiddleware(next HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)
	})
}
