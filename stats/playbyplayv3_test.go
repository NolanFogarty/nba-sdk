package stats_test

import (
	"context"
	"testing"

	"github.com/NolanFogarty/nba-sdk/stats"
)

func TestPlayByPlayV3(t *testing.T) {
	client, _ := newTestClient(t, "playbyplayv3.json", "/stats/playbyplayv3")

	resp, err := client.PlayByPlayV3(context.Background(), "0022400001")
	if err != nil {
		t.Fatalf("PlayByPlayV3: %v", err)
	}

	if got, want := resp.Game.GameID, "0022400001"; got != want {
		t.Errorf("GameID = %q, want %q", got, want)
	}
	if got, want := len(resp.Game.Actions), 4; got != want {
		t.Fatalf("len(Actions) = %d, want %d", got, want)
	}

	// Demonstrate filtering by quarter on the caller side.
	var q1 []stats.Action
	for _, a := range resp.Game.Actions {
		if a.Period == 1 {
			q1 = append(q1, a)
		}
	}
	if got, want := len(q1), 3; got != want {
		t.Errorf("Q1 actions = %d, want %d", got, want)
	}

	first := resp.Game.Actions[0]
	if first.ActionType != "period" || first.SubType != "start" {
		t.Errorf("first action = %s/%s, want period/start", first.ActionType, first.SubType)
	}
}
