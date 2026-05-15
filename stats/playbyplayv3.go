package stats

import (
	"context"
	"net/url"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
)

// PlayByPlayV3Response mirrors the full JSON response from
// stats.nba.com/stats/playbyplayv3.
type PlayByPlayV3Response struct {
	Meta Meta            `json:"meta"`
	Game PlayByPlayGame  `json:"game"`
}

// PlayByPlayGame is the game-level container for play-by-play data.
type PlayByPlayGame struct {
	GameID  string   `json:"gameId"`
	Actions []Action `json:"actions"`
}

// Action is a single play-by-play event. Each event carries a Period field,
// so callers can filter by quarter on their own.
type Action struct {
	ActionNumber       int     `json:"actionNumber"`
	Clock              string  `json:"clock"`
	TimeActual         string  `json:"timeActual,omitempty"`
	Period             int     `json:"period"`
	PeriodType         string  `json:"periodType,omitempty"`
	TeamID             int     `json:"teamId,omitempty"`
	TeamTricode        string  `json:"teamTricode,omitempty"`
	PersonID           int     `json:"personId,omitempty"`
	PlayerName         string  `json:"playerName,omitempty"`
	PlayerNameI        string  `json:"playerNameI,omitempty"`
	XLegacy            int     `json:"xLegacy,omitempty"`
	YLegacy            int     `json:"yLegacy,omitempty"`
	ShotDistance       float64 `json:"shotDistance,omitempty"`
	ShotResult         string  `json:"shotResult,omitempty"`
	IsFieldGoal        int     `json:"isFieldGoal,omitempty"`
	ScoreHome          string  `json:"scoreHome,omitempty"`
	ScoreAway          string  `json:"scoreAway,omitempty"`
	PointsTotal        int     `json:"pointsTotal,omitempty"`
	Location           string  `json:"location,omitempty"`
	Description        string  `json:"description"`
	ActionType         string  `json:"actionType"`
	SubType            string  `json:"subType,omitempty"`
	Descriptor         string  `json:"descriptor,omitempty"`
	Qualifiers         []string `json:"qualifiers,omitempty"`
	Possession         int     `json:"possession,omitempty"`
	OrderNumber        int     `json:"orderNumber,omitempty"`
	XLegacyValue       int     `json:"x,omitempty"`
	YLegacyValue       int     `json:"y,omitempty"`
	Side               string  `json:"side,omitempty"`
	VideoAvailable     int     `json:"videoAvailable,omitempty"`
	ActionID           int     `json:"actionId,omitempty"`
}

// PlayByPlayV3 fetches every play-by-play action for a single game.
// gameID is the 10-character NBA game identifier.
func (c *Client) PlayByPlayV3(ctx context.Context, gameID string) (*PlayByPlayV3Response, error) {
	q := url.Values{}
	q.Set("GameID", gameID)
	q.Set("LeagueID", "00")
	q.Set("endPeriod", "0")
	q.Set("endRange", "28800")
	q.Set("rangeType", "2")
	q.Set("startPeriod", "0")
	q.Set("startRange", "0")

	var out PlayByPlayV3Response
	if err := c.http.GetJSON(ctx, httpx.ProfileStats, c.baseURL+"/stats/playbyplayv3", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
