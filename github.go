package cursor

import (
	"context"
)

// ListRepositories returns GitHub repositories available to the authenticated user.
// Matches OpenAPI /v0/repositories without query parameters.
func (c *Client) ListRepositories(ctx context.Context) (*ListRepositoriesResponse, error) {
	var out ListRepositoriesResponse
	if err := c.do(ctx, "GET", "/v0/repositories", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
