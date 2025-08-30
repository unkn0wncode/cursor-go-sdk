package cursor

import (
	"fmt"
)

// APIError represents a non-2xx HTTP response from the Cursor API.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: status=%d body=%s", e.StatusCode, e.Body)
}
