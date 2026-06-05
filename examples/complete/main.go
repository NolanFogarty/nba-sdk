// Command complete is an exhaustive tour of the nba SDK. It calls every
// endpoint the SDK exposes and prints every field of every response struct,
// so it doubles as living documentation of "all available values".
//
// Endpoints exercised:
//
//	Live  (cdn.nba.com)     — Scoreboard, BoxScore, PlayByPlay
//	Stats (stats.nba.com)   — ScoreboardV3, BoxScoreTraditionalV3, PlayByPlayV3
//
// Game-ID resolution:
//
//   - The stats.* endpoints work for any historical game, so we resolve a
//     real game ID from ScoreboardV3 for a fixed past date.
//   - The live.* CDN endpoints only return data for games that are in
//     progress or very recently finished. We pick a live game from today's
//     scoreboard if one exists; otherwise we still attempt the call against
//     the first scheduled game and report the (expected) 404 so you can see
//     the call shape. Run during an active game for full live output.
//
// For collections that can be huge (players, play-by-play actions) we print
// the full field detail for one representative element and summarize the
// rest — every field name still appears at least once.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	nba "github.com/NolanFogarty/nba-sdk"
	"github.com/NolanFogarty/nba-sdk/live"
	"github.com/NolanFogarty/nba-sdk/stats"
)

// fallbackGameID is used if the date lookup returns nothing.
const fallbackGameID = "0022400001"

// pastDate is any date with a full slate of completed games; used to
// resolve a real historical game ID for the stats.* endpoints.
var pastDate = time.Date(2024, 12, 25, 0, 0, 0, 0, time.UTC)

