package stats_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
	"github.com/NolanFogarty/nba-sdk/stats"
)

func newTestClient(t *testing.T, fixturePath, expectedPath string) (*stats.Client, *httptest.Server) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", fixturePath))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, expectedPath) {
			t.Errorf("unexpected path %q, want prefix %q", r.URL.Path, expectedPath)
		}
		if r.Header.Get("x-nba-stats-token") != "true" {
			t.Errorf("missing x-nba-stats-token header")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	t.Cleanup(srv.Close)

	hx := &httpx.Client{
		HTTP:     &http.Client{Timeout: 5 * time.Second},
		Limiter:  httpx.NewLimiter(100, 100),
		MaxRetry: 0,
	}
	return stats.New(hx, srv.URL), srv
}

func TestBoxScoreTraditionalV3(t *testing.T) {
	client, _ := newTestClient(t, "boxscoretraditionalv3.json", "/stats/boxscoretraditionalv3")

	resp, err := client.BoxScoreTraditionalV3(context.Background(), "0022400001")
	if err != nil {
		t.Fatalf("BoxScoreTraditionalV3: %v", err)
	}

	if got, want := resp.BoxScoreTraditional.GameID, "0022400001"; got != want {
		t.Errorf("GameID = %q, want %q", got, want)
	}
	if got, want := resp.BoxScoreTraditional.HomeTeam.TeamTricode, "LAL"; got != want {
		t.Errorf("HomeTeam.TeamTricode = %q, want %q", got, want)
	}
	if got, want := resp.BoxScoreTraditional.AwayTeam.TeamTricode, "BOS"; got != want {
		t.Errorf("AwayTeam.TeamTricode = %q, want %q", got, want)
	}
	if got, want := resp.BoxScoreTraditional.HomeTeam.Statistics.Points, 110; got != want {
		t.Errorf("HomeTeam.Statistics.Points = %d, want %d", got, want)
	}
	if len(resp.BoxScoreTraditional.HomeTeam.Players) == 0 {
		t.Fatal("expected at least one home team player")
	}
	if got, want := resp.BoxScoreTraditional.HomeTeam.Players[0].FamilyName, "James"; got != want {
		t.Errorf("first home player FamilyName = %q, want %q", got, want)
	}
}
