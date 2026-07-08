# GreyhoundAPI SDKs

Official client libraries for the **[GreyhoundAPI](https://greyhoundapi.com)** — a REST + WebSocket API for greyhound racing in Great Britain and Australia: racecards, results, sectionals, starting prices and Betfair markets, plus dogs, trainers, owners, tracks and breeding.

| Language | Package | Directory |
| --- | --- | --- |
| JavaScript / TypeScript | `@greyhoundapi/sdk` (npm) | [`javascript/`](./javascript) |
| Python | `greyhoundapi` (PyPI) | [`python/`](./python) |
| PHP | `greyhoundapi/sdk` (Composer) | [`php/`](./php) |
| Go | `github.com/greyhoundapi/greyhoundapi-go` | [`go/`](./go) |

All four cover the same 35 REST endpoints, share the same `{ meta, data }` envelope and cursor pagination, and have **zero third-party dependencies**.

## Getting an API key

The API is authenticated with a key sent in the `X-API-Key` header. To get one:

1. Go to **[greyhoundapi.com](https://greyhoundapi.com)** and click **Get a free key** (top-right).
2. **Sign in** to create your account.
3. Open **Account → API keys** and click **Create key**.
4. Copy the key — **it's shown only once**. Sandbox keys start `gapi_test_`; live keys start `gapi_live_`.

**Free sandbox key** — no card required: the race & track endpoints over a rolling 7-day window (ending a few hours ago), 500 requests/day, one active key. Ideal for building and evaluating.

**Live plan** ($99/month) — every endpoint, the full historical archive, live-day data and the WebSocket results stream, 250,000 requests/month, and up to 5 active keys. See **[pricing](https://greyhoundapi.com/pricing)** for the full comparison.

**Keep your key secret.** Load it from an environment variable (e.g. `GREYHOUNDAPI_KEY`) rather than committing it. In browser code the key is visible to users — use a sandbox key for public demos, or proxy requests through your own server.

## Quick starts

### JavaScript / TypeScript
```sh
npm install @greyhoundapi/sdk
```
```js
import { GreyhoundAPI } from "@greyhoundapi/sdk";

const gapi = new GreyhoundAPI({ apiKey: process.env.GREYHOUNDAPI_KEY });
const race = await gapi.races.get(1229082);
console.log(race.data.runners[0].dog_name);
```

### Python
```sh
pip install greyhoundapi
```
```python
from greyhoundapi import GreyhoundAPI

gapi = GreyhoundAPI(api_key="gapi_live_...")
race = gapi.race(1229082)
print(race["data"]["runners"][0]["dog_name"])
```

### PHP
```sh
composer require greyhoundapi/sdk
```
```php
$gapi = new \GreyhoundApi\Client(getenv('GREYHOUNDAPI_KEY'));
$race = $gapi->race(1229082);
echo $race['data']['runners'][0]['dog_name'];
```

### Go
```sh
go get github.com/greyhoundapi/greyhoundapi-go
```
```go
gapi := greyhoundapi.New(os.Getenv("GREYHOUNDAPI_KEY"))
race, err := gapi.Race(1229082)
```

Full guides → [JavaScript](./javascript/README.md) · [Python](./python/README.md) · [PHP](./php/README.md) · [Go](./go/README.md)

## About the API

- **Base URL** — `https://api.greyhoundapi.com/v1`
- **Auth** — `X-API-Key` header
- **Envelope** — `{ "meta": { "request_id", "data_as_of" }, "data": … }`
- **Pagination** — opaque cursor in `meta.next_cursor`
- **Tiers** — a free sandbox key (race & track endpoints, rolling 7-day window, 500/day) or the live plan (every endpoint, full archive, results stream, 250k/month)

Read the docs at **[greyhoundapi.com/documentation](https://greyhoundapi.com/documentation)**.

## Contributing

Issues and pull requests are welcome. Please keep changes focused and preserve
the zero-dependency, build-step-free design of each client.

## License

[MIT](./LICENSE)
