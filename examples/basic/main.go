// Command basic is a tiny usage example for the nba SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	nba "github.com/NolanFogarty/nba-sdk"
)

func main() {
	client := nba.NewClient()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Today's scoreboard from the CDN. NBA's "today" rolls over on its own
	// schedule and may show yesterday's slate for several hours.
	sb, err := client.Live.Scoreboard(ctx)
	if err != nil {
		log.Fatalf("live scoreboard: %v", err)
	}
	fmt.Printf("Live scoreboard for %s:\n", sb.Scoreboard.GameDate)
	for _, g := range sb.Scoreboard.Games {
		fmt.Printf("  %s %s @ %s %s — %s\n",
			g.AwayTeam.TeamTricode, scoreOrDash(g.AwayTeam.Score, g.GameStatus),
			g.HomeTeam.TeamTricode, scoreOrDash(g.HomeTeam.Score, g.GameStatus),
			g.GameStatusText)
	}

	// Stats v3 scoreboard for an explicit date — bypasses the CDN's "today"
	// quirk. Pass time.Now() to force the real current calendar day.
	today, err := client.Stats.ScoreboardV3(ctx, time.Now())
	if err != nil {
		log.Printf("scoreboardV3: %v", err)
	} else {
		fmt.Printf("\nStats v3 scoreboard for %s:\n", today.Scoreboard.GameDate)
		for _, g := range today.Scoreboard.Games {
			fmt.Printf("  %s %s @ %s %s — %s\n",
				g.AwayTeam.TeamTricode, scoreOrDash(g.AwayTeam.Score, g.GameStatus),
				g.HomeTeam.TeamTricode, scoreOrDash(g.HomeTeam.Score, g.GameStatus),
				g.GameStatusText)
		}
	}

	// Look up games played on a past date to resolve a game ID. ScoreboardV3
	// accepts any calendar day — past, present, or future.
	gameID := "0022400001" // fallback if the date lookup returns nothing
	pastDate := time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)
	day, err := client.Stats.ScoreboardV3(ctx, pastDate)
	if err != nil {
		log.Printf("scoreboardV3 (past date): %v", err)
	} else {
		fmt.Printf("\nGames on %s:\n", day.Scoreboard.GameDate)
		for _, g := range day.Scoreboard.Games {
			fmt.Printf("  %s (game %s) — %s\n", g.GameCode, g.GameID, g.GameStatusText)
		}
		if len(day.Scoreboard.Games) > 0 {
			gameID = day.Scoreboard.Games[0].GameID
		}
	}

	// Traditional box score for the resolved game.
	box, err := client.Stats.BoxScoreTraditionalV3(ctx, gameID)
	if err != nil {
		log.Printf("box score: %v", err)
	} else {
		fmt.Printf("\n%s totals: %d pts, %d reb, %d ast\n",
			box.BoxScoreTraditional.HomeTeam.TeamTricode,
			box.BoxScoreTraditional.HomeTeam.Statistics.Points,
			box.BoxScoreTraditional.HomeTeam.Statistics.ReboundsTotal,
			box.BoxScoreTraditional.HomeTeam.Statistics.Assists)
	}

	// Play-by-play for the first game, filtered to Q1.
	pbp, err := client.Stats.PlayByPlayV3(ctx, gameID)
	if err != nil {
		log.Printf("play-by-play: %v", err)
	} else {
		count := 0
		for _, a := range pbp.Game.Actions {
			if a.Period == 1 {
				count++
			}
		}
		fmt.Printf("Q1 plays: %d\n", count)
	}
}

func scoreOrDash(score, status int) string {
	if status == 1 {
		return "-"
	}
	return fmt.Sprintf("%d", score)
}
