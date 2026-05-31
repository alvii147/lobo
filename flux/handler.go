package flux

import "net/http"

// HandlerFunc takes in a request and response writer and implements a handler.
type HandlerFunc func(w *ResponseWriter, r *http.Request)

// ServeHTTP calls the handler function.
func (f HandlerFunc) ServeHTTP(w *ResponseWriter, r *http.Request) {
	f(w, r)
}
