package cursor

import (
	"context"
	"fmt"
	"net/url"
)

// LaunchAgent starts a new background agent.
func (c *Client) LaunchAgent(ctx context.Context, req LaunchRequest) (*Agent, error) {
	var out Agent
	if err := c.do(ctx, "POST", "/v0/agents", nil, req, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// AddFollowup sends additional instructions to a running agent.
func (c *Client) AddFollowup(ctx context.Context, id string, req FollowupRequest) (string, error) {
	var out FollowupResponse
	path := fmt.Sprintf("/v0/agents/%s/followup", url.PathEscape(id))
	if err := c.do(ctx, "POST", path, nil, req, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

// GetAgent retrieves the current status of an agent.
func (c *Client) GetAgent(ctx context.Context, id string) (*Agent, error) {
	var out Agent
	path := fmt.Sprintf("/v0/agents/%s", url.PathEscape(id))
	if err := c.do(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListAgents retrieves multiple agents with optional pagination.
func (c *Client) ListAgents(ctx context.Context, limit int, cursor *string) (*ListAgentsResponse, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != nil && *cursor != "" {
		q.Set("cursor", *cursor)
	}
	var out ListAgentsResponse
	if err := c.do(ctx, "GET", "/v0/agents", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// DeleteAgent terminates and deletes an agent.
func (c *Client) DeleteAgent(ctx context.Context, id string) (string, error) {
	var out DeleteResponse
	path := fmt.Sprintf("/v0/agents/%s", url.PathEscape(id))
	if err := c.do(ctx, "DELETE", path, nil, nil, &out); err != nil {
		return "", err
	}
	return out.ID, nil
}

// GetConversation returns the conversation history for an agent.
func (c *Client) GetConversation(ctx context.Context, id string) (*Conversation, error) {
	var out Conversation
	path := fmt.Sprintf("/v0/agents/%s/conversation", url.PathEscape(id))
	if err := c.do(ctx, "GET", path, nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
