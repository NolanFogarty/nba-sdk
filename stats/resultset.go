package stats

import "encoding/json"

// resultSet is one named table in the classic stats.nba.com response shape
// ({"resultSets": [{"name": "...", "headers": [...], "rowSet": [[...]]}]}).
// It is internal — callers of the SDK interact with the typed projections
// each endpoint exposes, not the raw positional rows.
type resultSet struct {
	Name    string              `json:"name"`
	Headers []string            `json:"headers"`
	RowSet  [][]json.RawMessage `json:"rowSet"`
}

// indexHeaders returns a map from header name to column index, used to look
// up columns by name in a positional row.
func indexHeaders(headers []string) map[string]int {
	idx := make(map[string]int, len(headers))
	for i, h := range headers {
		idx[h] = i
	}
	return idx
}

// cell returns the raw value for column name in row, or nil if the column
// is absent or null.
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

func getFloat(row []json.RawMessage, idx map[string]int, name string, dst *float64) {
	if v := cell(row, idx, name); v != nil {
		_ = json.Unmarshal(v, dst)
	}
}
