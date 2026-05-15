package live_test

import (
	"context"
	"testing"

	"github.com/NolanFogarty/nba-sdk/live"
)

func TestPlayByPlay(t *testing.T) {
	client, _ := newTestClient(t, "playbyplay.json")

	resp, err := client.PlayByPlay(context.Background(), "0042500201")
	if err != nil {
		t.Fatalf("PlayByPlay: %v", err)
	}

	if got, want := resp.Game.GameID, "0042500201"; got != want {
		t.Errorf("GameID = %q, want %q", got, want)
	}
	if got, want := len(resp.Game.Actions), 7; got != want {
		t.Fatalf("len(Actions) = %d, want %d", got, want)
	}

	// Filter to period 1 — caller-side filtering is the intended pattern.
	var q1, q2 []live.Action
	for _, a := range resp.Game.Actions {
		switch a.Period {
		case 1:
			q1 = append(q1, a)
		case 2:
			q2 = append(q2, a)
		}
	}
	if got, want := len(q1), 5; got != want {
		t.Errorf("Q1 actions = %d, want %d", got, want)
	}
	if got, want := len(q2), 2; got != want {
		t.Errorf("Q2 actions = %d, want %d", got, want)
	}

	// Verify shot fields populate on Made Shot actions.
	var madeShots []live.Action
	for _, a := range resp.Game.Actions {
		if a.ActionType == "Made Shot" {
			madeShots = append(madeShots, a)
		}
	}
	if got, want := len(madeShots), 3; got != want {
		t.Fatalf("Made Shot actions = %d, want %d", got, want)
	}
	if madeShots[0].PlayerName != "Young" || madeShots[0].ShotResult != "Made" {
		t.Errorf("first made shot: PlayerName=%q ShotResult=%q, want Young/Made",
			madeShots[0].PlayerName, madeShots[0].ShotResult)
	}
	if madeShots[0].AssistPersonID == 0 {
		t.Errorf("first made shot AssistPersonID = 0, want non-zero (assisted)")
	}

	// Foul-specific fields populate on foul actions.
	for _, a := range resp.Game.Actions {
		if a.ActionType == "foul" {
			if a.FoulDrawnPersonID == 0 {
				t.Errorf("foul action missing FoulDrawnPersonID")
			}
			break
		}
	}
}
