package flux

import (
	"net/http"
	"time"
)

// HTTPClientDefaultTimeout is the default timeout for HTTP clients.
const HTTPClientDefaultTimeout = 60 * time.Second

// NewHTTPClient creates and returns a new HTTP client.
func NewHTTPClient(modifier func(c *http.Client)) *http.Client {
	client := &http.Client{
		Timeout: HTTPClientDefaultTimeout,
	}

	if modifier != nil {
		modifier(client)
	}

	return client
}
