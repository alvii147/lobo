package flux

import (
	"net/http"
)

// logger logs at debug, info, warn, and error levels.
//
//go:generate mockgen -package=mocks -source=$GOFILE -destination=./mocks/logger.go
type logger interface {
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

// logFunc is a function that logs a given message.
type logFunc func(msg string, args ...any)

// logRequest logs a single request given a logFunc.
func logRequest(lf logFunc, w *ResponseWriter, r *http.Request) {
	lf(
		"HTTP Traffic",
		"method",
		r.Method,
		"url",
		r.URL.String(),
		"protocol",
		r.Proto,
		"status",
		w.StatusCode,
		"status_text",
		http.StatusText(w.StatusCode),
	)
}

// logTraffic logs HTTP traffic, including http method, URL, protocol, and status code.
func logTraffic(l logger, w *ResponseWriter, r *http.Request) {
	switch {
	case w.StatusCode < http.StatusBadRequest:
		logRequest(l.Info, w, r)
	case w.StatusCode < http.StatusInternalServerError:
		logRequest(l.Warn, w, r)
	default:
		logRequest(l.Error, w, r)
	}
}

// NewLoggerMiddleware creates a middleware using LoggerMiddleware.
func NewLoggerMiddleware(lg logger) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return HandlerFunc(func(w *ResponseWriter, r *http.Request) {
			defer func() {
				logTraffic(lg, w, r)
			}()

			next.ServeHTTP(w, r)
		})
	}
}
