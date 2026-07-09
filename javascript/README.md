# @greyhoundapi/sdk

Official JavaScript / TypeScript client for the **[GreyhoundAPI](https://greyhoundapi.com)** — greyhound racing data for Great Britain and Australia: racecards, results, sectionals, starting prices and Betfair markets, plus dogs, trainers, owners, tracks and breeding.

- Works in **Node 18+** and modern **browsers** (uses the global `fetch`)
- **Zero dependencies**
- **TypeScript** types included
- Covers all 35 REST endpoints, with cursor pagination and typed errors

## Install

```sh
npm install @greyhoundapi/sdk
```

## Quick start

```js
import { GreyhoundAPI } from "@greyhoundapi/sdk";

const gapi = new GreyhoundAPI({ apiKey: process.env.GREYHOUNDAPI_KEY });

// today's GB racecards
const card = await gapi.racecards.today({ region: "GB" });

// one race, fully resolved
const race = await gapi.races.get(1229082);
console.log(race.data.runners[0].dog_name);
```

Every call resolves to the standard envelope:

```json
{ "meta": { "request_id": "req_…", "data_as_of": "2026-07-08T13:05:12Z" }, "data": { } }
```

## Getting an API key

The API is authenticated with a key sent in the `X-API-Key` header. To get one:

1. Go to **[greyhoundapi.com](https://greyhoundapi.com)** and click **Get a free key** (top-right).
2. **Sign in** to create your account.
3. Open **Account → API keys** and click **Create key**.
4. Copy the key — **it's shown only once**. Sandbox keys start `gapi_test_`; live keys start `gapi_live_`.

**Free sandbox key** — no card required: the race & track endpoints over a rolling 7-day window, 500 requests/day, one active key. Ideal for building and evaluating.

**Live plan** ($99/month) — every endpoint, the full historical archive, live-day data and the WebSocket results stream, 250,000 requests/month, and up to 5 active keys. See **[pricing](https://greyhoundapi.com/pricing)**.

**Keep your key secret.** Load it from an environment variable (e.g. `GREYHOUNDAPI_KEY`) rather than committing it. In browser code the key is visible to users — use a sandbox key for public demos, or proxy requests through your own server.

## Authentication

Pass your key to the constructor; it's sent on every request as the `X-API-Key` header.

```js
const gapi = new GreyhoundAPI({
  apiKey: "gapi_live_…",   // required
  // baseUrl: "https://api.greyhoundapi.com/v1",  // optional override
  // fetch: myFetch,        // optional (Node < 18, or for testing)
  // timeout: 30000,        // optional, ms
});
```

## Methods

Methods are grouped into namespaces that mirror the API. Each takes an optional
`params` object for query parameters (see [Common parameters](#common-parameters))
and returns a `Promise` of the response envelope.

| Method | Endpoint |
| --- | --- |
| `racecards.today(params?)` | `GET /racecards/today` |
| `racecards.upcoming(params?)` | `GET /racecards/upcoming` |
| `races.list(params?)` | `GET /races` |
| `races.get(raceId)` | `GET /races/{race_id}` |
| `races.runners(raceId, params?)` | `GET /races/{race_id}/runners` |
| `races.result(raceId)` | `GET /races/{race_id}/result` |
| `races.status(raceId)` | `GET /races/{race_id}/status` |
| `races.market(raceId)` | `GET /races/{race_id}/market` |
| `results.today(params?)` | `GET /results/today` |
| `results.search(params?)` | `GET /results` |
| `results.latest(params?)` | `GET /results/latest` |
| `meetings.list(params?)` | `GET /meetings` |
| `meetings.today(params?)` | `GET /meetings/today` |
| `meetings.get(meetingId)` | `GET /meetings/{meeting_id}` |
| `dogs.search(params?)` | `GET /dogs/search` |
| `dogs.get(dogId)` | `GET /dogs/{dog_id}` |
| `dogs.form(dogId, params?)` | `GET /dogs/{dog_id}/form` |
| `dogs.entries(dogId)` | `GET /dogs/{dog_id}/entries` |
| `dogs.prices(dogId)` | `GET /dogs/{dog_id}/prices` |
| `dogs.headToHead(dogId, rivalId)` | `GET /dogs/{dog_id}/head-to-head/{rival_id}` |
| `trainers.search(params?)` | `GET /trainers/search` |
| `trainers.get(trainerId)` | `GET /trainers/{trainer_id}` |
| `trainers.runners(trainerId)` | `GET /trainers/{trainer_id}/runners` |
| `trainers.results(trainerId, params?)` | `GET /trainers/{trainer_id}/results` |
| `owners.search(params?)` | `GET /owners/search` |
| `owners.get(ownerId)` | `GET /owners/{owner_id}` |
| `tracks.list(params?)` | `GET /tracks` |
| `tracks.get(trackId)` | `GET /tracks/{track_id}` |
| `tracks.races(trackId, params?)` | `GET /tracks/{track_id}/races` |
| `tracks.stats(trackId)` | `GET /tracks/{track_id}/stats` |
| `sires.progeny(name, params?)` | `GET /sires/{name}/progeny` |
| `dams.progeny(name, params?)` | `GET /dams/{name}/progeny` |
| `status()` | `GET /status` |
| `usage()` | `GET /usage` |
| `reference(params?)` | `GET /reference` |

There's also a low-level escape hatch for anything not wrapped:

```js
const res = await gapi.request("GET", "/races", { params: { region: "AU" } });
```

### Common parameters

Most list endpoints accept: `region` (`GB` \| `AU`), `date_from` / `date_to`
(`YYYY-MM-DD`, track-local), `track_id`, `grade`, `distance_m`, `limit` (1–200),
and `cursor`. See the [docs](https://greyhoundapi.com/documentation) for what
each endpoint supports.

## Pagination

List endpoints page with an opaque cursor in `meta.next_cursor`. Pass `cursor`
yourself, or let the SDK walk every page and yield each item:

```js
for await (const race of gapi.paginate("/races", { region: "GB", date_from: "2026-07-01" })) {
  console.log(race.race_id);
}
```

## Errors

Any non-2xx response (or a transport failure) throws `GreyhoundAPIError`:

```js
import { GreyhoundAPI, GreyhoundAPIError } from "@greyhoundapi/sdk";

try {
  await gapi.races.get(999999999);
} catch (err) {
  if (err instanceof GreyhoundAPIError) {
    console.error(err.status, err.code, err.message, err.requestId);
  } else {
    throw err;
  }
}
```

## Links

- **Documentation** — https://greyhoundapi.com/documentation
- **Service status** — https://greyhoundapi.com/status
- **Pricing** — https://greyhoundapi.com/pricing

## License

[MIT](./LICENSE)
