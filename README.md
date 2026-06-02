# nba-sdk

A Go SDK for the NBA's unofficial JSON endpoints.

## Install

```sh
go get github.com/NolanFogarty/nba-sdk
```

## Quickstart

```go
package main

import (
    "context"
    "fmt"
    "log"

    nba "github.com/NolanFogarty/nba-sdk"
)

func main() {
    client := nba.NewClient()
    ctx := context.Background()

    sb, err := client.Live.Scoreboard(ctx)
    if err != nil { log.Fatal(err) }
    for _, g := range sb.Scoreboard.Games {
        fmt.Printf("%s @ %s — %s\n", g.AwayTeam.TeamTricode, g.HomeTeam.TeamTricode, g.GameStatusText)
    }

    box, _ := client.Stats.BoxScoreTraditionalV3(ctx, "0022400001")
    fmt.Println("Home points:", box.BoxScoreTraditional.HomeTeam.Statistics.Points)

    pbp, _ := client.Stats.PlayByPlayV3(ctx, "0022400001")
    fmt.Println("Total plays:", len(pbp.Game.Actions))
}
```

## API surface

Two sub-clients:

- `client.Stats.*` — stats.nba.com endpoints
  - `ScoreboardV2(ctx, date)` — every game on a given date (`"YYYY-MM-DD"`), including past dates; use it to resolve game IDs for historical games
  - `BoxScoreTraditionalV3(ctx, gameID)` — full-game traditional box score
  - `PlayByPlayV3(ctx, gameID)` — every play-by-play action for the game
- `client.Live.*` — cdn.nba.com endpoints
  - `Scoreboard(ctx)` — every game scheduled today (upcoming, in-progress, finished)

Every method returns a fully populated typed struct mirroring the NBA JSON
response. Callers filter or project the data however they need — for
example, to get Q1-only plays, iterate `pbp.Game.Actions` and keep those
with `Period == 1`.

## Options

```go
client := nba.NewClient(
    nba.WithStatsRateLimit(1.0, 3),       // 1 req/sec, burst 3
    nba.WithLiveRateLimit(5.0, 10),       // CDN is friendlier
    nba.WithRetry(3, 500*time.Millisecond),
    nba.WithUserAgent("my-app/1.0"),
)
```

The SDK ships with a built-in token-bucket rate limiter and exponential
backoff on 429/5xx responses. stats.nba.com is aggressive about banning
IPs, so the defaults are conservative.

## Notes

- These endpoints are **unofficial**. There is no published SLA or schema,
  and the NBA may change or block them at any time. Callers are responsible
  for compliance with NBA's terms of use.
- stats.nba.com requires specific HTTP headers (`Referer`, `x-nba-stats-origin`,
  `x-nba-stats-token`, etc.) — the SDK handles this for you.
- The bundled golden test fixtures (`stats/testdata/`, `live/testdata/`) are
  hand-written to match the documented v3 response shapes. If you have a
  non-blocked network path to the NBA API, you can replace them with real
  captured responses.

## Testing

```sh
go test ./...
```

Tests use an `httptest.Server` that serves the bundled fixtures — no
network calls. CI is safe.
