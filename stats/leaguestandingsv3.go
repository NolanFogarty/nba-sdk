package stats

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/NolanFogarty/nba-sdk/internal/httpx"
)

// LeagueStandingsV3Response mirrors stats.nba.com/stats/leaguestandingsv3.
//
// Despite the "v3" name, the endpoint returns the classic resultSets shape
// (positional rows). The SDK projects the "Standings" result set into a
// typed []TeamStanding so callers get clean field access without decoding
// columns themselves.
type LeagueStandingsV3Response struct {
	Resource   string                      `json:"resource"`
	Parameters LeagueStandingsV3Parameters `json:"parameters"`

	// Standings is the typed projection of the "Standings" result set,
	// one row per team. Sort order matches what the NBA returns (typically
	// by conference, then playoff rank).
	Standings []TeamStanding `json:"-"`
}

// LeagueStandingsV3Parameters echoes the query parameters the endpoint resolved.
type LeagueStandingsV3Parameters struct {
	LeagueID   string `json:"LeagueID"`
	Season     string `json:"Season"`
	SeasonType string `json:"SeasonType"`
	SeasonYear string `json:"SeasonYear,omitempty"`
}

// TeamStanding is one team's row in the league standings. Records that the
// NBA returns as strings ("30-22", "5-3", etc.) are preserved as strings —
// they capture nuance like overtime games that simple W/L ints can't.
type TeamStanding struct {
	// Identity
	LeagueID  string `json:"leagueId"`
	SeasonID  string `json:"seasonId"`
	TeamID    int    `json:"teamId"`
	TeamCity  string `json:"teamCity"`
	TeamName  string `json:"teamName"`
	TeamSlug  string `json:"teamSlug"`

	// Conference / division
	Conference       string `json:"conference"`
	ConferenceRecord string `json:"conferenceRecord"`
	PlayoffRank      int    `json:"playoffRank"`
	ClinchIndicator  string `json:"clinchIndicator"`
	Division         string `json:"division"`
	DivisionRecord   string `json:"divisionRecord"`
	DivisionRank     int    `json:"divisionRank"`

	// Win/loss
	Wins       int     `json:"wins"`
	Losses     int     `json:"losses"`
	WinPCT     float64 `json:"winPct"`
	LeagueRank int     `json:"leagueRank"`
	Record     string  `json:"record"`

	// Splits
	Home       string `json:"home"`
	Road       string `json:"road"`
	L10        string `json:"l10"`
	Last10Home string `json:"last10Home"`
	Last10Road string `json:"last10Road"`
	OT         string `json:"ot"`

	// Scoring
	PointsPG     float64 `json:"pointsPg"`
	OppPointsPG  float64 `json:"oppPointsPg"`
	DiffPointsPG float64 `json:"diffPointsPg"`

	// Versus conferences / divisions
	VsEast      string `json:"vsEast"`
	VsAtlantic  string `json:"vsAtlantic"`
	VsCentral   string `json:"vsCentral"`
	VsSoutheast string `json:"vsSoutheast"`
	VsWest      string `json:"vsWest"`
	VsNorthwest string `json:"vsNorthwest"`
	VsPacific   string `json:"vsPacific"`
	VsSouthwest string `json:"vsSouthwest"`

	// Streaks
	LongHomeStreak       int    `json:"longHomeStreak"`
	StrLongHomeStreak    string `json:"strLongHomeStreak"`
	LongRoadStreak       int    `json:"longRoadStreak"`
	StrLongRoadStreak    string `json:"strLongRoadStreak"`
	LongWinStreak        int    `json:"longWinStreak"`
	LongLossStreak       int    `json:"longLossStreak"`
	CurrentHomeStreak    int    `json:"currentHomeStreak"`
	StrCurrentHomeStreak string `json:"strCurrentHomeStreak"`
	CurrentRoadStreak    int    `json:"currentRoadStreak"`
	StrCurrentRoadStreak string `json:"strCurrentRoadStreak"`
	CurrentStreak        int    `json:"currentStreak"`
	StrCurrentStreak     string `json:"strCurrentStreak"`

	// Games behind
	ConferenceGamesBack float64 `json:"conferenceGamesBack"`
	DivisionGamesBack   float64 `json:"divisionGamesBack"`

	// Clinching / elimination
	ClinchedConferenceTitle int `json:"clinchedConferenceTitle"`
	ClinchedDivisionTitle   int `json:"clinchedDivisionTitle"`
	ClinchedPlayoffBirth    int `json:"clinchedPlayoffBirth"`
	ClinchedPlayIn          int `json:"clinchedPlayIn"`
	EliminatedConference    int `json:"eliminatedConference"`
	EliminatedDivision      int `json:"eliminatedDivision"`

	// Situational
	AheadAtHalf      string `json:"aheadAtHalf"`
	BehindAtHalf     string `json:"behindAtHalf"`
	TiedAtHalf       string `json:"tiedAtHalf"`
	AheadAfterThree  string `json:"aheadAfterThree"`
	BehindAfterThree string `json:"behindAfterThree"`
	TiedAfterThree   string `json:"tiedAfterThree"`
	Score100PTS      string `json:"score100Pts"`
	OppScore100PTS   string `json:"oppScore100Pts"`
	OppOver500       string `json:"oppOver500"`
	LeadInFGPCT      string `json:"leadInFgPct"`
	LeadInReb        string `json:"leadInReb"`
	FewerTurnovers   string `json:"fewerTurnovers"`

	// Rankings
	PointsPGRank     int `json:"pointsPgRank"`
	OppPointsPGRank  int `json:"oppPointsPgRank"`
	DiffPointsPGRank int `json:"diffPointsPgRank"`
}

