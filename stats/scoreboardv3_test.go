package stats_test

import (
	"context"
	"testing"
	"time"
)

func TestScoreboardV3(t *testing.T) {
	client, _ := newTestClient(t, "scoreboardv3.json", "/stats/scoreboardv3")

	date := time.Date(2026, 6, 5, 0, 0, 0, 0, time.UTC)
	resp, err := client.ScoreboardV3(context.Background(), date)
	if err != nil {
		t.Fatalf("ScoreboardV3: %v", err)
	}

	if got, want := resp.Scoreboard.GameDate, "2026-06-05"; got != want {
		t.Errorf("GameDate = %q, want %q", got, want)
	}
	if got, want := len(resp.Scoreboard.Games), 1; got != want {
		t.Fatalf("len(Games) = %d, want %d", got, want)
	}

	g := resp.Scoreboard.Games[0]
	if g.GameID != "0042500404" {
		t.Errorf("GameID = %q, want 0042500404", g.GameID)
	}
	if g.GameStatus != 1 {
		t.Errorf("GameStatus = %d, want 1 (scheduled)", g.GameStatus)
	}
	if g.HomeTeam.TeamTricode != "OKC" || g.AwayTeam.TeamTricode != "IND" {
		t.Errorf("teams = %s vs %s, want OKC vs IND", g.HomeTeam.TeamTricode, g.AwayTeam.TeamTricode)
	}
	if g.SeriesText != "OKC leads 2-1" {
		t.Errorf("SeriesText = %q, want %q", g.SeriesText, "OKC leads 2-1")
	}
	if g.PoRoundDesc != "Finals" {
		t.Errorf("PoRoundDesc = %q, want Finals", g.PoRoundDesc)
	}
	if len(g.Broadcasters.NationalTvBroadcasters) == 0 || g.Broadcasters.NationalTvBroadcasters[0].BroadcasterDisplay != "ABC" {
		t.Errorf("expected national TV broadcaster ABC")
	}
}
