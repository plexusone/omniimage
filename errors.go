package omniimage

import "errors"

// Errors returned by omniimage.
var (
	// ErrNoProviders is returned when no providers are configured.
	ErrNoProviders = errors.New("no providers configured")

	// ErrUnknownProvider is returned for an unknown provider name.
	ErrUnknownProvider = errors.New("unknown provider")

	// ErrNotSupported is returned when an operation is not supported.
	ErrNotSupported = errors.New("operation not supported")

	// ErrInvalidRequest is returned for invalid request parameters.
	ErrInvalidRequest = errors.New("invalid request")

	// ErrRateLimited is returned when the API rate limit is exceeded.
	ErrRateLimited = errors.New("rate limited")

	// ErrContentPolicy is returned when content violates the provider's policy.
	ErrContentPolicy = errors.New("content policy violation")

	// ErrModelNotFound is returned when the requested model is not available.
	ErrModelNotFound = errors.New("model not found")
)

// APIError represents an error from the image generation API.
type APIError struct {
	// StatusCode is the HTTP status code.
	StatusCode int

	// Code is the provider-specific error code.
	Code string

	// Message is the error message.
	Message string

	// Provider is the provider that returned the error.
	Provider string
}

func (e *APIError) Error() string {
	if e.Code != "" {
		return e.Provider + ": " + e.Code + ": " + e.Message
	}
	return e.Provider + ": " + e.Message
}
