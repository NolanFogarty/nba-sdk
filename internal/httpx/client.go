// Package httpx is an internal HTTP client wrapper for the NBA SDK. It applies
// the headers stats.nba.com requires, enforces a token-bucket rate limit, and
// retries on transient failures.
package httpx

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"net/url"
	"time"
)

const (
	// Windows Chrome 120 — Akamai's bot detector currently accepts this fingerprint
	// when paired with the sec-ch-ua / sec-fetch headers below. A Mac UA was being
	// silently dropped (connection accepted, response held open indefinitely).
	defaultUserAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"
)

// Profile selects the header set used for a request. stats.nba.com requires a
// specific set of headers or returns 403; cdn.nba.com is happy with a normal UA.
type Profile int

const (
	ProfileStats Profile = iota
	ProfileLive
)

// Client wraps http.Client with rate limiting, retry, and NBA-specific headers.
type Client struct {
	HTTP      *http.Client
	UserAgent string
	Limiter   *Limiter
	MaxRetry  int
	Backoff   time.Duration
}

// ErrAPI is returned when the NBA API responds with a non-2xx status that we
// did not (or could not) retry past.
type ErrAPI struct {
	StatusCode int
	URL        string
	Body       string
}

func (e *ErrAPI) Error() string {
	return fmt.Sprintf("nba api: %s returned %d: %s", e.URL, e.StatusCode, e.Body)
}

// GetJSON issues a GET against rawURL with the headers appropriate for profile,
// retries transient failures, and decodes the JSON body into out.
func (c *Client) GetJSON(ctx context.Context, profile Profile, rawURL string, query url.Values, out any) error {
	if query != nil {
		rawURL = rawURL + "?" + query.Encode()
	}

	var lastErr error
	for attempt := 0; attempt <= c.MaxRetry; attempt++ {
		if err := c.Limiter.Wait(ctx); err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
		if err != nil {
			return err
		}
		c.applyHeaders(req, profile)

		resp, err := c.HTTP.Do(req)
		if err != nil {
			lastErr = err
			if !c.shouldRetry(ctx, attempt) {
				return err
			}
			c.sleepBackoff(ctx, attempt)
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			defer resp.Body.Close()
			return json.NewDecoder(resp.Body).Decode(out)
		}

		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		apiErr := &ErrAPI{StatusCode: resp.StatusCode, URL: rawURL, Body: truncate(string(body), 512)}

		if isRetryable(resp.StatusCode) && c.shouldRetry(ctx, attempt) {
			lastErr = apiErr
			c.sleepBackoff(ctx, attempt)
			continue
		}
		return apiErr
	}
	if lastErr == nil {
		lastErr = errors.New("nba api: exhausted retries with no error recorded")
	}
	return lastErr
}

func (c *Client) applyHeaders(req *http.Request, profile Profile) {
	ua := c.UserAgent
	if ua == "" {
		ua = defaultUserAgent
	}
	req.Header.Set("User-Agent", ua)
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Connection", "keep-alive")

	switch profile {
	case ProfileStats:
		req.Header.Set("Referer", "https://www.nba.com/")
		req.Header.Set("Origin", "https://www.nba.com")
		req.Header.Set("x-nba-stats-origin", "stats")
		req.Header.Set("x-nba-stats-token", "true")
		// Browser-fingerprint headers Akamai checks for. Without these the
		// server accepts the connection and then silently holds it open.
		req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Google Chrome";v="120"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-site")
	case ProfileLive:
		req.Header.Set("Referer", "https://www.nba.com/")
	}
}

func (c *Client) shouldRetry(ctx context.Context, attempt int) bool {
	if attempt >= c.MaxRetry {
		return false
	}
	return ctx.Err() == nil
}

func (c *Client) sleepBackoff(ctx context.Context, attempt int) {
	base := c.Backoff
	if base <= 0 {
		base = 500 * time.Millisecond
	}
	delay := base * time.Duration(1<<attempt)
	jitter := time.Duration(rand.Int64N(int64(base)))
	t := time.NewTimer(delay + jitter)
	defer t.Stop()
	select {
	case <-ctx.Done():
	case <-t.C:
	}
}

func isRetryable(status int) bool {
	if status == http.StatusTooManyRequests {
		return true
	}
	return status >= 500 && status <= 599
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
