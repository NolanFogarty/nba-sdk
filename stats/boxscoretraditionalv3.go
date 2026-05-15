package stats

import (
	"context"
	"net/url"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
)

// BoxScoreTraditionalV3Response mirrors the full JSON response from
// stats.nba.com/stats/boxscoretraditionalv3. Every documented field is
// included; callers can ignore what they don't need.
type BoxScoreTraditionalV3Response struct {
	Meta                 Meta                 `json:"meta"`
	BoxScoreTraditional  BoxScoreTraditional  `json:"boxScoreTraditional"`
}

// BoxScoreTraditional is the game-level container for traditional box score data.
type BoxScoreTraditional struct {
	GameID     string  `json:"gameId"`
	AwayTeamID int     `json:"awayTeamId"`
	HomeTeamID int     `json:"homeTeamId"`
	HomeTeam   BoxTeam `json:"homeTeam"`
	AwayTeam   BoxTeam `json:"awayTeam"`
}

// BoxTeam is one team's traditional box score (players plus team-level totals).
type BoxTeam struct {
	TeamID      int             `json:"teamId"`
	TeamCity    string          `json:"teamCity"`
	TeamName    string          `json:"teamName"`
	TeamTricode string          `json:"teamTricode"`
	TeamSlug    string          `json:"teamSlug"`
	Players     []BoxPlayer     `json:"players"`
	Statistics  BoxTeamStats    `json:"statistics"`
	Starters    BoxTeamStats    `json:"starters"`
	Bench       BoxTeamStats    `json:"bench"`
}

// BoxPlayer is a single player row in a traditional box score.
type BoxPlayer struct {
	PersonID      int            `json:"personId"`
	FirstName     string         `json:"firstName"`
	FamilyName    string         `json:"familyName"`
	NameI         string         `json:"nameI"`
	PlayerSlug    string         `json:"playerSlug"`
	Position      string         `json:"position"`
	Comment       string         `json:"comment"`
	JerseyNum     string         `json:"jerseyNum"`
	Starter       bool           `json:"starter,omitempty"`
	Statistics    BoxPlayerStats `json:"statistics"`
}

// BoxPlayerStats holds a single player's traditional stat line for the requested period range.
type BoxPlayerStats struct {
	Minutes                 string  `json:"minutes"`
	FieldGoalsMade          int     `json:"fieldGoalsMade"`
	FieldGoalsAttempted     int     `json:"fieldGoalsAttempted"`
	FieldGoalsPercentage    float64 `json:"fieldGoalsPercentage"`
	ThreePointersMade       int     `json:"threePointersMade"`
	ThreePointersAttempted  int     `json:"threePointersAttempted"`
	ThreePointersPercentage float64 `json:"threePointersPercentage"`
	FreeThrowsMade          int     `json:"freeThrowsMade"`
	FreeThrowsAttempted     int     `json:"freeThrowsAttempted"`
	FreeThrowsPercentage    float64 `json:"freeThrowsPercentage"`
	ReboundsOffensive       int     `json:"reboundsOffensive"`
	ReboundsDefensive       int     `json:"reboundsDefensive"`
	ReboundsTotal           int     `json:"reboundsTotal"`
	Assists                 int     `json:"assists"`
	Steals                  int     `json:"steals"`
	Blocks                  int     `json:"blocks"`
	Turnovers               int     `json:"turnovers"`
	FoulsPersonal           int     `json:"foulsPersonal"`
	Points                  int     `json:"points"`
	PlusMinusPoints         float64 `json:"plusMinusPoints"`
}

// BoxTeamStats holds team-level (or starters/bench) traditional stats.
type BoxTeamStats struct {
	Minutes                 string  `json:"minutes"`
	FieldGoalsMade          int     `json:"fieldGoalsMade"`
	FieldGoalsAttempted     int     `json:"fieldGoalsAttempted"`
	FieldGoalsPercentage    float64 `json:"fieldGoalsPercentage"`
	ThreePointersMade       int     `json:"threePointersMade"`
	ThreePointersAttempted  int     `json:"threePointersAttempted"`
	ThreePointersPercentage float64 `json:"threePointersPercentage"`
	FreeThrowsMade          int     `json:"freeThrowsMade"`
	FreeThrowsAttempted     int     `json:"freeThrowsAttempted"`
	FreeThrowsPercentage    float64 `json:"freeThrowsPercentage"`
	ReboundsOffensive       int     `json:"reboundsOffensive"`
	ReboundsDefensive       int     `json:"reboundsDefensive"`
	ReboundsTotal           int     `json:"reboundsTotal"`
	Assists                 int     `json:"assists"`
	Steals                  int     `json:"steals"`
	Blocks                  int     `json:"blocks"`
	Turnovers               int     `json:"turnovers"`
	FoulsPersonal           int     `json:"foulsPersonal"`
	Points                  int     `json:"points"`
	PlusMinusPoints         float64 `json:"plusMinusPoints"`
}

// Meta is the standard response envelope used by stats.nba.com v3 endpoints.
type Meta struct {
	Version int    `json:"version"`
	Request string `json:"request"`
	Time    string `json:"time"`
	Code    int    `json:"code,omitempty"`
}

// BoxScoreTraditionalV3 fetches the traditional box score for a single game.
// gameID is the 10-character NBA game identifier (e.g. "0022400001").
func (c *Client) BoxScoreTraditionalV3(ctx context.Context, gameID string) (*BoxScoreTraditionalV3Response, error) {
	q := url.Values{}
	q.Set("GameID", gameID)
	q.Set("LeagueID", "00")
	q.Set("endPeriod", "0")
	q.Set("endRange", "28800")
	q.Set("rangeType", "0")
	q.Set("startPeriod", "0")
	q.Set("startRange", "0")

	var out BoxScoreTraditionalV3Response
	if err := c.http.GetJSON(ctx, httpx.ProfileStats, c.baseURL+"/stats/boxscoretraditionalv3", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
