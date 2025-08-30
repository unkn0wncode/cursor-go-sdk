package cursor

import (
	"fmt"
)

// APIError represents a non-2xx HTTP response from the Cursor API.
// It attempts to map the OpenAPI error shape: { "error": { "message": string, "code": string } }.
type APIError struct {
	StatusCode int
	Message    string
	Code       string
	Body       string
}

func (e *APIError) Error() string {
	if e.Message != "" || e.Code != "" {
		if e.Code != "" {
			return fmt.Sprintf("API error: status=%d code=%s message=%s", e.StatusCode, e.Code, e.Message)
		}
		return fmt.Sprintf("API error: status=%d message=%s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("API error: status=%d body=%s", e.StatusCode, e.Body)
}
