package cursor

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

// Config contains settings for constructing a Client.
type Config struct {
	APIKey         string
	BaseURL        string
	UserAgent      string
	HTTPClient     *http.Client
	TimeoutSeconds *int
}

// ConfigFromEnv reads configuration from environment variables.
// Supported variables:
// - CURSOR_API_KEY (required)
// - CURSOR_BASE_URL (default: https://api.cursor.com)
// - CURSOR_USER_AGENT (optional)
// - CURSOR_TIMEOUT_SECONDS (optional)
func ConfigFromEnv() (Config, error) {
	c := Config{}
	c.APIKey = os.Getenv("CURSOR_API_KEY")
	c.BaseURL = os.Getenv("CURSOR_BASE_URL")
	if c.BaseURL == "" {
		c.BaseURL = "https://api.cursor.com"
	}
	if ua := os.Getenv("CURSOR_USER_AGENT"); ua != "" {
		c.UserAgent = ua
	}
	if ts := os.Getenv("CURSOR_TIMEOUT_SECONDS"); ts != "" {
		if v, err := strconv.Atoi(ts); err == nil && v > 0 {
			c.TimeoutSeconds = &v
		}
	}
	return c, nil
}

// applyDefaults fills unset fields with sensible defaults.
func (c *Config) applyDefaults() {
	if c.BaseURL == "" {
		c.BaseURL = "https://api.cursor.com"
	}
	if c.HTTPClient == nil {
		c.HTTPClient = &http.Client{Timeout: 60 * time.Second}
	}
	if c.TimeoutSeconds != nil {
		c.HTTPClient.Timeout = time.Duration(*c.TimeoutSeconds) * time.Second
	}
	if c.UserAgent == "" {
		c.UserAgent = "cursor-go-sdk/0.1"
	}
}
