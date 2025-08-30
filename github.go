package cursor

import (
	"context"
)

// ListRepositories returns GitHub repositories available to the authenticated user.
// This endpoint has very strict rate limits.
// Limit requests to 1 / user / minute, and 30 / user / hour.
// This request can take tens of seconds to respond for users with access to many repositories.
// Make sure to handle this information not being available gracefully.
func (c *Client) ListRepositories(ctx context.Context) (*ListRepositoriesResponse, error) {
	var out ListRepositoriesResponse
	if err := c.do(ctx, "GET", "/v0/repositories", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
