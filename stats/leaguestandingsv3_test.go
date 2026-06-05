package stats_test

import (
	"context"
	"testing"
)

func TestLeagueStandingsV3(t *testing.T) {
	client, _ := newTestClient(t, "leaguestandingsv3.json", "/stats/leaguestandingsv3")

	resp, err := client.LeagueStandingsV3(context.Background(), "2025-26")
	if err != nil {
		t.Fatalf("LeagueStandingsV3: %v", err)
	}

	if got, want := resp.Resource, "leaguestandingsv3"; got != want {
		t.Errorf("Resource = %q, want %q", got, want)
	}
	if got, want := resp.Parameters.Season, "2025-26"; got != want {
		t.Errorf("Parameters.Season = %q, want %q", got, want)
	}
	if got, want := resp.Parameters.SeasonType, "Regular Season"; got != want {
		t.Errorf("Parameters.SeasonType = %q, want %q", got, want)
	}
	if got, want := len(resp.Standings), 2; got != want {
		t.Fatalf("len(Standings) = %d, want %d", got, want)
	}

	// First row: Thunder, #1 in West.
	okc := resp.Standings[0]
	if okc.TeamCity != "Oklahoma City" || okc.TeamName != "Thunder" {
		t.Errorf("first team = %s %s, want Oklahoma City Thunder", okc.TeamCity, okc.TeamName)
	}
	if okc.Wins != 64 || okc.Losses != 18 {
		t.Errorf("OKC record = %d-%d, want 64-18", okc.Wins, okc.Losses)
	}
	if okc.WinPCT < 0.77 || okc.WinPCT > 0.79 {
		t.Errorf("OKC WinPCT = %v, want ~0.78", okc.WinPCT)
	}
	if okc.Conference != "West" {
		t.Errorf("OKC Conference = %q, want West", okc.Conference)
	}
	if okc.PlayoffRank != 1 {
		t.Errorf("OKC PlayoffRank = %d, want 1", okc.PlayoffRank)
	}
	if okc.ClinchedPlayoffBirth != 1 {
		t.Errorf("OKC ClinchedPlayoffBirth = %d, want 1", okc.ClinchedPlayoffBirth)
	}
	if okc.L10 != "8-2" {
		t.Errorf("OKC L10 = %q, want 8-2", okc.L10)
	}
	if okc.PointsPG <= 0 {
		t.Errorf("OKC PointsPG = %v, want > 0", okc.PointsPG)
	}
	if okc.StrCurrentStreak != "8W" {
		t.Errorf("OKC StrCurrentStreak = %q, want 8W", okc.StrCurrentStreak)
	}

	// Second row: Nuggets.
	den := resp.Standings[1]
	if den.TeamName != "Nuggets" || den.PlayoffRank != 4 {
		t.Errorf("second team = %s rank %d, want Nuggets rank 4", den.TeamName, den.PlayoffRank)
	}
	if den.ConferenceGamesBack != 7.5 {
		t.Errorf("DEN ConferenceGamesBack = %v, want 7.5", den.ConferenceGamesBack)
	}
}
