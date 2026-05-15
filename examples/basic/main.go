// Command basic is a tiny usage example for the nba SDK.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	nba "github.com/NolanFogarty/nba-sdk"
	"github.com/NolanFogarty/nba-sdk/live"
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

	// Pick a finished game from the scoreboard if one exists (gameStatus == 3),
	// otherwise fall back to a known finished game ID so the stats demo always
	// has real data to fetch. stats.nba.com returns nothing for a game that
	// hasn't tipped off yet.
	gameID := pickFinishedGame(sb.Scoreboard.Games)
	if gameID == "" {
		gameID = "0022400001" // 2024-25 opening night, Knicks @ Celtics
		fmt.Printf("\n(No finished games today — using %s for stats demo.)\n", gameID)
	} else {
		fmt.Printf("\nUsing finished game %s for stats demo.\n", gameID)
	}

	// Traditional box score for the first game.
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

// pickFinishedGame returns the first GameID with status 3 (final), or "" if none.
func pickFinishedGame(games []live.Game) string {
	for _, g := range games {
		if g.GameStatus == 3 {
			return g.GameID
		}
	}
	return ""
}
