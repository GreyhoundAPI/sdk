# greyhoundapi

Official Python client for the **[GreyhoundAPI](https://greyhoundapi.com)** — greyhound racing data for Great Britain and Australia: racecards, results, sectionals, starting prices and Betfair markets, plus dogs, trainers, owners, tracks and breeding.

- **Python 3.8+**
- **Zero dependencies** — standard library only
- Covers all 35 REST endpoints, with cursor pagination and a typed error

## Install

```sh
pip install greyhoundapi
```

## Quick start

```python
from greyhoundapi import GreyhoundAPI

gapi = GreyhoundAPI(api_key="gapi_live_...")

# today's GB racecards
card = gapi.racecards_today(region="GB")

# one race, fully resolved
race = gapi.race(1229082)
print(race["data"]["runners"][0]["dog_name"])
```

Every call returns the standard envelope as a dict:

```python
{"meta": {"request_id": "req_…", "data_as_of": "2026-07-08T13:05:12Z"}, "data": {...}}
```

Get a free key at **[greyhoundapi.com](https://greyhoundapi.com)** — no card required.

## Authentication

Pass your key to the constructor; it's sent on every request as the `X-API-Key` header.

```python
gapi = GreyhoundAPI(
    api_key="gapi_live_...",                        # required
    # base_url="https://api.greyhoundapi.com/v1",   # optional override
    # timeout=30.0,                                 # optional, seconds
)
```

## Methods

Query parameters are passed as keyword arguments; path parameters are positional.
Each method returns the response envelope as a dict.

| Method | Endpoint |
| --- | --- |
| `racecards_today(**params)` | `GET /racecards/today` |
| `racecards_upcoming(**params)` | `GET /racecards/upcoming` |
| `races(**params)` | `GET /races` |
| `race(race_id)` | `GET /races/{race_id}` |
| `race_runners(race_id, **params)` | `GET /races/{race_id}/runners` |
| `race_result(race_id)` | `GET /races/{race_id}/result` |
| `race_status(race_id)` | `GET /races/{race_id}/status` |
| `race_market(race_id)` | `GET /races/{race_id}/market` |
| `results_today(**params)` | `GET /results/today` |
| `results(**params)` | `GET /results` |
| `latest_results(**params)` | `GET /results/latest` |
| `meetings(**params)` | `GET /meetings` |
| `meetings_today(**params)` | `GET /meetings/today` |
| `meeting(meeting_id)` | `GET /meetings/{meeting_id}` |
| `search_dogs(**params)` | `GET /dogs/search` |
| `dog(dog_id)` | `GET /dogs/{dog_id}` |
| `dog_form(dog_id, **params)` | `GET /dogs/{dog_id}/form` |
| `dog_entries(dog_id)` | `GET /dogs/{dog_id}/entries` |
| `dog_prices(dog_id)` | `GET /dogs/{dog_id}/prices` |
| `dog_head_to_head(dog_id, rival_id)` | `GET /dogs/{dog_id}/head-to-head/{rival_id}` |
| `search_trainers(**params)` | `GET /trainers/search` |
| `trainer(trainer_id)` | `GET /trainers/{trainer_id}` |
| `trainer_runners(trainer_id)` | `GET /trainers/{trainer_id}/runners` |
| `trainer_results(trainer_id, **params)` | `GET /trainers/{trainer_id}/results` |
| `search_owners(**params)` | `GET /owners/search` |
| `owner(owner_id)` | `GET /owners/{owner_id}` |
| `tracks(**params)` | `GET /tracks` |
| `track(track_id)` | `GET /tracks/{track_id}` |
| `track_races(track_id, **params)` | `GET /tracks/{track_id}/races` |
| `track_stats(track_id)` | `GET /tracks/{track_id}/stats` |
| `sire_progeny(name, **params)` | `GET /sires/{name}/progeny` |
| `dam_progeny(name, **params)` | `GET /dams/{name}/progeny` |
| `status()` | `GET /status` |
| `usage()` | `GET /usage` |
| `reference(**params)` | `GET /reference` |

### Common parameters

Most list endpoints accept: `region` (`"GB"` or `"AU"`), `date_from` / `date_to`
(`"YYYY-MM-DD"`, track-local), `track_id`, `grade`, `distance_m`, `limit` (1–200),
and `cursor`. See the [docs](https://greyhoundapi.com/documentation) for what
each endpoint supports.

```python
races = gapi.races(region="GB", date_from="2026-07-01", grade="A2", limit=100)
```

## Pagination

List endpoints page with an opaque cursor in `meta["next_cursor"]`. Pass `cursor`
yourself, or let the SDK walk every page and yield each item:

```python
for race in gapi.paginate("/races", region="GB", date_from="2026-07-01"):
    print(race["race_id"])
```

## Errors

Any non-2xx response (or a transport failure) raises `GreyhoundAPIError`:

```python
from greyhoundapi import GreyhoundAPI, GreyhoundAPIError

try:
    gapi.race(999999999)
except GreyhoundAPIError as err:
    print(err.status, err.code, err, err.request_id)
```

## Links

- **Documentation** — https://greyhoundapi.com/documentation
- **Service status** — https://greyhoundapi.com/status
- **Pricing** — https://greyhoundapi.com/pricing

## License

[MIT](./LICENSE)
