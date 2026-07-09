<?php

namespace GreyhoundApi;

/**
 * Official PHP client for the GreyhoundAPI: greyhound racing data for Great
 * Britain and Australia.
 *
 *   $gapi = new \GreyhoundApi\Client('gapi_live_...');
 *   $card = $gapi->racecardsToday(['region' => 'GB']);
 *   $race = $gapi->race(1229082);
 *   echo $race['data']['runners'][0]['dog_name'];
 *
 * Requires ext-curl and ext-json (both standard). Docs: https://greyhoundapi.com/documentation
 */
class Client
{
    /** @var string */
    private $apiKey;
    /** @var string */
    private $baseUrl;
    /** @var float */
    private $timeout;

    /**
     * @param string $apiKey  Your API key (live or sandbox).
     * @param array  $options ['base_url' => string, 'timeout' => float(seconds)]
     */
    public function __construct(string $apiKey, array $options = [])
    {
        if ($apiKey === '') {
            throw new \InvalidArgumentException('GreyhoundApi\\Client: apiKey is required.');
        }
        $this->apiKey  = $apiKey;
        $this->baseUrl = rtrim($options['base_url'] ?? 'https://api.greyhoundapi.com/v1', '/');
        $this->timeout = (float) ($options['timeout'] ?? 30.0);
    }

    // ---- racecards ----
    public function racecardsToday(array $params = []): array { return $this->get('/racecards/today', $params); }
    public function racecardsUpcoming(array $params = []): array { return $this->get('/racecards/upcoming', $params); }

    // ---- races ----
    public function races(array $params = []): array { return $this->get('/races', $params); }
    public function race($raceId): array { return $this->get('/races/' . self::enc($raceId)); }
    public function raceRunners($raceId, array $params = []): array { return $this->get('/races/' . self::enc($raceId) . '/runners', $params); }
    public function raceResult($raceId): array { return $this->get('/races/' . self::enc($raceId) . '/result'); }
    public function raceStatus($raceId): array { return $this->get('/races/' . self::enc($raceId) . '/status'); }
    public function raceMarket($raceId): array { return $this->get('/races/' . self::enc($raceId) . '/market'); }

    // ---- results ----
    public function resultsToday(array $params = []): array { return $this->get('/results/today', $params); }
    public function results(array $params = []): array { return $this->get('/results', $params); }
    public function latestResults(array $params = []): array { return $this->get('/results/latest', $params); }

    // ---- meetings ----
    public function meetings(array $params = []): array { return $this->get('/meetings', $params); }
    public function meetingsToday(array $params = []): array { return $this->get('/meetings/today', $params); }
    public function meeting($meetingId): array { return $this->get('/meetings/' . self::enc($meetingId)); }

    // ---- dogs ----
    public function searchDogs(array $params = []): array { return $this->get('/dogs/search', $params); }
    public function dog($dogId): array { return $this->get('/dogs/' . self::enc($dogId)); }
    public function dogForm($dogId, array $params = []): array { return $this->get('/dogs/' . self::enc($dogId) . '/form', $params); }
    public function dogEntries($dogId): array { return $this->get('/dogs/' . self::enc($dogId) . '/entries'); }
    public function dogPrices($dogId): array { return $this->get('/dogs/' . self::enc($dogId) . '/prices'); }
    public function dogHeadToHead($dogId, $rivalId): array { return $this->get('/dogs/' . self::enc($dogId) . '/head-to-head/' . self::enc($rivalId)); }

    // ---- trainers ----
    public function searchTrainers(array $params = []): array { return $this->get('/trainers/search', $params); }
    public function trainer($trainerId): array { return $this->get('/trainers/' . self::enc($trainerId)); }
    public function trainerRunners($trainerId): array { return $this->get('/trainers/' . self::enc($trainerId) . '/runners'); }
    public function trainerResults($trainerId, array $params = []): array { return $this->get('/trainers/' . self::enc($trainerId) . '/results', $params); }

