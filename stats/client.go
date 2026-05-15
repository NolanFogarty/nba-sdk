// Package stats wraps stats.nba.com endpoints.
package stats

import "github.com/NolanFogarty/nba-sdk/internal/httpx"

// Client groups all stats.nba.com endpoint methods.
type Client struct {
	http    *httpx.Client
	baseURL string
}

// New constructs a stats Client. Most callers should use the root nba.NewClient
// constructor instead of building this directly.
func New(http *httpx.Client, baseURL string) *Client {
	return &Client{http: http, baseURL: baseURL}
}
