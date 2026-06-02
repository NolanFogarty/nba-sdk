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

	// Today's scoreboard.
	sb, err := client.Live.Scoreboard(ctx)
	if err != nil {
		log.Fatalf("scoreboard: %v", err)
	}
	fmt.Printf("Games on %s:\n", sb.Scoreboard.GameDate)
	for _, g := range sb.Scoreboard.Games {
		fmt.Printf("  %s %s @ %s %s — %s\n",
			g.AwayTeam.TeamTricode, scoreOrDash(g.AwayTeam.Score, g.GameStatus),
			g.HomeTeam.TeamTricode, scoreOrDash(g.HomeTeam.Score, g.GameStatus),
			g.GameStatusText)
	}

	// Look up the games played on a past date to resolve a game ID. Unlike
	// the live scoreboard above, ScoreboardV2 works for any date.
	gameID := "0022400001" // fallback if the date lookup returns nothing
	day, err := client.Stats.ScoreboardV2(ctx, "2024-12-25")
	if err != nil {
		log.Printf("scoreboardV2: %v", err)
	} else {
		fmt.Printf("\nGames on %s:\n", day.Parameters.GameDate)
		for _, g := range day.Games {
			fmt.Printf("  %s (game %s) — %s\n", g.GameCode, g.GameID, g.GameStatusText)
		}
		if len(day.Games) > 0 {
			gameID = day.Games[0].GameID
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