    // ---- owners ----
    public function searchOwners(array $params = []): array { return $this->get('/owners/search', $params); }
    public function owner($ownerId): array { return $this->get('/owners/' . self::enc($ownerId)); }

    // ---- tracks ----
    public function tracks(array $params = []): array { return $this->get('/tracks', $params); }
    public function track($trackId): array { return $this->get('/tracks/' . self::enc($trackId)); }
    public function trackRaces($trackId, array $params = []): array { return $this->get('/tracks/' . self::enc($trackId) . '/races', $params); }
    public function trackStats($trackId): array { return $this->get('/tracks/' . self::enc($trackId) . '/stats'); }

    // ---- breeding ----
    public function sireProgeny($name, array $params = []): array { return $this->get('/sires/' . self::enc($name) . '/progeny', $params); }
    public function damProgeny($name, array $params = []): array { return $this->get('/dams/' . self::enc($name) . '/progeny', $params); }

    // ---- platform & reference ----
    public function status(): array { return $this->get('/status'); }
    public function usage(): array { return $this->get('/usage'); }
    public function reference(array $params = []): array { return $this->get('/reference', $params); }

    /** Convenience GET. */
    public function get(string $path, array $params = []): array
    {
        return $this->request('GET', $path, $params);
    }

    /**
     * Low-level request. Returns the decoded envelope as an associative array.
     * Throws ApiException on a non-2xx response or transport failure.
     */
    public function request(string $method, string $path, array $params = []): array
    {
        $url = $this->baseUrl . $path;
        $params = array_filter($params, static function ($v) { return $v !== null && $v !== ''; });
        if (!empty($params)) {
            $url .= '?' . http_build_query($params);
        }

        $ch = curl_init();
        curl_setopt_array($ch, [
            CURLOPT_URL            => $url,
            CURLOPT_CUSTOMREQUEST  => $method,
            CURLOPT_RETURNTRANSFER => true,
            CURLOPT_CONNECTTIMEOUT => 10,
            CURLOPT_TIMEOUT        => (int) ceil($this->timeout),
            CURLOPT_USERAGENT      => 'greyhoundapi-php/1.0.0',
            CURLOPT_HTTPHEADER     => ['X-API-Key: ' . $this->apiKey, 'Accept: application/json'],
        ]);
        $body   = curl_exec($ch);
        $status = (int) curl_getinfo($ch, CURLINFO_HTTP_CODE);
        $errno  = curl_errno($ch);
        $errmsg = curl_error($ch);
        curl_close($ch);

        if ($errno !== 0) {
            throw new ApiException('Request failed: ' . $errmsg, null, 'network_error');
        }

        $decoded = ($body !== false && $body !== '') ? json_decode($body, true) : null;

        if ($status < 200 || $status >= 300) {
            $err = (is_array($decoded) && isset($decoded['error'])) ? $decoded['error'] : [];
            $rid = $decoded['meta']['request_id'] ?? null;
            throw new ApiException($err['message'] ?? ('HTTP ' . $status), $status, $err['code'] ?? null, $rid, is_array($err) ? $err : []);
        }

        return is_array($decoded) ? $decoded : [];
    }

    /**
     * Yield every item across all pages, following meta.next_cursor.
     *
     *   foreach ($gapi->paginate('/races', ['region' => 'GB']) as $race) { ... }
     *
     * @return \Generator
     */
    public function paginate(string $path, array $params = []): \Generator
    {
        $cursor = $params['cursor'] ?? null;
        do {
            $page = $this->request('GET', $path, array_merge($params, ['cursor' => $cursor]));
            $data = $page['data'] ?? [];
            if (is_array($data)) {
                $isList = $data === [] || array_keys($data) === range(0, count($data) - 1);
                $items  = $isList ? $data : ($data['items'] ?? []);
            } else {
                $items = [];
            }
            foreach ($items as $item) {
                yield $item;
            }
            $cursor = $page['meta']['next_cursor'] ?? null;
        } while ($cursor);
    }

    private static function enc($value): string
    {
        return rawurlencode((string) $value);
    }
}
