package live_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
	"github.com/NolanFogarty/nba-sdk/live"
)

func newTestClient(t *testing.T, fixturePath string) (*live.Client, *httptest.Server) {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", fixturePath))
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(body)
	}))
	t.Cleanup(srv.Close)

	hx := &httpx.Client{
		HTTP:     &http.Client{Timeout: 5 * time.Second},
		Limiter:  httpx.NewLimiter(100, 100),
		MaxRetry: 0,
	}
	return live.New(hx, srv.URL), srv
}

func TestScoreboard(t *testing.T) {
	client, _ := newTestClient(t, "scoreboard.json")

	resp, err := client.Scoreboard(context.Background())
	if err != nil {
		t.Fatalf("Scoreboard: %v", err)
	}

	if got, want := resp.Scoreboard.GameDate, "2026-05-15"; got != want {
		t.Errorf("GameDate = %q, want %q", got, want)
	}
	if got, want := len(resp.Scoreboard.Games), 2; got != want {
		t.Fatalf("len(Games) = %d, want %d", got, want)
	}

	game := resp.Scoreboard.Games[0]
	if game.GameID != "0042500201" {
		t.Errorf("GameID = %q, want 0042500201", game.GameID)
	}
	if game.HomeTeam.TeamTricode != "ATL" || game.AwayTeam.TeamTricode != "MIA" {
		t.Errorf("teams = %s vs %s, want ATL vs MIA", game.HomeTeam.TeamTricode, game.AwayTeam.TeamTricode)
	}
	if got, want := game.HomeTeam.Score, 48; got != want {
		t.Errorf("HomeTeam.Score = %d, want %d", got, want)
	}
	if game.GameStatus != 2 {
		t.Errorf("GameStatus = %d, want 2 (in progress)", game.GameStatus)
	}
}
