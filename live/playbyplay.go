package live

import (
	"context"
	"fmt"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
)

// PlayByPlayResponse mirrors the JSON response from
// cdn.nba.com/static/json/liveData/playbyplay/playbyplay_<gameId>.json.
//
// The CDN play-by-play updates in near real-time during a game and is the
// preferred source for live or recently finished games. For older games,
// use Stats.PlayByPlayV3 instead.
type PlayByPlayResponse struct {
	Meta Meta           `json:"meta"`
	Game PlayByPlayGame `json:"game"`
}

// PlayByPlayGame is the game-level container for live play-by-play data.
type PlayByPlayGame struct {
	GameID  string   `json:"gameId"`
	Actions []Action `json:"actions"`
}

// Action is a single play-by-play event. Not every field is populated for
// every action — most fields are present only for the action types that
// produce them (e.g. shot fields only on Made Shot / Missed Shot actions).
type Action struct {
	ActionNumber           int      `json:"actionNumber"`
	Clock                  string   `json:"clock"`
	TimeActual             string   `json:"timeActual,omitempty"`
	Period                 int      `json:"period"`
	PeriodType             string   `json:"periodType,omitempty"`
	TeamID                 int      `json:"teamId,omitempty"`
	TeamTricode            string   `json:"teamTricode,omitempty"`
	ActionType             string   `json:"actionType"`
	SubType                string   `json:"subType,omitempty"`
	Descriptor             string   `json:"descriptor,omitempty"`
	Qualifiers             []string `json:"qualifiers,omitempty"`
	PersonID               int      `json:"personId,omitempty"`
	PlayerName             string   `json:"playerName,omitempty"`
	PlayerNameI            string   `json:"playerNameI,omitempty"`
	X                      float64  `json:"x,omitempty"`
	Y                      float64  `json:"y,omitempty"`
	XLegacy                int      `json:"xLegacy,omitempty"`
	YLegacy                int      `json:"yLegacy,omitempty"`
	Side                   string   `json:"side,omitempty"`
	ShotDistance           float64  `json:"shotDistance,omitempty"`
	ShotResult             string   `json:"shotResult,omitempty"`
	IsFieldGoal            int      `json:"isFieldGoal,omitempty"`
	ScoreHome              string   `json:"scoreHome,omitempty"`
	ScoreAway              string   `json:"scoreAway,omitempty"`
	PointsTotal            int      `json:"pointsTotal,omitempty"`
	Possession             int      `json:"possession,omitempty"`
	Description            string   `json:"description"`
	OrderNumber            int      `json:"orderNumber,omitempty"`
	Edited                 string   `json:"edited,omitempty"`
	VideoAvailable         int      `json:"videoAvailable,omitempty"`
	AssistPersonID         int      `json:"assistPersonId,omitempty"`
	AssistPlayerNameI      string   `json:"assistPlayerNameInitial,omitempty"`
	AssistTotal            int      `json:"assistTotal,omitempty"`
	BlockPersonID          int      `json:"blockPersonId,omitempty"`
	BlockPlayerName        string   `json:"blockPlayerName,omitempty"`
	StealPersonID          int      `json:"stealPersonId,omitempty"`
	StealPlayerName        string   `json:"stealPlayerName,omitempty"`
	FoulPersonalTotal      int      `json:"foulPersonalTotal,omitempty"`
	FoulTechnicalTotal     int      `json:"foulTechnicalTotal,omitempty"`
	FoulDrawnPersonID      int      `json:"foulDrawnPersonId,omitempty"`
	FoulDrawnPlayerName    string   `json:"foulDrawnPlayerName,omitempty"`
	ReboundTotal           int      `json:"reboundTotal,omitempty"`
	ReboundDefensiveTotal  int      `json:"reboundDefensiveTotal,omitempty"`
	ReboundOffensiveTotal  int      `json:"reboundOffensiveTotal,omitempty"`
	TurnoverTotal          int      `json:"turnoverTotal,omitempty"`
	OfficialID             int      `json:"officialId,omitempty"`
	PersonIDsFilter        []int    `json:"personIdsFilter,omitempty"`
}

// PlayByPlay fetches the live play-by-play for a single game from cdn.nba.com.
// Works for games currently in progress or recently finished. For games
// long since concluded, use Stats.PlayByPlayV3 instead.
//
// Every action carries a Period field, so callers can filter by quarter
// (or by player, team, action type, etc.) on their own.
func (c *Client) PlayByPlay(ctx context.Context, gameID string) (*PlayByPlayResponse, error) {
	url := fmt.Sprintf("%s/static/json/liveData/playbyplay/playbyplay_%s.json", c.baseURL, gameID)
	var out PlayByPlayResponse
	if err := c.http.GetJSON(ctx, httpx.ProfileLive, url, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
