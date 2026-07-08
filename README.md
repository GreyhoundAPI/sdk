# GreyhoundAPI SDKs

Official client libraries for the **[GreyhoundAPI](https://greyhoundapi.com)** — a REST + WebSocket API for greyhound racing in Great Britain and Australia: racecards, results, sectionals, starting prices and Betfair markets, plus dogs, trainers, owners, tracks and breeding.

| Language | Package | Directory |
| --- | --- | --- |
| JavaScript / TypeScript | `@greyhoundapi/sdk` | [`javascript/`](./javascript) |
| Python | `greyhoundapi` | [`python/`](./python) |

Both cover all 35 REST endpoints, share the same response envelope and cursor pagination, and have **zero runtime dependencies**.

## JavaScript / TypeScript

```sh
npm install @greyhoundapi/sdk
```

```js
import { GreyhoundAPI } from "@greyhoundapi/sdk";

const gapi = new GreyhoundAPI({ apiKey: process.env.GREYHOUNDAPI_KEY });
const race = await gapi.races.get(1229082);
console.log(race.data.runners[0].dog_name);
```

Full guide → [javascript/README.md](./javascript/README.md)

## Python

```sh
pip install greyhoundapi
```

```python
from greyhoundapi import GreyhoundAPI

gapi = GreyhoundAPI(api_key="gapi_live_...")
race = gapi.race(1229082)
print(race["data"]["runners"][0]["dog_name"])
```

Full guide → [python/README.md](./python/README.md)

## About the API

- **Base URL** — `https://api.greyhoundapi.com/v1`
- **Auth** — `X-API-Key` header
- **Envelope** — `{ "meta": { "request_id", "data_as_of" }, "data": … }`
- **Pagination** — opaque cursor in `meta.next_cursor`
- A free **sandbox** key covers the race & track endpoints over a rolling 7-day window; the **live** plan opens every endpoint and the full archive.

Get a key and read the docs at **[greyhoundapi.com](https://greyhoundapi.com)**.

## Contributing

Issues and pull requests are welcome. Please keep changes focused and, for the
JavaScript client, avoid adding runtime dependencies — it's intentionally
dependency-free and build-step-free.

## License

[MIT](./LICENSE)
