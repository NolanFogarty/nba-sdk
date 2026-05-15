// Package nba is a Go SDK for the NBA's unofficial JSON endpoints.
//
// The SDK exposes two sub-clients:
//
//   - Client.Stats wraps stats.nba.com (boxscoretraditionalv3, playbyplayv3, ...)
//   - Client.Live wraps cdn.nba.com (today's scoreboard, live game data)
//
// Endpoint methods always take a context.Context first and return a fully
// populated response struct mirroring the NBA JSON response. Callers filter
// the data however they like.
//
// stats.nba.com is unofficial and aggressively rate-limited. Callers are
// responsible for compliance with NBA's terms of use.
package nba
