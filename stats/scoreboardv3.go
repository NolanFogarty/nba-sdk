package stats

import (
	"context"
	"net/url"
	"time"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
)

// ScoreboardV3Response mirrors the JSON response from
// stats.nba.com/stats/scoreboardv3. Unlike the CDN scoreboard (which
// returns whatever NBA considers "today" and rolls over on its own
// schedule), this endpoint accepts an explicit date and returns every
// game on that calendar day — scheduled, in-progress, and final.
type ScoreboardV3Response struct {
	Meta       Meta       `json:"meta"`
	Scoreboard ScoreboardV3 `json:"scoreboard"`
}

// ScoreboardV3 is the date-level container for games on a given day.
type ScoreboardV3 struct {
	GameDate   string           `json:"gameDate"`
	LeagueID   string           `json:"leagueId"`
	LeagueName string           `json:"leagueName"`
	Games      []ScoreboardV3Game `json:"games"`
}

// ScoreboardV3Game is a single game on the scoreboard. GameStatus values:
// 1 = scheduled, 2 = in progress, 3 = final.
type ScoreboardV3Game struct {
	GameID            string                 `json:"gameId"`
	GameCode          string                 `json:"gameCode"`
	GameStatus        int                    `json:"gameStatus"`
	GameStatusText    string                 `json:"gameStatusText"`
	Period            int                    `json:"period"`
	GameClock         string                 `json:"gameClock"`
	GameTimeUTC       string                 `json:"gameTimeUTC"`
	GameEt            string                 `json:"gameEt"`
	RegulationPeriods int                    `json:"regulationPeriods"`
	SeriesGameNumber  string                 `json:"seriesGameNumber,omitempty"`
	SeriesText        string                 `json:"seriesText,omitempty"`
	SeriesConference  string                 `json:"seriesConference,omitempty"`
	PoRoundDesc       string                 `json:"poRoundDesc,omitempty"`
	GameSubtype       string                 `json:"gameSubtype,omitempty"`
	IfNecessary       bool                   `json:"ifNecessary,omitempty"`
	HomeTeam          ScoreboardV3Team         `json:"homeTeam"`
	AwayTeam          ScoreboardV3Team         `json:"awayTeam"`
	GameLeaders       ScoreboardV3Leaders      `json:"gameLeaders,omitempty"`
	PbOdds            ScoreboardV3PbOdds       `json:"pbOdds,omitempty"`
	Broadcasters      ScoreboardV3Broadcasters `json:"broadcasters,omitempty"`
}

// ScoreboardV3Team is one team's scoreboard entry for a game.
type ScoreboardV3Team struct {
	TeamID            int                `json:"teamId"`
	TeamName          string             `json:"teamName"`
	TeamCity          string             `json:"teamCity"`
	TeamTricode       string             `json:"teamTricode"`
	TeamSlug          string             `json:"teamSlug,omitempty"`
	Wins              int                `json:"wins"`
	Losses            int                `json:"losses"`
	Score             int                `json:"score"`
	Seed              int                `json:"seed,omitempty"`
	InBonus           string             `json:"inBonus,omitempty"`
	TimeoutsRemaining int                `json:"timeoutsRemaining"`
	Periods           []ScoreboardV3Period `json:"periods"`
}

// ScoreboardV3Period is one team's score for a single period.
type ScoreboardV3Period struct {
	Period     int    `json:"period"`
	PeriodType string `json:"periodType"`
	Score      int    `json:"score"`
}

// ScoreboardV3Leaders holds the top scoring leaders for each team in a game.
type ScoreboardV3Leaders struct {
	HomeLeaders ScoreboardV3Leader`json:"homeLeaders"`
	AwayLeaders ScoreboardV3Leader`json:"awayLeaders"`
}

// ScoreboardV3Leader is a single statistical leader entry.
type ScoreboardV3Leader struct {
	PersonID    int    `json:"personId"`
	Name        string `json:"name"`
	JerseyNum   string `json:"jerseyNum"`
	Position    string `json:"position"`
	TeamTricode string `json:"teamTricode"`
	PlayerSlug  string `json:"playerSlug"`
	Points      int    `json:"points"`
	Rebounds    int    `json:"rebounds"`
	Assists     int    `json:"assists"`
}

// ScoreboardV3PbOdds carries pre-game odds when available.
type ScoreboardV3PbOdds struct {
	Team        string  `json:"team,omitempty"`
	Odds        float64 `json:"odds,omitempty"`
	Suspended   int     `json:"suspended,omitempty"`
}

// ScoreboardV3Broadcasters lists national, home, and away broadcasters.
type ScoreboardV3Broadcasters struct {
	NationalTvBroadcasters    []ScoreboardV3Broadcaster`json:"nationalTvBroadcasters,omitempty"`
	NationalRadioBroadcasters []ScoreboardV3Broadcaster`json:"nationalRadioBroadcasters,omitempty"`
	HomeTvBroadcasters        []ScoreboardV3Broadcaster`json:"homeTvBroadcasters,omitempty"`
	HomeRadioBroadcasters     []ScoreboardV3Broadcaster`json:"homeRadioBroadcasters,omitempty"`
	AwayTvBroadcasters        []ScoreboardV3Broadcaster`json:"awayTvBroadcasters,omitempty"`
	AwayRadioBroadcasters     []ScoreboardV3Broadcaster`json:"awayRadioBroadcasters,omitempty"`
}

// ScoreboardV3Broadcaster is a single broadcaster entry.
type ScoreboardV3Broadcaster struct {
	BroadcasterScope    string `json:"broadcasterScope,omitempty"`
	BroadcasterMedia    string `json:"broadcasterMedia,omitempty"`
	BroadcasterID       int    `json:"broadcasterId,omitempty"`
	BroadcasterDisplay  string `json:"broadcasterDisplay,omitempty"`
	BroadcasterAbbrev   string `json:"broadcasterAbbreviation,omitempty"`
	BroadcasterDescr    string `json:"broadcasterDescription,omitempty"`
	TapeDelayComments   string `json:"tapeDelayComments,omitempty"`
	RegionID            int    `json:"regionId,omitempty"`
}

// ScoreboardV3 fetches every game on the given calendar date, regardless
// of game status (scheduled, in-progress, or final). Use this when you
// need an explicit date rather than the CDN's idea of "today."
//
// date is interpreted in NBA local terms — pass a time.Time and we'll
// format it as YYYY-MM-DD.
func (c *Client) ScoreboardV3(ctx context.Context, date time.Time) (*ScoreboardV3Response, error) {
	q := url.Values{}
	q.Set("GameDate", date.Format("2006-01-02"))
	q.Set("LeagueID", "00")

	var out ScoreboardV3Response
	if err := c.http.GetJSON(ctx, httpx.ProfileStats, c.baseURL+"/stats/scoreboardv3", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
