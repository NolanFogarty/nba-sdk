package live_test

import (
	"context"
	"strings"
	"testing"
)

func TestBoxScore(t *testing.T) {
	client, _ := newTestClient(t, "boxscore.json")

	resp, err := client.BoxScore(context.Background(), "0042500201")
	if err != nil {
		t.Fatalf("BoxScore: %v", err)
	}

	if got, want := resp.Game.GameID, "0042500201"; got != want {
		t.Errorf("GameID = %q, want %q", got, want)
	}
	if got, want := resp.Game.GameStatus, 2; got != want {
		t.Errorf("GameStatus = %d, want %d", got, want)
	}
	if got, want := resp.Game.Arena.ArenaName, "State Farm Arena"; got != want {
		t.Errorf("Arena.ArenaName = %q, want %q", got, want)
	}
	if len(resp.Game.Officials) == 0 {
		t.Fatal("expected at least one official")
	}

	home := resp.Game.HomeTeam
	if home.TeamTricode != "ATL" || home.Score != 48 {
		t.Errorf("HomeTeam = %s/%d, want ATL/48", home.TeamTricode, home.Score)
	}
	if len(home.Players) == 0 {
		t.Fatal("expected at least one home player")
	}

	young := home.Players[0]
	if young.FamilyName != "Young" {
		t.Errorf("first home player = %q, want Young", young.FamilyName)
	}
	if young.Starter != "1" {
		t.Errorf("Starter = %q, want \"1\" (NBA returns these as strings)", young.Starter)
	}
	if young.Statistics.Points != 18 {
		t.Errorf("Young points = %d, want 18", young.Statistics.Points)
	}
	if !strings.HasPrefix(young.Statistics.Minutes, "PT") {
		t.Errorf("Minutes = %q, want ISO-8601 duration", young.Statistics.Minutes)
	}

	if home.Statistics.BiggestLead != 9 {
		t.Errorf("HomeTeam.Statistics.BiggestLead = %d, want 9", home.Statistics.BiggestLead)
	}
	if home.Statistics.TrueShootingPercentage <= 0 {
		t.Errorf("TrueShootingPercentage = %v, want > 0", home.Statistics.TrueShootingPercentage)
	}
}
