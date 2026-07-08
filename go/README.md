# greyhoundapi-go

Official Go client for the **[GreyhoundAPI](https://greyhoundapi.com)** — greyhound racing data for Great Britain and Australia: racecards, results, sectionals, starting prices and Betfair markets, plus dogs, trainers, owners, tracks and breeding.

- **Go 1.18+**, standard library only (zero dependencies)
- Covers all 35 REST endpoints, with cursor pagination and a typed error

## Install

```sh
go get github.com/greyhoundapi/greyhoundapi-go
```

```go
import greyhoundapi "github.com/greyhoundapi/greyhoundapi-go"
```

## Getting an API key

The API is authenticated with a key sent in the `X-API-Key` header. To get one:

1. Go to **[greyhoundapi.com](https://greyhoundapi.com)** and click **Get a free key** (top-right).
2. **Sign in** to create your account.
3. Open **Account → API keys** and click **Create key**.
4. Copy the key — **it's shown only once**. Sandbox keys start `gapi_test_`; live keys start `gapi_live_`.

**Free sandbox key** — no card required: the race & track endpoints over a rolling 7-day window, 500 requests/day, one active key. Ideal for building and evaluating.

**Live plan** ($99/month) — every endpoint, the full historical archive, live-day data and the WebSocket results stream, 250,000 requests/month, and up to 5 active keys. See **[pricing](https://greyhoundapi.com/pricing)**.

**Keep your key secret.** Load it from an environment variable (e.g. `GREYHOUNDAPI_KEY`) rather than committing it.

## Quick start

```go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	greyhoundapi "github.com/greyhoundapi/greyhoundapi-go"
)

func main() {
	gapi := greyhoundapi.New(os.Getenv("GREYHOUNDAPI_KEY"))

	// today's GB racecards
	if _, err := gapi.RacecardsToday(greyhoundapi.Params{"region": "GB"}); err != nil {
		log.Fatal(err)
	}

	// one race, fully resolved
	race, err := gapi.Race(1229082)
	if err != nil {
		log.Fatal(err)
	}
	var data struct {
		Runners []struct {
			DogName string `json:"dog_name"`
		} `json:"runners"`
	}
	if err := json.Unmarshal(race.Data, &data); err != nil {
		log.Fatal(err)
	}
	fmt.Println(data.Runners[0].DogName)
}
```

Every call returns an `*Envelope`. `Meta` is decoded for you; `Data` is raw JSON
you unmarshal into your own type:

```go
type Envelope struct {
	Meta Meta            // request_id, data_as_of, count, next_cursor
	Data json.RawMessage // unmarshal into your own struct
}
```

## Authentication

Pass your key to `New`; it's sent on every request as `X-API-Key`. Override
`BaseURL` or `HTTP` (e.g. a custom timeout) on the returned client if needed.

```go
gapi := greyhoundapi.New("gapi_live_...")
// gapi.BaseURL = "https://api.greyhoundapi.com/v1"
// gapi.HTTP = &http.Client{Timeout: 10 * time.Second}
```

## Methods

Query parameters are passed as `Params` (a `map[string]interface{}`); path
parameters are positional. Each method returns `(*Envelope, error)`.

| Method | Endpoint |
| --- | --- |
| `RacecardsToday(p)` | `GET /racecards/today` |
| `RacecardsUpcoming(p)` | `GET /racecards/upcoming` |
| `Races(p)` | `GET /races` |
| `Race(raceID)` | `GET /races/{race_id}` |
| `RaceRunners(raceID, p)` | `GET /races/{race_id}/runners` |
| `RaceResult(raceID)` | `GET /races/{race_id}/result` |
| `RaceStatus(raceID)` | `GET /races/{race_id}/status` |
| `RaceMarket(raceID)` | `GET /races/{race_id}/market` |
| `ResultsToday(p)` | `GET /results/today` |
| `Results(p)` | `GET /results` |
| `LatestResults(p)` | `GET /results/latest` |
| `Meetings(p)` | `GET /meetings` |
| `MeetingsToday(p)` | `GET /meetings/today` |
| `Meeting(meetingID)` | `GET /meetings/{meeting_id}` |
| `SearchDogs(p)` | `GET /dogs/search` |
| `Dog(dogID)` | `GET /dogs/{dog_id}` |
| `DogForm(dogID, p)` | `GET /dogs/{dog_id}/form` |
| `DogEntries(dogID)` | `GET /dogs/{dog_id}/entries` |
| `DogPrices(dogID)` | `GET /dogs/{dog_id}/prices` |
| `DogHeadToHead(dogID, rivalID)` | `GET /dogs/{dog_id}/head-to-head/{rival_id}` |
| `SearchTrainers(p)` | `GET /trainers/search` |
| `Trainer(trainerID)` | `GET /trainers/{trainer_id}` |
| `TrainerRunners(trainerID)` | `GET /trainers/{trainer_id}/runners` |
| `TrainerResults(trainerID, p)` | `GET /trainers/{trainer_id}/results` |
| `SearchOwners(p)` | `GET /owners/search` |
| `Owner(ownerID)` | `GET /owners/{owner_id}` |
| `Tracks(p)` | `GET /tracks` |
| `Track(trackID)` | `GET /tracks/{track_id}` |
| `TrackRaces(trackID, p)` | `GET /tracks/{track_id}/races` |
| `TrackStats(trackID)` | `GET /tracks/{track_id}/stats` |
| `SireProgeny(name, p)` | `GET /sires/{name}/progeny` |
| `DamProgeny(name, p)` | `GET /dams/{name}/progeny` |
| `Status()` | `GET /status` |
| `Usage()` | `GET /usage` |
| `Reference(p)` | `GET /reference` |

### Common parameters

Most list endpoints accept: `region` (`GB` \| `AU`), `date_from` / `date_to`
(`YYYY-MM-DD`, track-local), `track_id`, `grade`, `distance_m`, `limit` (1–200)
and `cursor`. See the [docs](https://greyhoundapi.com/documentation).

```go
races, err := gapi.Races(greyhoundapi.Params{"region": "GB", "date_from": "2026-07-01", "limit": 100})
```

## Pagination

List endpoints page with an opaque cursor in `Meta.NextCursor`. `Paginate`
follows it across every page, calling your function for each item:

```go
err := gapi.Paginate("/races", greyhoundapi.Params{"region": "GB"}, func(item json.RawMessage) error {
	var r struct {
		RaceID int `json:"race_id"`
	}
	json.Unmarshal(item, &r)
	fmt.Println(r.RaceID)
	return nil // return a non-nil error to stop early
})
```

## Errors

Any non-2xx response (or transport failure) returns a `*greyhoundapi.Error`:

```go
_, err := gapi.Race(999999999)
var apiErr *greyhoundapi.Error
if errors.As(err, &apiErr) {
	fmt.Println(apiErr.Status, apiErr.Code, apiErr.Message, apiErr.RequestID)
}
```

## Links

- **Documentation** — https://greyhoundapi.com/documentation
- **Service status** — https://greyhoundapi.com/status
- **Pricing** — https://greyhoundapi.com/pricing

## License

[MIT](./LICENSE)
