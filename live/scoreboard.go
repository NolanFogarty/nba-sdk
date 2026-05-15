package live

import (
	"context"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
)

// ScoreboardResponse mirrors the JSON response from
// cdn.nba.com/static/json/liveData/scoreboard/todaysScoreboard_00.json.
type ScoreboardResponse struct {
	Meta       Meta       `json:"meta"`
	Scoreboard Scoreboard `json:"scoreboard"`
}

// Scoreboard is the date-level container for the games scheduled today.
type Scoreboard struct {
	GameDate   string `json:"gameDate"`
	LeagueID   string `json:"leagueId"`
	LeagueName string `json:"leagueName"`
	Games      []Game `json:"games"`
}

// Game is a single scheduled or in-progress game on the scoreboard.
type Game struct {
	GameID            string      `json:"gameId"`
	GameCode          string      `json:"gameCode"`
	GameStatus        int         `json:"gameStatus"`
	GameStatusText    string      `json:"gameStatusText"`
	Period            int         `json:"period"`
	GameClock         string      `json:"gameClock"`
	GameTimeUTC       string      `json:"gameTimeUTC"`
	GameEt            string      `json:"gameEt"`
	RegulationPeriods int         `json:"regulationPeriods"`
	SeriesGameNumber  string      `json:"seriesGameNumber,omitempty"`
	SeriesText        string      `json:"seriesText,omitempty"`
	IfNecessary       bool        `json:"ifNecessary,omitempty"`
	HomeTeam          GameTeam    `json:"homeTeam"`
	AwayTeam          GameTeam    `json:"awayTeam"`
	GameLeaders       GameLeaders `json:"gameLeaders,omitempty"`
}

// GameTeam is one team's scoreboard entry for a game.
type GameTeam struct {
	TeamID            int            `json:"teamId"`
	TeamName          string         `json:"teamName"`
	TeamCity          string         `json:"teamCity"`
	TeamTricode       string         `json:"teamTricode"`
	Wins              int            `json:"wins"`
	Losses            int            `json:"losses"`
	Score             int            `json:"score"`
	Seed              int            `json:"seed,omitempty"`
	InBonus           string         `json:"inBonus,omitempty"`
	TimeoutsRemaining int            `json:"timeoutsRemaining"`
	Periods           []GamePeriod   `json:"periods"`
}

// GamePeriod is one team's score for a single period.
type GamePeriod struct {
	Period     int    `json:"period"`
	PeriodType string `json:"periodType"`
	Score      int    `json:"score"`
}

// GameLeaders holds the top scoring leaders for each team in a game.
type GameLeaders struct {
	HomeLeaders Leader `json:"homeLeaders"`
	AwayLeaders Leader `json:"awayLeaders"`
}

// Leader is a single statistical leader entry.
type Leader struct {
	PersonID    int     `json:"personId"`
	Name        string  `json:"name"`
	JerseyNum   string  `json:"jerseyNum"`
	Position    string  `json:"position"`
	TeamTricode string  `json:"teamTricode"`
	PlayerSlug  string  `json:"playerSlug"`
	Points      int     `json:"points"`
	Rebounds    int     `json:"rebounds"`
	Assists     int     `json:"assists"`
}

// Meta is the standard response envelope used by cdn.nba.com live endpoints.
type Meta struct {
	Version int    `json:"version"`
	Request string `json:"request"`
	Time    string `json:"time"`
	Code    int    `json:"code,omitempty"`
}

// Scoreboard fetches every game scheduled today. The response contains
// upcoming, in-progress, and completed games for the current NBA day.
func (c *Client) Scoreboard(ctx context.Context) (*ScoreboardResponse, error) {
	var out ScoreboardResponse
	url := c.baseURL + "/static/json/liveData/scoreboard/todaysScoreboard_00.json"
	if err := c.http.GetJSON(ctx, httpx.ProfileLive, url, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