// UnmarshalJSON decodes the classic stats.nba.com envelope and then projects
// the "Standings" result set into typed TeamStanding rows. Columns are
// matched by header name, so the SDK is resilient to the NBA reordering or
// adding columns.
func (r *LeagueStandingsV3Response) UnmarshalJSON(data []byte) error {
	type envelope struct {
		Resource   string                      `json:"resource"`
		Parameters LeagueStandingsV3Parameters `json:"parameters"`
		ResultSets []resultSet                 `json:"resultSets"`
	}
	var raw envelope
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	r.Resource = raw.Resource
	r.Parameters = raw.Parameters

	for _, rs := range raw.ResultSets {
		if rs.Name != "Standings" {
			continue
		}
		idx := indexHeaders(rs.Headers)
		for _, row := range rs.RowSet {
			var t TeamStanding
			getString(row, idx, "LeagueID", &t.LeagueID)
			getString(row, idx, "SeasonID", &t.SeasonID)
			getInt(row, idx, "TeamID", &t.TeamID)
			getString(row, idx, "TeamCity", &t.TeamCity)
			getString(row, idx, "TeamName", &t.TeamName)
			getString(row, idx, "TeamSlug", &t.TeamSlug)

			getString(row, idx, "Conference", &t.Conference)
			getString(row, idx, "ConferenceRecord", &t.ConferenceRecord)
			getInt(row, idx, "PlayoffRank", &t.PlayoffRank)
			getString(row, idx, "ClinchIndicator", &t.ClinchIndicator)
			getString(row, idx, "Division", &t.Division)
			getString(row, idx, "DivisionRecord", &t.DivisionRecord)
			getInt(row, idx, "DivisionRank", &t.DivisionRank)

			getInt(row, idx, "WINS", &t.Wins)
			getInt(row, idx, "LOSSES", &t.Losses)
			getFloat(row, idx, "WinPCT", &t.WinPCT)
			getInt(row, idx, "LeagueRank", &t.LeagueRank)
			getString(row, idx, "Record", &t.Record)

			getString(row, idx, "HOME", &t.Home)
			getString(row, idx, "ROAD", &t.Road)
			getString(row, idx, "L10", &t.L10)
			getString(row, idx, "Last10Home", &t.Last10Home)
			getString(row, idx, "Last10Road", &t.Last10Road)
			getString(row, idx, "OT", &t.OT)

			getFloat(row, idx, "PointsPG", &t.PointsPG)
			getFloat(row, idx, "OppPointsPG", &t.OppPointsPG)
			getFloat(row, idx, "DiffPointsPG", &t.DiffPointsPG)

			getString(row, idx, "vsEast", &t.VsEast)
			getString(row, idx, "vsAtlantic", &t.VsAtlantic)
			getString(row, idx, "vsCentral", &t.VsCentral)
			getString(row, idx, "vsSoutheast", &t.VsSoutheast)
			getString(row, idx, "vsWest", &t.VsWest)
			getString(row, idx, "vsNorthwest", &t.VsNorthwest)
			getString(row, idx, "vsPacific", &t.VsPacific)
			getString(row, idx, "vsSouthwest", &t.VsSouthwest)

			getInt(row, idx, "LongHomeStreak", &t.LongHomeStreak)
			getString(row, idx, "strLongHomeStreak", &t.StrLongHomeStreak)
			getInt(row, idx, "LongRoadStreak", &t.LongRoadStreak)
			getString(row, idx, "strLongRoadStreak", &t.StrLongRoadStreak)
			getInt(row, idx, "LongWinStreak", &t.LongWinStreak)
			getInt(row, idx, "LongLossStreak", &t.LongLossStreak)
			getInt(row, idx, "CurrentHomeStreak", &t.CurrentHomeStreak)
			getString(row, idx, "strCurrentHomeStreak", &t.StrCurrentHomeStreak)
			getInt(row, idx, "CurrentRoadStreak", &t.CurrentRoadStreak)
			getString(row, idx, "strCurrentRoadStreak", &t.StrCurrentRoadStreak)
			getInt(row, idx, "CurrentStreak", &t.CurrentStreak)
			getString(row, idx, "strCurrentStreak", &t.StrCurrentStreak)

			getFloat(row, idx, "ConferenceGamesBack", &t.ConferenceGamesBack)
			getFloat(row, idx, "DivisionGamesBack", &t.DivisionGamesBack)

			getInt(row, idx, "ClinchedConferenceTitle", &t.ClinchedConferenceTitle)
			getInt(row, idx, "ClinchedDivisionTitle", &t.ClinchedDivisionTitle)
			getInt(row, idx, "ClinchedPlayoffBirth", &t.ClinchedPlayoffBirth)
			getInt(row, idx, "ClinchedPlayIn", &t.ClinchedPlayIn)
			getInt(row, idx, "EliminatedConference", &t.EliminatedConference)
			getInt(row, idx, "EliminatedDivision", &t.EliminatedDivision)

			getString(row, idx, "AheadAtHalf", &t.AheadAtHalf)
			getString(row, idx, "BehindAtHalf", &t.BehindAtHalf)
			getString(row, idx, "TiedAtHalf", &t.TiedAtHalf)
			getString(row, idx, "AheadAfterThree", &t.AheadAfterThree)
			getString(row, idx, "BehindAfterThree", &t.BehindAfterThree)
			getString(row, idx, "TiedAfterThree", &t.TiedAfterThree)
			getString(row, idx, "Score100PTS", &t.Score100PTS)
			getString(row, idx, "OppScore100PTS", &t.OppScore100PTS)
			getString(row, idx, "OppOver500", &t.OppOver500)
			getString(row, idx, "LeadInFGPCT", &t.LeadInFGPCT)
			getString(row, idx, "LeadInReb", &t.LeadInReb)
			getString(row, idx, "FewerTurnovers", &t.FewerTurnovers)

			getInt(row, idx, "PointsPG_Rank", &t.PointsPGRank)
			getInt(row, idx, "OppPointsPG_Rank", &t.OppPointsPGRank)
			getInt(row, idx, "DiffPointsPG_Rank", &t.DiffPointsPGRank)

			r.Standings = append(r.Standings, t)
		}
		break
	}
	return nil
}

// LeagueStandingsV3 returns the regular-season standings for the given
// season. season must be in NBA "YYYY-YY" form (e.g. "2025-26").
//
// The response contains one TeamStanding row per team, with the full record,
// splits, streaks, clinching indicators, and scoring averages. Sort order
// matches what the NBA returns (by conference, then playoff rank).
func (c *Client) LeagueStandingsV3(ctx context.Context, season string) (*LeagueStandingsV3Response, error) {
	if season == "" {
		return nil, fmt.Errorf("stats: LeagueStandingsV3 requires a season in YYYY-YY form (e.g. \"2025-26\")")
	}
	q := url.Values{}
	q.Set("LeagueID", "00")
	q.Set("Season", season)
	q.Set("SeasonType", "Regular Season")

	var out LeagueStandingsV3Response
	if err := c.http.GetJSON(ctx, httpx.ProfileStats, c.baseURL+"/stats/leaguestandingsv3", q, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
