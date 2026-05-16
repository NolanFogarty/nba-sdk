// Command livetest fetches today's scoreboard, finds a game that is
// currently in progress (gameStatus == 2), and exercises the live CDN
// endpoints (BoxScore and PlayByPlay) against it.
//
// If no game is in progress when this runs, it prints the schedule and
// exits — the CDN endpoints only return data for live or recently
// finished games. Run during an active game to test.
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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	sb, err := client.Live.Scoreboard(ctx)
	if err != nil {
		log.Fatalf("scoreboard: %v", err)
	}

	fmt.Printf("Scoreboard for %s — %d games\n", sb.Scoreboard.GameDate, len(sb.Scoreboard.Games))
	for _, g := range sb.Scoreboard.Games {
		fmt.Printf("  [%s] %s @ %s — %s\n", statusLabel(g.GameStatus), g.AwayTeam.TeamTricode, g.HomeTeam.TeamTricode, g.GameStatusText)
	}

	game := firstInProgress(sb.Scoreboard.Games)
	if game == nil {
		fmt.Println("\nNo games currently in progress. Live CDN endpoints will 404.")
		fmt.Println("Re-run during an active game to test Live.BoxScore and Live.PlayByPlay.")
		return
	}

	fmt.Printf("\nUsing live game %s (%s @ %s) for live endpoints.\n",
		game.GameID, game.AwayTeam.TeamTricode, game.HomeTeam.TeamTricode)

	// Live box score.
	box, err := client.Live.BoxScore(ctx, game.GameID)
	if err != nil {
		log.Printf("live boxscore: %v", err)
	} else {
		h := box.Game.HomeTeam
		a := box.Game.AwayTeam
		fmt.Printf("\nLive box score — Q%d %s\n", box.Game.Period, box.Game.GameClock)
		fmt.Printf("  %s %d (TS%% %.3f, biggest lead %d)\n", a.TeamTricode, a.Score, a.Statistics.TrueShootingPercentage, a.Statistics.BiggestLead)
		fmt.Printf("  %s %d (TS%% %.3f, biggest lead %d)\n", h.TeamTricode, h.Score, h.Statistics.TrueShootingPercentage, h.Statistics.BiggestLead)
		fmt.Printf("  Arena: %s, %s — attendance %d\n", box.Game.Arena.ArenaName, box.Game.Arena.ArenaCity, box.Game.Attendance)

		var leader *live.BoxScorePlayer
		for i := range h.Players {
			p := &h.Players[i]
			if leader == nil || p.Statistics.Points > leader.Statistics.Points {
				leader = p
			}
		}
		if leader != nil {
			fmt.Printf("  Home top scorer: %s — %d pts, %d reb, %d ast\n",
				leader.Name, leader.Statistics.Points, leader.Statistics.ReboundsTotal, leader.Statistics.Assists)
		}
	}

	// Live play-by-play — count actions and show the last few.
	pbp, err := client.Live.PlayByPlay(ctx, game.GameID)
	if err != nil {
		log.Printf("live play-by-play: %v", err)
	} else {
		fmt.Printf("\nLive play-by-play — %d total actions\n", len(pbp.Game.Actions))
		tail := pbp.Game.Actions
		if len(tail) > 5 {
			tail = tail[len(tail)-5:]
		}
		fmt.Println("  Most recent:")
		for _, a := range tail {
			fmt.Printf("    Q%d %s — %s\n", a.Period, a.Clock, a.Description)
		}
	}
}

func firstInProgress(games []live.Game) *live.Game {
	for i := range games {
		if games[i].GameStatus == 2 {
			return &games[i]
		}
	}
	return nil
}

func statusLabel(s int) string {
	switch s {
	case 1:
		return "scheduled"
	case 2:
		return "LIVE"
	case 3:
		return "final"
	default:
		return fmt.Sprintf("status=%d", s)
	}
}