func main() {
	client := nba.NewClient(
		// Defaults are conservative; shown here for completeness.
		nba.WithStatsRateLimit(1.0, 3),
		nba.WithLiveRateLimit(5.0, 10),
		nba.WithRetry(3, 500*time.Millisecond),
		nba.WithUserAgent("nba-sdk-complete-example/1.0"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// ----- LIVE: today's scoreboard -------------------------------------
	header("LIVE  client.Live.Scoreboard(ctx)")
	sb, err := client.Live.Scoreboard(ctx)
	if err != nil {
		log.Fatalf("live scoreboard: %v", err)
	}
	printLiveScoreboard(sb)

	// Pick a game for the live endpoints: prefer one in progress.
	liveGame := firstInProgress(sb.Scoreboard.Games)
	liveLabel := "in-progress"
	if liveGame == nil && len(sb.Scoreboard.Games) > 0 {
		liveGame = &sb.Scoreboard.Games[0]
		liveLabel = "scheduled (live CDN endpoints will likely 404)"
	}

	// ----- LIVE: box score + play-by-play -------------------------------
	if liveGame == nil {
		fmt.Println("\nNo games on today's scoreboard; skipping Live.BoxScore / Live.PlayByPlay.")
	} else {
		fmt.Printf("\nUsing game %s (%s @ %s) — %s\n",
			liveGame.GameID, liveGame.AwayTeam.TeamTricode, liveGame.HomeTeam.TeamTricode, liveLabel)

		header("LIVE  client.Live.BoxScore(ctx, gameID)")
		if box, err := client.Live.BoxScore(ctx, liveGame.GameID); err != nil {
			fmt.Printf("(no live box score available: %v)\n", err)
		} else {
			printLiveBoxScore(box)
		}

		header("LIVE  client.Live.PlayByPlay(ctx, gameID)")
		if pbp, err := client.Live.PlayByPlay(ctx, liveGame.GameID); err != nil {
			fmt.Printf("(no live play-by-play available: %v)\n", err)
		} else {
			printLivePlayByPlay(pbp)
		}
	}

	// ----- STATS: scoreboard v3 -----------------------------------------
	header("STATS client.Stats.ScoreboardV3(ctx, date)")
	gameID := fallbackGameID
	if day, err := client.Stats.ScoreboardV3(ctx, pastDate); err != nil {
		log.Printf("scoreboardV3: %v", err)
	} else {
		printScoreboardV3(day)
		if len(day.Scoreboard.Games) > 0 {
			gameID = day.Scoreboard.Games[0].GameID
		}
	}
	fmt.Printf("\nResolved historical game ID for stats endpoints: %s\n", gameID)

	// ----- STATS: traditional box score (v3) ----------------------------
	header("STATS client.Stats.BoxScoreTraditionalV3(ctx, gameID)")
	if box, err := client.Stats.BoxScoreTraditionalV3(ctx, gameID); err != nil {
		log.Printf("boxscoretraditionalv3: %v", err)
	} else {
		printBoxScoreTraditionalV3(box)
	}

	// ----- STATS: play-by-play (v3) -------------------------------------
	header("STATS client.Stats.PlayByPlayV3(ctx, gameID)")
	if pbp, err := client.Stats.PlayByPlayV3(ctx, gameID); err != nil {
		log.Printf("playbyplayv3: %v", err)
	} else {
		printPlayByPlayV3(pbp)
	}
}

// ===================================================================
// LIVE: Scoreboard
// ===================================================================

func printLiveScoreboard(r *live.ScoreboardResponse) {
	printMetaLive(r.Meta)
	s := r.Scoreboard
	fmt.Println("scoreboard:")
	kv(1, "gameDate", s.GameDate)
	kv(1, "leagueId", s.LeagueID)
	kv(1, "leagueName", s.LeagueName)
	kv(1, "games", fmt.Sprintf("%d total", len(s.Games)))

	for i := range s.Games {
		g := &s.Games[i]
		fmt.Printf("  game[%d]:\n", i)
		kv(2, "gameId", g.GameID)
		kv(2, "gameCode", g.GameCode)
		kv(2, "gameStatus", g.GameStatus)
		kv(2, "gameStatusText", g.GameStatusText)
		kv(2, "period", g.Period)
		kv(2, "gameClock", g.GameClock)
		kv(2, "gameTimeUTC", g.GameTimeUTC)
		kv(2, "gameEt", g.GameEt)
		kv(2, "regulationPeriods", g.RegulationPeriods)
		kv(2, "seriesGameNumber", g.SeriesGameNumber)
		kv(2, "seriesText", g.SeriesText)
		kv(2, "ifNecessary", g.IfNecessary)
		fmt.Println("    awayTeam:")
		printGameTeam(g.AwayTeam)
		fmt.Println("    homeTeam:")
		printGameTeam(g.HomeTeam)
		fmt.Println("    gameLeaders:")
		fmt.Println("      awayLeaders:")
		printLeader(g.GameLeaders.AwayLeaders)
		fmt.Println("      homeLeaders:")
		printLeader(g.GameLeaders.HomeLeaders)
		if i == 0 && len(s.Games) > 1 {
			fmt.Printf("  (... %d more games, same shape ...)\n", len(s.Games)-1)
			break
		}
	}
}

func printGameTeam(t live.GameTeam) {
	kv(3, "teamId", t.TeamID)
	kv(3, "teamName", t.TeamName)
	kv(3, "teamCity", t.TeamCity)
	kv(3, "teamTricode", t.TeamTricode)
	kv(3, "wins", t.Wins)
	kv(3, "losses", t.Losses)
	kv(3, "score", t.Score)
	kv(3, "seed", t.Seed)
	kv(3, "inBonus", t.InBonus)
	kv(3, "timeoutsRemaining", t.TimeoutsRemaining)
	kv(3, "periods", fmt.Sprintf("%d", len(t.Periods)))
	for _, p := range t.Periods {
		kv(4, fmt.Sprintf("P%d (%s)", p.Period, p.PeriodType), p.Score)
	}
}

func printLeader(l live.Leader) {
	kv(4, "personId", l.PersonID)
	kv(4, "name", l.Name)
	kv(4, "jerseyNum", l.JerseyNum)
	kv(4, "position", l.Position)
	kv(4, "teamTricode", l.TeamTricode)
	kv(4, "playerSlug", l.PlayerSlug)
	kv(4, "points", l.Points)
	kv(4, "rebounds", l.Rebounds)
	kv(4, "assists", l.Assists)
}

// ===================================================================
// LIVE: BoxScore
// ===================================================================

func printLiveBoxScore(r *live.BoxScoreResponse) {
	printMetaLive(r.Meta)
	g := r.Game
	fmt.Println("game:")
	kv(1, "gameId", g.GameID)
	kv(1, "gameTimeLocal", g.GameTimeLocal)
	kv(1, "gameTimeUTC", g.GameTimeUTC)
	kv(1, "gameTimeHome", g.GameTimeHome)
	kv(1, "gameTimeAway", g.GameTimeAway)
	kv(1, "gameEt", g.GameEt)
	kv(1, "duration", g.Duration)
	kv(1, "gameCode", g.GameCode)
	kv(1, "gameStatusText", g.GameStatusText)
	kv(1, "gameStatus", g.GameStatus)
	kv(1, "regulationPeriods", g.RegulationPeriods)
	kv(1, "period", g.Period)
	kv(1, "gameClock", g.GameClock)
	kv(1, "attendance", g.Attendance)
	kv(1, "sellout", g.Sellout)

	fmt.Println("  arena:")
	kv(2, "arenaId", g.Arena.ArenaID)
	kv(2, "arenaName", g.Arena.ArenaName)
	kv(2, "arenaCity", g.Arena.ArenaCity)
	kv(2, "arenaState", g.Arena.ArenaState)
	kv(2, "arenaCountry", g.Arena.ArenaCountry)
	kv(2, "arenaTimezone", g.Arena.ArenaTimezone)

	fmt.Printf("  officials: %d\n", len(g.Officials))
	for _, o := range g.Officials {
		kv(2, "personId", o.PersonID)
		kv(2, "name", o.Name)
		kv(2, "nameI", o.NameI)
		kv(2, "firstName", o.FirstName)
		kv(2, "familyName", o.FamilyName)
		kv(2, "jerseyNum", o.JerseyNum)
		kv(2, "assignment", o.Assignment)
	}

	fmt.Println("  awayTeam:")
	printBoxScoreTeam(g.AwayTeam)
	fmt.Println("  homeTeam:")
	printBoxScoreTeam(g.HomeTeam)
}

func printBoxScoreTeam(t live.BoxScoreTeam) {
	kv(2, "teamId", t.TeamID)
	kv(2, "teamName", t.TeamName)
	kv(2, "teamCity", t.TeamCity)
	kv(2, "teamTricode", t.TeamTricode)
	kv(2, "score", t.Score)
	kv(2, "inBonus", t.InBonus)
	kv(2, "timeoutsRemaining", t.TimeoutsRemaining)
	kv(2, "periods", fmt.Sprintf("%d", len(t.Periods)))
	for _, p := range t.Periods {
		kv(3, fmt.Sprintf("P%d (%s)", p.Period, p.PeriodType), p.Score)
	}

	fmt.Printf("    statistics (team, %d fields):\n", 60)
	printLiveTeamStats(t.Statistics)

	fmt.Printf("    players: %d\n", len(t.Players))
	if len(t.Players) > 0 {
		fmt.Println("    players[0] (full detail, every field):")
		printBoxScorePlayer(t.Players[0])
		if len(t.Players) > 1 {
			fmt.Printf("    (... %d more players, same shape ...)\n", len(t.Players)-1)
		}
	}
}

func printBoxScorePlayer(p live.BoxScorePlayer) {
	kv(3, "status", p.Status)
	kv(3, "order", p.Order)
	kv(3, "personId", p.PersonID)
	kv(3, "jerseyNum", p.JerseyNum)
	kv(3, "position", p.Position)
	kv(3, "starter", p.Starter) // string "1"/"0"
	kv(3, "oncourt", p.OnCourt) // string "1"/"0"
	kv(3, "played", p.Played)   // string "1"/"0"
	kv(3, "name", p.Name)
	kv(3, "nameI", p.NameI)
	kv(3, "firstName", p.FirstName)
	kv(3, "familyName", p.FamilyName)
	kv(3, "notPlayingReason", p.NotPlayingReason)
	kv(3, "notPlayingDescription", p.NotPlayingDescription)
	fmt.Println("      statistics:")
	printLivePlayerStats(p.Statistics)
}

func printLivePlayerStats(s live.BoxScorePlayerStats) {
	kv(4, "minutes", s.Minutes)
	kv(4, "minutesCalculated", s.MinutesCalculated)
	kv(4, "points", s.Points)
	kv(4, "assists", s.Assists)
	kv(4, "blocks", s.Blocks)
	kv(4, "blocksReceived", s.BlocksReceived)
	kv(4, "fieldGoalsMade", s.FieldGoalsMade)
	kv(4, "fieldGoalsAttempted", s.FieldGoalsAttempted)
	kv(4, "fieldGoalsPercentage", s.FieldGoalsPercentage)
	kv(4, "threePointersMade", s.ThreePointersMade)
	kv(4, "threePointersAttempted", s.ThreePointersAttempted)
	kv(4, "threePointersPercentage", s.ThreePointersPercentage)
	kv(4, "twoPointersMade", s.TwoPointersMade)
	kv(4, "twoPointersAttempted", s.TwoPointersAttempted)
	kv(4, "twoPointersPercentage", s.TwoPointersPercentage)
	kv(4, "freeThrowsMade", s.FreeThrowsMade)
	kv(4, "freeThrowsAttempted", s.FreeThrowsAttempted)
	kv(4, "freeThrowsPercentage", s.FreeThrowsPercentage)
	kv(4, "reboundsOffensive", s.ReboundsOffensive)
	kv(4, "reboundsDefensive", s.ReboundsDefensive)
	kv(4, "reboundsTotal", s.ReboundsTotal)
	kv(4, "steals", s.Steals)
	kv(4, "turnovers", s.Turnovers)
	kv(4, "foulsPersonal", s.FoulsPersonal)
	kv(4, "foulsOffensive", s.FoulsOffensive)
	kv(4, "foulsDrawn", s.FoulsDrawn)
	kv(4, "foulsTechnical", s.FoulsTechnical)
	kv(4, "plus", s.Plus)
	kv(4, "minus", s.Minus)
	kv(4, "plusMinusPoints", s.PlusMinusPoints)
	kv(4, "pointsFastBreak", s.PointsFastBreak)
	kv(4, "pointsInThePaint", s.PointsInThePaint)
	kv(4, "pointsSecondChance", s.PointsSecondChance)
}

func printLiveTeamStats(s live.BoxScoreTeamStats) {
	kv(3, "minutes", s.Minutes)
	kv(3, "minutesCalculated", s.MinutesCalculated)
	kv(3, "points", s.Points)
	kv(3, "pointsAgainst", s.PointsAgainst)
	kv(3, "assists", s.Assists)
	kv(3, "assistsTurnoverRatio", s.AssistsTurnoverRatio)
	kv(3, "benchPoints", s.BenchPoints)
	kv(3, "biggestLead", s.BiggestLead)
	kv(3, "biggestLeadScore", s.BiggestLeadScore)
	kv(3, "biggestScoringRun", s.BiggestScoringRun)
	kv(3, "biggestScoringRunScore", s.BiggestScoringRunScore)
	kv(3, "blocks", s.Blocks)
	kv(3, "blocksReceived", s.BlocksReceived)
	kv(3, "fastBreakPointsMade", s.FastBreakPointsMade)
	kv(3, "fastBreakPointsAttempted", s.FastBreakPointsAttempted)
	kv(3, "fastBreakPointsPercentage", s.FastBreakPointsPercentage)
	kv(3, "fieldGoalsMade", s.FieldGoalsMade)
	kv(3, "fieldGoalsAttempted", s.FieldGoalsAttempted)
	kv(3, "fieldGoalsPercentage", s.FieldGoalsPercentage)
	kv(3, "fieldGoalsEffectiveAdjusted", s.FieldGoalsEffectiveAdjusted)
	kv(3, "foulsOffensive", s.FoulsOffensive)
	kv(3, "foulsDrawn", s.FoulsDrawn)
	kv(3, "foulsPersonal", s.FoulsPersonal)
	kv(3, "foulsTeam", s.FoulsTeam)
	kv(3, "foulsTechnical", s.FoulsTechnical)
	kv(3, "foulsTeamTechnical", s.FoulsTeamTechnical)
	kv(3, "freeThrowsMade", s.FreeThrowsMade)
	kv(3, "freeThrowsAttempted", s.FreeThrowsAttempted)
	kv(3, "freeThrowsPercentage", s.FreeThrowsPercentage)
	kv(3, "leadChanges", s.LeadChanges)
	kv(3, "pointsFastBreak", s.PointsFastBreak)
	kv(3, "pointsFromTurnovers", s.PointsFromTurnovers)
	kv(3, "pointsInThePaint", s.PointsInThePaint)
	kv(3, "pointsInThePaintAttempted", s.PointsInThePaintAttempted)
	kv(3, "pointsInThePaintMade", s.PointsInThePaintMade)
	kv(3, "pointsInThePaintPercentage", s.PointsInThePaintPercentage)
	kv(3, "pointsSecondChance", s.PointsSecondChance)
	kv(3, "reboundsOffensive", s.ReboundsOffensive)
	kv(3, "reboundsDefensive", s.ReboundsDefensive)
	kv(3, "reboundsPersonal", s.ReboundsPersonal)
	kv(3, "reboundsTeam", s.ReboundsTeam)
	kv(3, "reboundsTeamDefensive", s.ReboundsTeamDefensive)
	kv(3, "reboundsTeamOffensive", s.ReboundsTeamOffensive)
	kv(3, "reboundsTotal", s.ReboundsTotal)
	kv(3, "secondChancePointsMade", s.SecondChancePointsMade)
	kv(3, "secondChancePointsAttempted", s.SecondChancePointsAttempted)
	kv(3, "secondChancePointsPercentage", s.SecondChancePointsPercentage)
	kv(3, "steals", s.Steals)
	kv(3, "threePointersMade", s.ThreePointersMade)
	kv(3, "threePointersAttempted", s.ThreePointersAttempted)
	kv(3, "threePointersPercentage", s.ThreePointersPercentage)
	kv(3, "timeLeading", s.TimeLeading)
	kv(3, "timesTied", s.TimesTied)
	kv(3, "trueShootingAttempts", s.TrueShootingAttempts)
	kv(3, "trueShootingPercentage", s.TrueShootingPercentage)
	kv(3, "turnovers", s.Turnovers)
	kv(3, "turnoversTeam", s.TurnoversTeam)
	kv(3, "turnoversTotal", s.TurnoversTotal)
	kv(3, "twoPointersMade", s.TwoPointersMade)
	kv(3, "twoPointersAttempted", s.TwoPointersAttempted)
	kv(3, "twoPointersPercentage", s.TwoPointersPercentage)
}

// ===================================================================
// LIVE: PlayByPlay
// ===================================================================

func printLivePlayByPlay(r *live.PlayByPlayResponse) {
	printMetaLive(r.Meta)
	fmt.Println("game:")
	kv(1, "gameId", r.Game.GameID)
	kv(1, "actions", fmt.Sprintf("%d total", len(r.Game.Actions)))
	if len(r.Game.Actions) == 0 {
		return
	}
	fmt.Println("  actions[last] (full detail, every field):")
	printLiveAction(r.Game.Actions[len(r.Game.Actions)-1])
	fmt.Println("  (most fields are omitempty — present only for the action types that produce them)")
}

func printLiveAction(a live.Action) {
	kv(2, "actionNumber", a.ActionNumber)
	kv(2, "clock", a.Clock)
	kv(2, "timeActual", a.TimeActual)
	kv(2, "period", a.Period)
	kv(2, "periodType", a.PeriodType)
	kv(2, "teamId", a.TeamID)
	kv(2, "teamTricode", a.TeamTricode)
	kv(2, "actionType", a.ActionType)
	kv(2, "subType", a.SubType)
	kv(2, "descriptor", a.Descriptor)
	kv(2, "qualifiers", strings.Join(a.Qualifiers, ","))
	kv(2, "personId", a.PersonID)
	kv(2, "playerName", a.PlayerName)
	kv(2, "playerNameI", a.PlayerNameI)
	kv(2, "x", a.X)
	kv(2, "y", a.Y)
	kv(2, "xLegacy", a.XLegacy)
	kv(2, "yLegacy", a.YLegacy)
	kv(2, "side", a.Side)
	kv(2, "shotDistance", a.ShotDistance)
	kv(2, "shotResult", a.ShotResult)
	kv(2, "isFieldGoal", a.IsFieldGoal)
	kv(2, "scoreHome", a.ScoreHome)
	kv(2, "scoreAway", a.ScoreAway)
	kv(2, "pointsTotal", a.PointsTotal)
	kv(2, "possession", a.Possession)
	kv(2, "description", a.Description)
	kv(2, "orderNumber", a.OrderNumber)
	kv(2, "edited", a.Edited)
	kv(2, "videoAvailable", a.VideoAvailable)
	kv(2, "assistPersonId", a.AssistPersonID)
	kv(2, "assistPlayerNameInitial", a.AssistPlayerNameI)
	kv(2, "assistTotal", a.AssistTotal)
	kv(2, "blockPersonId", a.BlockPersonID)
	kv(2, "blockPlayerName", a.BlockPlayerName)
	kv(2, "stealPersonId", a.StealPersonID)
	kv(2, "stealPlayerName", a.StealPlayerName)
	kv(2, "foulPersonalTotal", a.FoulPersonalTotal)
	kv(2, "foulTechnicalTotal", a.FoulTechnicalTotal)
	kv(2, "foulDrawnPersonId", a.FoulDrawnPersonID)
	kv(2, "foulDrawnPlayerName", a.FoulDrawnPlayerName)
	kv(2, "reboundTotal", a.ReboundTotal)
	kv(2, "reboundDefensiveTotal", a.ReboundDefensiveTotal)
	kv(2, "reboundOffensiveTotal", a.ReboundOffensiveTotal)
	kv(2, "turnoverTotal", a.TurnoverTotal)
	kv(2, "officialId", a.OfficialID)
	kv(2, "personIdsFilter", fmt.Sprintf("%v", a.PersonIDsFilter))
}

// ===================================================================
// STATS: ScoreboardV3
// ===================================================================

func printScoreboardV3(r *stats.ScoreboardV3Response) {
	printMetaStats(r.Meta)
	s := r.Scoreboard
	fmt.Println("scoreboard:")
	kv(1, "gameDate", s.GameDate)
	kv(1, "leagueId", s.LeagueID)
	kv(1, "leagueName", s.LeagueName)
	kv(1, "games", fmt.Sprintf("%d total", len(s.Games)))

	for i := range s.Games {
		g := &s.Games[i]
		fmt.Printf("  game[%d]:\n", i)
		kv(2, "gameId", g.GameID)
		kv(2, "gameCode", g.GameCode)
		kv(2, "gameStatus", g.GameStatus)
		kv(2, "gameStatusText", g.GameStatusText)
		kv(2, "period", g.Period)
		kv(2, "gameClock", g.GameClock)
		kv(2, "gameTimeUTC", g.GameTimeUTC)
		kv(2, "gameEt", g.GameEt)
		kv(2, "regulationPeriods", g.RegulationPeriods)
		kv(2, "seriesGameNumber", g.SeriesGameNumber)
		kv(2, "seriesText", g.SeriesText)
		kv(2, "seriesConference", g.SeriesConference)
		kv(2, "poRoundDesc", g.PoRoundDesc)
		kv(2, "gameSubtype", g.GameSubtype)
		kv(2, "ifNecessary", g.IfNecessary)
		fmt.Println("    awayTeam:")
		printScoreboardV3Team(g.AwayTeam)
		fmt.Println("    homeTeam:")
		printScoreboardV3Team(g.HomeTeam)
		fmt.Println("    gameLeaders:")
		fmt.Println("      awayLeaders:")
		printScoreboardV3Leader(g.GameLeaders.AwayLeaders)
		fmt.Println("      homeLeaders:")
		printScoreboardV3Leader(g.GameLeaders.HomeLeaders)
		fmt.Println("    broadcasters:")
		printScoreboardV3Broadcasters(g.Broadcasters)
		if i == 0 && len(s.Games) > 1 {
			fmt.Printf("  (... %d more games, same shape ...)\n", len(s.Games)-1)
			break
		}
	}
}

func printScoreboardV3Team(t stats.ScoreboardV3Team) {
	kv(3, "teamId", t.TeamID)
	kv(3, "teamName", t.TeamName)
	kv(3, "teamCity", t.TeamCity)
	kv(3, "teamTricode", t.TeamTricode)
	kv(3, "teamSlug", t.TeamSlug)
	kv(3, "wins", t.Wins)
	kv(3, "losses", t.Losses)
	kv(3, "score", t.Score)
	kv(3, "seed", t.Seed)
	kv(3, "inBonus", t.InBonus)
	kv(3, "timeoutsRemaining", t.TimeoutsRemaining)
	kv(3, "periods", fmt.Sprintf("%d", len(t.Periods)))
	for _, p := range t.Periods {
		kv(4, fmt.Sprintf("P%d (%s)", p.Period, p.PeriodType), p.Score)
	}
}

func printScoreboardV3Leader(l stats.ScoreboardV3Leader) {
	kv(4, "personId", l.PersonID)
	kv(4, "name", l.Name)
	kv(4, "jerseyNum", l.JerseyNum)
	kv(4, "position", l.Position)
	kv(4, "teamTricode", l.TeamTricode)
	kv(4, "playerSlug", l.PlayerSlug)
	kv(4, "points", l.Points)
	kv(4, "rebounds", l.Rebounds)
	kv(4, "assists", l.Assists)
}

func printScoreboardV3Broadcasters(b stats.ScoreboardV3Broadcasters) {
	printBroadcasterList(3, "nationalTv", b.NationalTvBroadcasters)
	printBroadcasterList(3, "nationalRadio", b.NationalRadioBroadcasters)
	printBroadcasterList(3, "homeTv", b.HomeTvBroadcasters)
	printBroadcasterList(3, "homeRadio", b.HomeRadioBroadcasters)
	printBroadcasterList(3, "awayTv", b.AwayTvBroadcasters)
	printBroadcasterList(3, "awayRadio", b.AwayRadioBroadcasters)
}

func printBroadcasterList(indent int, label string, list []stats.ScoreboardV3Broadcaster) {
	if len(list) == 0 {
		return
	}
	names := make([]string, 0, len(list))
	for _, b := range list {
		names = append(names, b.BroadcasterDisplay)
	}
	kv(indent, label, strings.Join(names, ", "))
}

// ===================================================================
// STATS: BoxScoreTraditionalV3
// ===================================================================

func printBoxScoreTraditionalV3(r *stats.BoxScoreTraditionalV3Response) {
	printMetaStats(r.Meta)
	b := r.BoxScoreTraditional
	fmt.Println("boxScoreTraditional:")
	kv(1, "gameId", b.GameID)
	kv(1, "awayTeamId", b.AwayTeamID)
	kv(1, "homeTeamId", b.HomeTeamID)
	fmt.Println("  awayTeam:")
	printBoxTeam(b.AwayTeam)
	fmt.Println("  homeTeam:")
	printBoxTeam(b.HomeTeam)
}

func printBoxTeam(t stats.BoxTeam) {
	kv(2, "teamId", t.TeamID)
	kv(2, "teamCity", t.TeamCity)
	kv(2, "teamName", t.TeamName)
	kv(2, "teamTricode", t.TeamTricode)
	kv(2, "teamSlug", t.TeamSlug)

	fmt.Println("    statistics (team totals):")
	printTradTeamStats(t.Statistics)
	fmt.Println("    starters (subtotal):")
	printTradTeamStats(t.Starters)
	fmt.Println("    bench (subtotal):")
	printTradTeamStats(t.Bench)

	fmt.Printf("    players: %d\n", len(t.Players))
	if len(t.Players) > 0 {
		fmt.Println("    players[0] (full detail, every field):")
		printBoxPlayer(t.Players[0])
		if len(t.Players) > 1 {
			fmt.Printf("    (... %d more players, same shape ...)\n", len(t.Players)-1)
		}
	}
}

func printBoxPlayer(p stats.BoxPlayer) {
	kv(3, "personId", p.PersonID)
	kv(3, "firstName", p.FirstName)
	kv(3, "familyName", p.FamilyName)
	kv(3, "nameI", p.NameI)
	kv(3, "playerSlug", p.PlayerSlug)
	kv(3, "position", p.Position)
	kv(3, "comment", p.Comment)
	kv(3, "jerseyNum", p.JerseyNum)
	kv(3, "starter", p.Starter)
	fmt.Println("      statistics:")
	printTradPlayerStats(p.Statistics)
}

func printTradPlayerStats(s stats.BoxPlayerStats) {
	kv(4, "minutes", s.Minutes)
	kv(4, "fieldGoalsMade", s.FieldGoalsMade)
	kv(4, "fieldGoalsAttempted", s.FieldGoalsAttempted)
	kv(4, "fieldGoalsPercentage", s.FieldGoalsPercentage)
	kv(4, "threePointersMade", s.ThreePointersMade)
	kv(4, "threePointersAttempted", s.ThreePointersAttempted)
	kv(4, "threePointersPercentage", s.ThreePointersPercentage)
	kv(4, "freeThrowsMade", s.FreeThrowsMade)
	kv(4, "freeThrowsAttempted", s.FreeThrowsAttempted)
	kv(4, "freeThrowsPercentage", s.FreeThrowsPercentage)
	kv(4, "reboundsOffensive", s.ReboundsOffensive)
	kv(4, "reboundsDefensive", s.ReboundsDefensive)
	kv(4, "reboundsTotal", s.ReboundsTotal)
	kv(4, "assists", s.Assists)
	kv(4, "steals", s.Steals)
	kv(4, "blocks", s.Blocks)
	kv(4, "turnovers", s.Turnovers)
	kv(4, "foulsPersonal", s.FoulsPersonal)
	kv(4, "points", s.Points)
	kv(4, "plusMinusPoints", s.PlusMinusPoints)
}

func printTradTeamStats(s stats.BoxTeamStats) {
	kv(3, "minutes", s.Minutes)
	kv(3, "fieldGoalsMade", s.FieldGoalsMade)
	kv(3, "fieldGoalsAttempted", s.FieldGoalsAttempted)
	kv(3, "fieldGoalsPercentage", s.FieldGoalsPercentage)
	kv(3, "threePointersMade", s.ThreePointersMade)
	kv(3, "threePointersAttempted", s.ThreePointersAttempted)
	kv(3, "threePointersPercentage", s.ThreePointersPercentage)
	kv(3, "freeThrowsMade", s.FreeThrowsMade)
	kv(3, "freeThrowsAttempted", s.FreeThrowsAttempted)
	kv(3, "freeThrowsPercentage", s.FreeThrowsPercentage)
	kv(3, "reboundsOffensive", s.ReboundsOffensive)
	kv(3, "reboundsDefensive", s.ReboundsDefensive)
	kv(3, "reboundsTotal", s.ReboundsTotal)
	kv(3, "assists", s.Assists)
	kv(3, "steals", s.Steals)
	kv(3, "blocks", s.Blocks)
	kv(3, "turnovers", s.Turnovers)
	kv(3, "foulsPersonal", s.FoulsPersonal)
	kv(3, "points", s.Points)
	kv(3, "plusMinusPoints", s.PlusMinusPoints)
}

// ===================================================================
// STATS: PlayByPlayV3
// ===================================================================

func printPlayByPlayV3(r *stats.PlayByPlayV3Response) {
	printMetaStats(r.Meta)
	fmt.Println("game:")
	kv(1, "gameId", r.Game.GameID)
	kv(1, "actions", fmt.Sprintf("%d total", len(r.Game.Actions)))
	if len(r.Game.Actions) == 0 {
		return
	}
	fmt.Println("  actions[0] (full detail, every field):")
	printStatsAction(r.Game.Actions[0])

	// Demonstrate the documented "filter by period yourself" pattern.
	q1 := 0
	for _, a := range r.Game.Actions {
		if a.Period == 1 {
			q1++
		}
	}
	fmt.Printf("  example filter — actions in Q1: %d\n", q1)
}

func printStatsAction(a stats.Action) {
	kv(2, "actionNumber", a.ActionNumber)
	kv(2, "actionId", a.ActionID)
	kv(2, "clock", a.Clock)
	kv(2, "timeActual", a.TimeActual)
	kv(2, "period", a.Period)
	kv(2, "periodType", a.PeriodType)
	kv(2, "teamId", a.TeamID)
	kv(2, "teamTricode", a.TeamTricode)
	kv(2, "personId", a.PersonID)
	kv(2, "playerName", a.PlayerName)
	kv(2, "playerNameI", a.PlayerNameI)
	kv(2, "x", a.XLegacyValue)
	kv(2, "y", a.YLegacyValue)
	kv(2, "xLegacy", a.XLegacy)
	kv(2, "yLegacy", a.YLegacy)
	kv(2, "side", a.Side)
	kv(2, "shotDistance", a.ShotDistance)
	kv(2, "shotResult", a.ShotResult)
	kv(2, "isFieldGoal", a.IsFieldGoal)
	kv(2, "scoreHome", a.ScoreHome)
	kv(2, "scoreAway", a.ScoreAway)
	kv(2, "pointsTotal", a.PointsTotal)
	kv(2, "location", a.Location)
	kv(2, "possession", a.Possession)
	kv(2, "description", a.Description)
	kv(2, "actionType", a.ActionType)
	kv(2, "subType", a.SubType)
	kv(2, "descriptor", a.Descriptor)
	kv(2, "qualifiers", strings.Join(a.Qualifiers, ","))
	kv(2, "orderNumber", a.OrderNumber)
	kv(2, "videoAvailable", a.VideoAvailable)
}

// ===================================================================
// helpers
// ===================================================================

func printMetaLive(m live.Meta) {
	fmt.Println("meta:")
	kv(1, "version", m.Version)
	kv(1, "request", m.Request)
	kv(1, "time", m.Time)
	kv(1, "code", m.Code)
}

func printMetaStats(m stats.Meta) {
	fmt.Println("meta:")
	kv(1, "version", m.Version)
	kv(1, "request", m.Request)
	kv(1, "time", m.Time)
	kv(1, "code", m.Code)
}

func firstInProgress(games []live.Game) *live.Game {
	for i := range games {
		if games[i].GameStatus == 2 { // 1=scheduled, 2=live, 3=final
			return &games[i]
		}
	}
	return nil
}

// kv prints a "key: value" line at the given indent depth (2 spaces each).
func kv(indent int, key string, val any) {
	fmt.Printf("%s%s: %v\n", strings.Repeat("  ", indent+1), key, val)
}

// header prints a section banner.
func header(title string) {
	bar := strings.Repeat("=", len(title)+4)
	fmt.Printf("\n%s\n  %s\n%s\n", bar, title, bar)
}

// dumpJSON is an alternative to the field-by-field printers above: it marshals
// any response struct to indented JSON, which also reveals every populated
// value. Kept here as a documented convenience for callers who'd rather not
// hand-write printers. Unused by main; reference it from your own code.
func dumpJSON(label string, v any) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("%s: <marshal error: %v>\n", label, err)
		return
	}
	fmt.Printf("%s:\n%s\n", label, b)
}
