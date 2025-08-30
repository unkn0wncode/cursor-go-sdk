# Cursor Go SDK

Go client for Cursor Background Agents — typed, minimal, and easy to use.

Official API reference that this SDK is based on:
https://docs.cursor.com/en/background-agent/api/overview


## Features

- Agents: launch, follow up, get status, list, delete, and fetch conversation.
- Models: list available model names.
- GitHub: list repositories available via your Cursor GitHub integration.
- Webhooks: HMAC-SHA256 signature verification helper for agent events.
- Config: sensible defaults + environment variable support.
- Errors: rich `APIError` with HTTP status, code, and message.


## Installation

```bash
go get https://github.com/unkn0wncode/cursor-go-sdk
```

## Quickstart

Initialize a client with just an API key, or via environment variables.

```go
package main

import (
    "context"
    "fmt"
    "github.com/unkn0wncode/cursor-go-sdk"
)

func main() {
    ctx := context.Background()

    // Option A: direct key
    c := cursor.New("YOUR_CURSOR_API_KEY")

    // Option B: from environment
    // Set up environment in shell:
    //   export CURSOR_API_KEY=...  (required)
    //   export CURSOR_USER_AGENT=my-app/1.0 (optional)
    //   export CURSOR_BASE_URL=https://api.cursor.com (default)
    //   export CURSOR_TIMEOUT_SECONDS=60 (optional)
    // cfg, _ := cursor.ConfigFromEnv()
    // c = cursor.NewClientFromConfig(cfg)

    me, err := c.Me(ctx)
    if err != nil { panic(err) }
    fmt.Println("Hello,", me.UserEmail)
}
```


## Common Tasks

### Launch an Agent

```go
agent, err := c.LaunchAgent(ctx, cursor.LaunchRequest{
    Prompt: cursor.Prompt{Text: "Add a README.md file with installation instructions."},
    Source: cursor.Source{Repository: "https://github.com/owner/repo"},
    Target: &cursor.LaunchTarget{AutoCreatePR: false},
})
if err != nil { /* handle */ }
fmt.Println("agent:", agent.ID, agent.Status)
```

### Add a Follow-up Instruction

```go
_, err = c.AddFollowup(ctx, agent.ID, cursor.FollowupRequest{
    Prompt: cursor.Prompt{Text: "Also add a section for troubleshooting."},
})
if err != nil { /* handle */ }
```

### Check Status / Conversation

```go
got, _ := c.GetAgent(ctx, agent.ID)
fmt.Println("status:", got.Status)

conv, _ := c.GetConversation(ctx, agent.ID)
for _, m := range conv.Messages {
    fmt.Println(m.Type+":", m.Text)
}
```

### List Agents (with pagination)

```go
resp, _ := c.ListAgents(ctx, 25, nil)
fmt.Println("agents:", len(resp.Agents), "next:", resp.NextCursor)
```

### Delete an Agent

```go
_, _ = c.DeleteAgent(ctx, agent.ID)
```

### List Models

```go
models, _ := c.ListModels(ctx)
fmt.Println("models:", models.Models)
```

### List GitHub Repositories

```go
repos, _ := c.ListRepositories(ctx)
fmt.Println("repos:", len(repos.Repositories))
```

Note: `ListRepositories` is rate-limited and can be slow for users with access to many repositories. Cache results and call sparingly.


## Webhooks

Background agent events can be delivered to your server via webhooks. Use the built-in signature verification helpers to check the `X-Webhook-Signature` header (HMAC-SHA256 over the raw body):

```go
// Low-level verification
ok := cursor.VerifySignature(secret, rawBody, signatureHeader)

// HTTP handler wrapper
http.Handle("/webhook", cursor.SignatureHandleWrapper(secret, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    var ev cursor.WebhookEvent
    if err := json.NewDecoder(r.Body).Decode(&ev); err != nil {
        http.Error(w, "bad request", http.StatusBadRequest)
        return
    }
    // ... handle event ...
    w.WriteHeader(http.StatusNoContent)
})))
```


## Configuration

Environment variables supported by `ConfigFromEnv`:

- `CURSOR_API_KEY`: required
- `CURSOR_BASE_URL`: default `https://api.cursor.com`
- `CURSOR_USER_AGENT`: optional custom User-Agent
- `CURSOR_TIMEOUT_SECONDS`: optional HTTP timeout override


## Errors

Non-2xx responses return `*cursor.APIError` containing:

- `StatusCode`: HTTP status code
- `Code`: API error code (if provided)
- `Message`: API error message (if provided)
- `Body`: raw response body


## Integration Tests

The file `api_endpoints_test.go` contains live integration tests. To run them, set an API key and (optionally) the repository to use for agent tests:

```bash
export CURSOR_API_KEY=... # do not commit!
# Optional: override autodetection for agent tests
export CURSOR_TEST_REPOSITORY="https://github.com/owner/repo"

go test -v
```


## Notes

- Background Agent API is currently main feature of Cursor API, but this SDK will be updated to support other features as they are added.
- Be mindful of rate limits, especially for repository listing.
- Remember to never commit real API keys.


## Reference

- Background Agents API: https://docs.cursor.com/en/background-agent/api/overview

