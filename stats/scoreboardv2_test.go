package stats_test

import (
	"context"
	"testing"
)

func TestScoreboardV2(t *testing.T) {
	client, _ := newTestClient(t, "scoreboardv2.json", "/stats/scoreboardv2")

	resp, err := client.ScoreboardV2(context.Background(), "2024-12-25")
	if err != nil {
		t.Fatalf("ScoreboardV2: %v", err)
	}

	if got, want := resp.Parameters.GameDate, "2024-12-25"; got != want {
		t.Errorf("Parameters.GameDate = %q, want %q", got, want)
	}
	if got, want := len(resp.Games), 2; got != want {
		t.Fatalf("len(Games) = %d, want %d", got, want)
	}

	first := resp.Games[0]
	if got, want := first.GameID, "0022400384"; got != want {
		t.Errorf("Games[0].GameID = %q, want %q", got, want)
	}
	if got, want := first.GameStatusText, "Final"; got != want {
		t.Errorf("Games[0].GameStatusText = %q, want %q", got, want)
	}
	if got, want := first.HomeTeamID, 1610612759; got != want {
		t.Errorf("Games[0].HomeTeamID = %d, want %d", got, want)
	}
	if got, want := first.GameCode, "20241225/NYKSAS"; got != want {
		t.Errorf("Games[0].GameCode = %q, want %q", got, want)
	}

	// The game IDs surfaced here can be fed straight into the v3 endpoints.
	if got, want := resp.Games[1].GameID, "0022400385"; got != want {
		t.Errorf("Games[1].GameID = %q, want %q", got, want)
	}

	// Raw result sets remain available for callers who need more.
	if len(resp.ResultSets) != 2 {
		t.Errorf("len(ResultSets) = %d, want 2", len(resp.ResultSets))
	}
}
