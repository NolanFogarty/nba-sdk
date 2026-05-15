package nba

import (
	"net/http"
	"time"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
	"github.com/NolanFogarty/nba-sdk/live"
	"github.com/NolanFogarty/nba-sdk/stats"
)

// Client is the top-level NBA SDK entry point. Construct one with NewClient
// and access endpoint groups via the Stats and Live fields.
type Client struct {
	Stats *stats.Client
	Live  *live.Client
}

// Option configures a Client at construction time.
type Option func(*config)

type config struct {
	httpClient    *http.Client
	userAgent     string
	statsRate     float64
	statsBurst    int
	liveRate      float64
	liveBurst     int
	maxRetry      int
	retryBackoff  time.Duration
	statsBaseURL  string
	liveBaseURL   string
}

func defaults() *config {
	return &config{
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		statsRate:    1.0,
		statsBurst:   3,
		liveRate:     5.0,
		liveBurst:    10,
		maxRetry:     3,
		retryBackoff: 500 * time.Millisecond,
		statsBaseURL: "https://stats.nba.com",
		liveBaseURL:  "https://cdn.nba.com",
	}
}

// WithHTTPClient swaps in a custom *http.Client. Useful for tests or for
// callers who want to inject their own transport (proxy, logging, etc.).
func WithHTTPClient(c *http.Client) Option {
	return func(cfg *config) { cfg.httpClient = c }
}

// WithUserAgent overrides the User-Agent header.
func WithUserAgent(ua string) Option {
	return func(cfg *config) { cfg.userAgent = ua }
}

// WithStatsRateLimit sets the rate (requests/sec) and burst for stats.nba.com.
func WithStatsRateLimit(rate float64, burst int) Option {
	return func(cfg *config) { cfg.statsRate = rate; cfg.statsBurst = burst }
}

// WithLiveRateLimit sets the rate (requests/sec) and burst for cdn.nba.com.
func WithLiveRateLimit(rate float64, burst int) Option {
	return func(cfg *config) { cfg.liveRate = rate; cfg.liveBurst = burst }
}

// WithRetry configures the maximum retry count and base backoff for transient failures.
func WithRetry(maxRetry int, backoff time.Duration) Option {
	return func(cfg *config) { cfg.maxRetry = maxRetry; cfg.retryBackoff = backoff }
}

// WithStatsBaseURL overrides the stats.nba.com base URL. Intended for tests.
func WithStatsBaseURL(u string) Option {
	return func(cfg *config) { cfg.statsBaseURL = u }
}

// WithLiveBaseURL overrides the cdn.nba.com base URL. Intended for tests.
func WithLiveBaseURL(u string) Option {
	return func(cfg *config) { cfg.liveBaseURL = u }
}

// NewClient builds a Client with the supplied options applied over defaults.
func NewClient(opts ...Option) *Client {
	cfg := defaults()
	for _, opt := range opts {
		opt(cfg)
	}

	statsHTTP := &httpx.Client{
		HTTP:      cfg.httpClient,
		UserAgent: cfg.userAgent,
		Limiter:   httpx.NewLimiter(cfg.statsRate, cfg.statsBurst),
		MaxRetry:  cfg.maxRetry,
		Backoff:   cfg.retryBackoff,
	}
	liveHTTP := &httpx.Client{
		HTTP:      cfg.httpClient,
		UserAgent: cfg.userAgent,
		Limiter:   httpx.NewLimiter(cfg.liveRate, cfg.liveBurst),
		MaxRetry:  cfg.maxRetry,
		Backoff:   cfg.retryBackoff,
	}

	return &Client{
		Stats: stats.New(statsHTTP, cfg.statsBaseURL),
		Live:  live.New(liveHTTP, cfg.liveBaseURL),
	}
}
