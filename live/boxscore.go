package live

import (
	"context"
	"fmt"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
)

// BoxScoreResponse mirrors the JSON response from
// cdn.nba.com/static/json/liveData/boxscore/boxscore_<gameId>.json.
//
// The CDN box score is richer than stats.nba.com's traditional box score:
// it includes arena, officials, every advanced stat per player, and the
// current period/clock for in-progress games. It's the preferred source
// for live or recently finished games.
type BoxScoreResponse struct {
	Meta Meta         `json:"meta"`
	Game BoxScoreGame `json:"game"`
}

// BoxScoreGame is the game-level container for live box score data.
type BoxScoreGame struct {
	GameID            string         `json:"gameId"`
	GameTimeLocal     string         `json:"gameTimeLocal"`
	GameTimeUTC       string         `json:"gameTimeUTC"`
	GameTimeHome      string         `json:"gameTimeHome"`
	GameTimeAway      string         `json:"gameTimeAway"`
	GameEt            string         `json:"gameEt"`
	Duration          int            `json:"duration"`
	GameCode          string         `json:"gameCode"`
	GameStatusText    string         `json:"gameStatusText"`
	GameStatus        int            `json:"gameStatus"`
	RegulationPeriods int            `json:"regulationPeriods"`
	Period            int            `json:"period"`
	GameClock         string         `json:"gameClock"`
	Attendance        int            `json:"attendance"`
	Sellout           string         `json:"sellout"`
	Arena             Arena          `json:"arena"`
	Officials         []Official     `json:"officials"`
	HomeTeam          BoxScoreTeam   `json:"homeTeam"`
	AwayTeam          BoxScoreTeam   `json:"awayTeam"`
}

// Arena describes the venue where the game is being played.
type Arena struct {
	ArenaID       int    `json:"arenaId"`
	ArenaName     string `json:"arenaName"`
	ArenaCity     string `json:"arenaCity"`
	ArenaState    string `json:"arenaState"`
	ArenaCountry  string `json:"arenaCountry"`
	ArenaTimezone string `json:"arenaTimezone"`
}

// Official is a single referee assigned to the game.
type Official struct {
	PersonID    int    `json:"personId"`
	Name        string `json:"name"`
	NameI       string `json:"nameI"`
	FirstName   string `json:"firstName"`
	FamilyName  string `json:"familyName"`
	JerseyNum   string `json:"jerseyNum"`
	Assignment  string `json:"assignment"`
}

// BoxScoreTeam is one team's full live box score (players + advanced totals).
type BoxScoreTeam struct {
	TeamID            int                `json:"teamId"`
	TeamName          string             `json:"teamName"`
	TeamCity          string             `json:"teamCity"`
	TeamTricode       string             `json:"teamTricode"`
	Score             int                `json:"score"`
	InBonus           string             `json:"inBonus"`
	TimeoutsRemaining int                `json:"timeoutsRemaining"`
	Periods           []GamePeriod       `json:"periods"`
	Players           []BoxScorePlayer   `json:"players"`
	Statistics        BoxScoreTeamStats  `json:"statistics"`
}

// BoxScorePlayer is a single player row in a live box score. Note that
// `Starter`, `OnCourt`, and `Played` are returned as strings ("1"/"0")
// by the NBA, not booleans.
type BoxScorePlayer struct {
	Status        string               `json:"status"`
	Order         int                  `json:"order"`
	PersonID      int                  `json:"personId"`
	JerseyNum     string               `json:"jerseyNum"`
	Position      string               `json:"position,omitempty"`
	Starter       string               `json:"starter"`
	OnCourt       string               `json:"oncourt"`
	Played        string               `json:"played"`
	Statistics    BoxScorePlayerStats  `json:"statistics"`
	Name          string               `json:"name"`
	NameI         string               `json:"nameI"`
	FirstName     string               `json:"firstName"`
	FamilyName    string               `json:"familyName"`
	NotPlayingReason      string       `json:"notPlayingReason,omitempty"`
	NotPlayingDescription string       `json:"notPlayingDescription,omitempty"`
}

// BoxScorePlayerStats holds a single player's full live stat line.
type BoxScorePlayerStats struct {
	Assists                 int     `json:"assists"`
	Blocks                  int     `json:"blocks"`
	BlocksReceived          int     `json:"blocksReceived"`
	FieldGoalsAttempted     int     `json:"fieldGoalsAttempted"`
	FieldGoalsMade          int     `json:"fieldGoalsMade"`
	FieldGoalsPercentage    float64 `json:"fieldGoalsPercentage"`
	FoulsOffensive          int     `json:"foulsOffensive"`
	FoulsDrawn              int     `json:"foulsDrawn"`
	FoulsPersonal           int     `json:"foulsPersonal"`
	FoulsTechnical          int     `json:"foulsTechnical"`
	FreeThrowsAttempted     int     `json:"freeThrowsAttempted"`
	FreeThrowsMade          int     `json:"freeThrowsMade"`
	FreeThrowsPercentage    float64 `json:"freeThrowsPercentage"`
	Minus                   float64 `json:"minus"`
	Minutes                 string  `json:"minutes"`
	MinutesCalculated       string  `json:"minutesCalculated"`
	Plus                    float64 `json:"plus"`
	PlusMinusPoints         float64 `json:"plusMinusPoints"`
	Points                  int     `json:"points"`
	PointsFastBreak         int     `json:"pointsFastBreak"`
	PointsInThePaint        int     `json:"pointsInThePaint"`
	PointsSecondChance      int     `json:"pointsSecondChance"`
	ReboundsDefensive       int     `json:"reboundsDefensive"`
	ReboundsOffensive       int     `json:"reboundsOffensive"`
	ReboundsTotal           int     `json:"reboundsTotal"`
	Steals                  int     `json:"steals"`
	ThreePointersAttempted  int     `json:"threePointersAttempted"`
	ThreePointersMade       int     `json:"threePointersMade"`
	ThreePointersPercentage float64 `json:"threePointersPercentage"`
	Turnovers               int     `json:"turnovers"`
	TwoPointersAttempted    int     `json:"twoPointersAttempted"`
	TwoPointersMade         int     `json:"twoPointersMade"`
	TwoPointersPercentage   float64 `json:"twoPointersPercentage"`
}

