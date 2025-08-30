package cursor

import (
	"time"
)

// Possible values for Agent.Status.
const (
	AgentStatusRunning  = "RUNNING"
	AgentStatusFinished = "FINISHED"
	AgentStatusError    = "ERROR"
	AgentStatusCreating = "CREATING"
	AgentStatusExpired  = "EXPIRED"
)

// Agent represents a background agent task running in Cursor.
type Agent struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	Status    string    `json:"status"`
	Source    Source    `json:"source"`
	Target    Target    `json:"target"`
	Summary   *string   `json:"summary,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

// Source identifies the input code repository and ref.
type Source struct {
	Repository string `json:"repository"`
	Ref        string `json:"ref,omitempty"`
}

// Target identifies the output branch and PR details, if any.
type Target struct {
	BranchName   string  `json:"branchName,omitempty"`
	URL          string  `json:"url"`
	PRURL        *string `json:"prUrl,omitempty"`
	AutoCreatePR bool    `json:"autoCreatePr,omitempty"`
}

// Prompt contains instructions for the agent and optional images.
type Prompt struct {
	Text   string  `json:"text"`
	Images []Image `json:"images,omitempty"`
}

// Image provides base64-encoded data and dimensions for visual context.
type Image struct {
	Data      string     `json:"data"`
	Dimension *Dimension `json:"dimension,omitempty"`
}

// Dimension specifies width and height in pixels.
type Dimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Conversation holds the chat history of an agent.
type Conversation struct {
	ID       string    `json:"id"`
	Messages []Message `json:"messages"`
}

// Message is a single message in a conversation.
type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}

// LaunchRequest launches a new background agent.
type LaunchRequest struct {
	Prompt  Prompt         `json:"prompt"`
	Source  Source         `json:"source"`
	Model   string         `json:"model,omitempty"`
	Target  *LaunchTarget  `json:"target,omitempty"`
	Webhook *LaunchWebhook `json:"webhook,omitempty"`
}

// FollowupRequest adds instructions to a running agent.
type FollowupRequest struct {
	Prompt Prompt `json:"prompt"`
}

// FollowupResponse contains the agent identifier.
type FollowupResponse struct {
	ID string `json:"id"`
}

// DeleteResponse contains the identifier of the deleted agent.
type DeleteResponse struct {
	ID string `json:"id"`
}

// ListAgentsResponse is returned by GET /v0/agents.
type ListAgentsResponse struct {
	Agents     []Agent `json:"agents"`
	NextCursor *string `json:"nextCursor,omitempty"`
}

// LaunchTarget configures branch and PR behavior for a new agent.
type LaunchTarget struct {
	AutoCreatePR bool   `json:"autoCreatePr,omitempty"`
	BranchName   string `json:"branchName,omitempty"`
}

// LaunchWebhook configures webhook delivery for a new agent.
type LaunchWebhook struct {
	URL    string `json:"url"`
	Secret string `json:"secret,omitempty"`
}

// Repository describes a GitHub repository accessible via integration.
type Repository struct {
	Owner      string `json:"owner"`
	Name       string `json:"name"`
	Repository string `json:"repository"`
}

// ListRepositoriesResponse is returned by GET /v0/repositories.
type ListRepositoriesResponse struct {
	Repositories []Repository `json:"repositories"`
}

// ListModelsResponse contains available model names.
type ListModelsResponse struct {
	Models []string `json:"models"`
}

// MeResponse contains metadata about the API key.
type MeResponse struct {
	APIKeyName string    `json:"apiKeyName"`
	CreatedAt  time.Time `json:"createdAt"`
	UserEmail  string    `json:"userEmail"`
}

// WebhookEvent is the payload for background agent webhook notifications.
type WebhookEvent struct {
	Event     string    `json:"event"`
	Timestamp time.Time `json:"timestamp"`
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Source    Source    `json:"source"`
	Target    Target    `json:"target"`
	Summary   *string   `json:"summary,omitempty"`
}
