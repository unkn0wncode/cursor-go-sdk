package cursor

import (
	"context"
)

// ListModels retrieves available model names.
func (c *Client) ListModels(ctx context.Context) (*ListModelsResponse, error) {
	var out ListModelsResponse
	if err := c.do(ctx, "GET", "/v0/models", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