// BoxScoreTeamStats holds team-level live stats. Includes advanced metrics
// (biggest lead, lead changes, points off turnovers, etc.) not present on
// the stats.nba.com traditional box score.
type BoxScoreTeamStats struct {
	Assists                       int     `json:"assists"`
	AssistsTurnoverRatio          float64 `json:"assistsTurnoverRatio"`
	BenchPoints                   int     `json:"benchPoints"`
	BiggestLead                   int     `json:"biggestLead"`
	BiggestLeadScore              string  `json:"biggestLeadScore"`
	BiggestScoringRun             int     `json:"biggestScoringRun"`
	BiggestScoringRunScore        string  `json:"biggestScoringRunScore"`
	Blocks                        int     `json:"blocks"`
	BlocksReceived                int     `json:"blocksReceived"`
	FastBreakPointsAttempted      int     `json:"fastBreakPointsAttempted"`
	FastBreakPointsMade           int     `json:"fastBreakPointsMade"`
	FastBreakPointsPercentage     float64 `json:"fastBreakPointsPercentage"`
	FieldGoalsAttempted           int     `json:"fieldGoalsAttempted"`
	FieldGoalsEffectiveAdjusted   float64 `json:"fieldGoalsEffectiveAdjusted"`
	FieldGoalsMade                int     `json:"fieldGoalsMade"`
	FieldGoalsPercentage          float64 `json:"fieldGoalsPercentage"`
	FoulsOffensive                int     `json:"foulsOffensive"`
	FoulsDrawn                    int     `json:"foulsDrawn"`
	FoulsPersonal                 int     `json:"foulsPersonal"`
	FoulsTeam                     int     `json:"foulsTeam"`
	FoulsTechnical                int     `json:"foulsTechnical"`
	FoulsTeamTechnical            int     `json:"foulsTeamTechnical"`
	FreeThrowsAttempted           int     `json:"freeThrowsAttempted"`
	FreeThrowsMade                int     `json:"freeThrowsMade"`
	FreeThrowsPercentage          float64 `json:"freeThrowsPercentage"`
	LeadChanges                   int     `json:"leadChanges"`
	Minutes                       string  `json:"minutes"`
	MinutesCalculated             string  `json:"minutesCalculated"`
	Points                        int     `json:"points"`
	PointsAgainst                 int     `json:"pointsAgainst"`
	PointsFastBreak               int     `json:"pointsFastBreak"`
	PointsFromTurnovers           int     `json:"pointsFromTurnovers"`
	PointsInThePaint              int     `json:"pointsInThePaint"`
	PointsInThePaintAttempted     int     `json:"pointsInThePaintAttempted"`
	PointsInThePaintMade          int     `json:"pointsInThePaintMade"`
	PointsInThePaintPercentage    float64 `json:"pointsInThePaintPercentage"`
	PointsSecondChance            int     `json:"pointsSecondChance"`
	ReboundsDefensive             int     `json:"reboundsDefensive"`
	ReboundsOffensive             int     `json:"reboundsOffensive"`
	ReboundsPersonal              int     `json:"reboundsPersonal"`
	ReboundsTeam                  int     `json:"reboundsTeam"`
	ReboundsTeamDefensive         int     `json:"reboundsTeamDefensive"`
	ReboundsTeamOffensive         int     `json:"reboundsTeamOffensive"`
	ReboundsTotal                 int     `json:"reboundsTotal"`
	SecondChancePointsAttempted   int     `json:"secondChancePointsAttempted"`
	SecondChancePointsMade        int     `json:"secondChancePointsMade"`
	SecondChancePointsPercentage  float64 `json:"secondChancePointsPercentage"`
	Steals                        int     `json:"steals"`
	ThreePointersAttempted        int     `json:"threePointersAttempted"`
	ThreePointersMade             int     `json:"threePointersMade"`
	ThreePointersPercentage       float64 `json:"threePointersPercentage"`
	TimeLeading                   string  `json:"timeLeading"`
	TimesTied                     int     `json:"timesTied"`
	TrueShootingAttempts          float64 `json:"trueShootingAttempts"`
	TrueShootingPercentage        float64 `json:"trueShootingPercentage"`
	Turnovers                     int     `json:"turnovers"`
	TurnoversTeam                 int     `json:"turnoversTeam"`
	TurnoversTotal                int     `json:"turnoversTotal"`
	TwoPointersAttempted          int     `json:"twoPointersAttempted"`
	TwoPointersMade               int     `json:"twoPointersMade"`
	TwoPointersPercentage         float64 `json:"twoPointersPercentage"`
}

// BoxScore fetches the live box score for a game from cdn.nba.com.
// Works for games currently in progress or recently finished. For games
// long since concluded, use Stats.BoxScoreTraditionalV3 instead.
func (c *Client) BoxScore(ctx context.Context, gameID string) (*BoxScoreResponse, error) {
	url := fmt.Sprintf("%s/static/json/liveData/boxscore/boxscore_%s.json", c.baseURL, gameID)
	var out BoxScoreResponse
	if err := c.http.GetJSON(ctx, httpx.ProfileLive, url, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
