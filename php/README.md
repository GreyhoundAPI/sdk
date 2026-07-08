# greyhoundapi/sdk (PHP)

Official PHP client for the **[GreyhoundAPI](https://greyhoundapi.com)** — greyhound racing data for Great Britain and Australia: racecards, results, sectionals, starting prices and Betfair markets, plus dogs, trainers, owners, tracks and breeding.

- **PHP 7.4+**, requires only `ext-curl` and `ext-json`
- No third-party dependencies
- Covers all 35 REST endpoints, with cursor pagination and a typed exception

## Install

```sh
composer require greyhoundapi/sdk
```

No Composer? Clone the repo and `require` the two files in `src/` directly — there are no dependencies.

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

```php
<?php
require 'vendor/autoload.php';

use GreyhoundApi\Client;

$gapi = new Client(getenv('GREYHOUNDAPI_KEY'));

// today's GB racecards
$card = $gapi->racecardsToday(['region' => 'GB']);

// one race, fully resolved
$race = $gapi->race(1229082);
echo $race['data']['runners'][0]['dog_name'];
```

Every call returns the decoded envelope as an associative array:

```php
['meta' => ['request_id' => 'req_…', 'data_as_of' => '2026-07-08T13:05:12Z'], 'data' => [/* … */]]
```

## Authentication

Pass your key to the constructor; it's sent on every request as `X-API-Key`.

```php
$gapi = new Client('gapi_live_...', [
    // 'base_url' => 'https://api.greyhoundapi.com/v1',
    // 'timeout'  => 30.0,
]);
```

## Methods

Query parameters are passed as an array; path parameters are positional. Each
method returns the response envelope as an array.

| Method | Endpoint |
| --- | --- |
| `racecardsToday($params)` | `GET /racecards/today` |
| `racecardsUpcoming($params)` | `GET /racecards/upcoming` |
| `races($params)` | `GET /races` |
| `race($raceId)` | `GET /races/{race_id}` |
| `raceRunners($raceId, $params)` | `GET /races/{race_id}/runners` |
| `raceResult($raceId)` | `GET /races/{race_id}/result` |
| `raceStatus($raceId)` | `GET /races/{race_id}/status` |
| `raceMarket($raceId)` | `GET /races/{race_id}/market` |
| `resultsToday($params)` | `GET /results/today` |
| `results($params)` | `GET /results` |
| `latestResults($params)` | `GET /results/latest` |
| `meetings($params)` | `GET /meetings` |
| `meetingsToday($params)` | `GET /meetings/today` |
| `meeting($meetingId)` | `GET /meetings/{meeting_id}` |
| `searchDogs($params)` | `GET /dogs/search` |
| `dog($dogId)` | `GET /dogs/{dog_id}` |
| `dogForm($dogId, $params)` | `GET /dogs/{dog_id}/form` |
| `dogEntries($dogId)` | `GET /dogs/{dog_id}/entries` |
| `dogPrices($dogId)` | `GET /dogs/{dog_id}/prices` |
| `dogHeadToHead($dogId, $rivalId)` | `GET /dogs/{dog_id}/head-to-head/{rival_id}` |
| `searchTrainers($params)` | `GET /trainers/search` |
| `trainer($trainerId)` | `GET /trainers/{trainer_id}` |
| `trainerRunners($trainerId)` | `GET /trainers/{trainer_id}/runners` |
| `trainerResults($trainerId, $params)` | `GET /trainers/{trainer_id}/results` |
| `searchOwners($params)` | `GET /owners/search` |
| `owner($ownerId)` | `GET /owners/{owner_id}` |
| `tracks($params)` | `GET /tracks` |
| `track($trackId)` | `GET /tracks/{track_id}` |
| `trackRaces($trackId, $params)` | `GET /tracks/{track_id}/races` |
| `trackStats($trackId)` | `GET /tracks/{track_id}/stats` |
| `sireProgeny($name, $params)` | `GET /sires/{name}/progeny` |
| `damProgeny($name, $params)` | `GET /dams/{name}/progeny` |
| `status()` | `GET /status` |
| `usage()` | `GET /usage` |
| `reference($params)` | `GET /reference` |

### Common parameters

Most list endpoints accept: `region` (`GB` \| `AU`), `date_from` / `date_to`
(`YYYY-MM-DD`, track-local), `track_id`, `grade`, `distance_m`, `limit` (1–200)
and `cursor`. See the [docs](https://greyhoundapi.com/documentation).

```php
$races = $gapi->races(['region' => 'GB', 'date_from' => '2026-07-01', 'grade' => 'A2', 'limit' => 100]);
```

## Pagination

List endpoints page with an opaque cursor in `meta.next_cursor`. Pass `cursor`
yourself, or let the SDK walk every page and yield each item:

```php
foreach ($gapi->paginate('/races', ['region' => 'GB', 'date_from' => '2026-07-01']) as $race) {
    echo $race['race_id'], "\n";
}
```

## Errors

Any non-2xx response (or transport failure) throws `GreyhoundApi\ApiException`:

```php
use GreyhoundApi\ApiException;

try {
    $gapi->race(999999999);
} catch (ApiException $e) {
    echo $e->status, ' ', $e->code, ' ', $e->getMessage(), ' ', $e->requestId;
}
```

## Links

- **Documentation** — https://greyhoundapi.com/documentation
- **Service status** — https://greyhoundapi.com/status
- **Pricing** — https://greyhoundapi.com/pricing

## License

[MIT](./LICENSE)
