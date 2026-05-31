package flux

// MiddlewareFunc takes in a HandlerFunc and returns a HandlerFunc.
// Middleware is used to add a processing step to handlers.
type MiddlewareFunc func(HandlerFunc) HandlerFunc
