package flux

import (
	"net/http"
	"strings"
)

const (
	// HTTPHeaderAccessControlAllowHeaders is the CORS header that lists the allowed header names.
	HTTPHeaderAccessControlAllowHeaders = "Access-Control-Allow-Headers"
	// HTTPHeaderAccessControlAllowMethods is the CORS header that lists the allowed methods.
	HTTPHeaderAccessControlAllowMethods = "Access-Control-Allow-Methods"
	// HTTPHeaderAccessControlAllowOrigin is the CORS header that specifies the allowed origin URL.
	HTTPHeaderAccessControlAllowOrigin = "Access-Control-Allow-Origin"
	// HTTPHeaderContentType is the header that defines the content type of the body.
	HTTPHeaderContentType = "Content-Type"
	// HTTPHeaderAuthorization is the header used for authentication/authorization credentials.
	HTTPHeaderAuthorization = "Authorization"
)

// GetAuthorizationHeader parses HTTP authorization header.
func GetAuthorizationHeader(header http.Header, authType string) (string, bool) {
	token, ok := strings.CutPrefix(strings.TrimSpace(header.Get(HTTPHeaderAuthorization)), authType)
	if !ok {
		return "", false
	}

	return strings.TrimSpace(token), true
}
