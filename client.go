package cursor

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
)

// Client is the entrypoint for interacting with the Cursor Background Agents API.
// It wraps HTTP details and provides typed methods for each endpoint.
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiKey     string
	userAgent  string
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) { c.baseURL = baseURL }
}

// WithHTTPClient sets a custom http.Client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) { c.httpClient = hc }
}

// WithUserAgent sets a custom User-Agent header value.
func WithUserAgent(ua string) Option {
	return func(c *Client) { c.userAgent = ua }
}

// New creates a new Client with the provided API key.
func New(apiKey string, opts ...Option) *Client {
	cfg := Config{APIKey: apiKey}
	return NewClientFromConfig(cfg, opts...)
}

// NewClientFromConfig creates a client from Config and Options.
func NewClientFromConfig(cfg Config, opts ...Option) *Client {
	cfg.applyDefaults()
	c := &Client{
		baseURL:    cfg.BaseURL,
		httpClient: cfg.HTTPClient,
		apiKey:     cfg.APIKey,
		userAgent:  cfg.UserAgent,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// do performs an HTTP request and decodes the JSON response into out if non-nil.
func (c *Client) do(ctx context.Context, method, path string, query url.Values, body any, out any) error {
	fullURL, err := url.JoinPath(c.baseURL, path)
	if err != nil {
		return err
	}
	if query != nil {
		u, err := url.Parse(fullURL)
		if err != nil {
			return err
		}
		u.RawQuery = query.Encode()
		fullURL = u.String()
	}

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return &APIError{StatusCode: resp.StatusCode, Body: string(b)}
	}

	if out == nil {
		io.Copy(io.Discard, resp.Body)
		return nil
	}
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(out)
}
