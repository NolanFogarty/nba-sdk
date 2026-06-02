package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
)

// ScoreboardV2Response mirrors the JSON response from
// stats.nba.com/stats/scoreboardv2. Unlike the v3 endpoints, scoreboardv2
// returns the classic "resultSets" shape (named tables of headers + row
// arrays). The SDK parses the GameHeader table into typed ScoreboardGame
// values so callers can resolve game IDs for a date without decoding the
// positional rows themselves. The untouched result sets are preserved in
// ResultSets for callers who need LineScore, SeriesStandings, etc.
type ScoreboardV2Response struct {
	Resource   string                 `json:"resource"`
	Parameters ScoreboardV2Parameters `json:"parameters"`
	ResultSets []ResultSet            `json:"resultSets"`

	// Games is derived from the GameHeader result set. It is the primary way
	// to enumerate games (and their IDs) scheduled for the requested date.
	Games []ScoreboardGame `json:"-"`
}

// ScoreboardV2Parameters echoes the query parameters the endpoint resolved.
type ScoreboardV2Parameters struct {
	GameDate  string `json:"GameDate"`
	LeagueID  string `json:"LeagueID"`
	DayOffset string `json:"DayOffset"`
}

// ResultSet is one named table in a classic stats.nba.com response. Each row
// in RowSet is a positional list of values aligned to Headers.
type ResultSet struct {
	Name    string              `json:"name"`
	Headers []string            `json:"headers"`
	RowSet  [][]json.RawMessage `json:"rowSet"`
}

// ScoreboardGame is one game from the GameHeader result set, with the columns
// callers most commonly need to identify and look up a game.
type ScoreboardGame struct {
	GameDateEst    string `json:"gameDateEst"`
	GameSequence   int    `json:"gameSequence"`
	GameID         string `json:"gameId"`
	GameStatusID   int    `json:"gameStatusId"`
	GameStatusText string `json:"gameStatusText"`
	GameCode       string `json:"gameCode"`
	HomeTeamID     int    `json:"homeTeamId"`
	VisitorTeamID  int    `json:"visitorTeamId"`
	Season         string `json:"season"`
	LivePeriod     int    `json:"livePeriod"`
}

// UnmarshalJSON decodes the raw envelope and then projects the GameHeader
// result set into Games. Columns are matched by header name, so the SDK is
// resilient to the NBA reordering or adding columns.
func (r *ScoreboardV2Response) UnmarshalJSON(data []byte) error {
	type alias ScoreboardV2Response
	var raw alias
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*r = ScoreboardV2Response(raw)

	for _, rs := range r.ResultSets {
		if rs.Name != "GameHeader" {
			continue
		}
		idx := make(map[string]int, len(rs.Headers))
		for i, h := range rs.Headers {
			idx[h] = i
		}
		for _, row := range rs.RowSet {
			var g ScoreboardGame
			getString(row, idx, "GAME_DATE_EST", &g.GameDateEst)
			getInt(row, idx, "GAME_SEQUENCE", &g.GameSequence)
			getString(row, idx, "GAME_ID", &g.GameID)
			getInt(row, idx, "GAME_STATUS_ID", &g.GameStatusID)
			getString(row, idx, "GAME_STATUS_TEXT", &g.GameStatusText)
			getString(row, idx, "GAMECODE", &g.GameCode)
			getInt(row, idx, "HOME_TEAM_ID", &g.HomeTeamID)
			getInt(row, idx, "VISITOR_TEAM_ID", &g.VisitorTeamID)
			getString(row, idx, "SEASON", &g.Season)
			getInt(row, idx, "LIVE_PERIOD", &g.LivePeriod)
			r.Games = append(r.Games, g)
		}
		break
	}
	return nil
}

// cell returns the raw value for column name in row, or nil if the column is
// absent or null.
func cell(row []json.RawMessage, idx map[string]int, name string) json.RawMessage {
	i, ok := idx[name]
	if !ok || i >= len(row) {
		return nil
	}
	v := row[i]
	if len(v) == 0 || string(v) == "null" {
		return nil
	}
	return v
}

func getString(row []json.RawMessage, idx map[string]int, name string, dst *string) {
	if v := cell(row, idx, name); v != nil {
		_ = json.Unmarshal(v, dst)
	}
}

func getInt(row []json.RawMessage, idx map[string]int, name string, dst *int) {
	if v := cell(row, idx, name); v != nil {
		_ = json.Unmarshal(v, dst)
	}
}

// ScoreboardV2 lists every game scheduled on a given date, including the game
// IDs needed to fetch historical box scores and play-by-play. date must be in
// "YYYY-MM-DD" form (e.g. "2024-12-25"). Unlike Live.Scoreboard, which is
// fixed to today, this works for any past or future date.
func (c *Client) ScoreboardV2(ctx context.Context, date string) (*ScoreboardV2Response, error) {
	if date == "" {
		return nil, fmt.Errorf("stats: ScoreboardV2 requires a date in YYYY-MM-DD form")
	}
	q := url.Values{}
	q.Set("GameDate", date)
	q.Set("LeagueID", "00")
	q.Set("DayOffset", "0")

	var out ScoreboardV2Response
	if err := c.http.GetJSON(ctx, httpx.ProfileStats, c.baseURL+"/stats/scoreboardv2", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
