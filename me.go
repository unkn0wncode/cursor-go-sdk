package cursor

import (
	"context"
)

// Me returns metadata about the current API key.
func (c *Client) Me(ctx context.Context) (*MeResponse, error) {
	var out MeResponse
	if err := c.do(ctx, "GET", "/v0/me", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
