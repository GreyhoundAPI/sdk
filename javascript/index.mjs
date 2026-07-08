/**
 * @greyhoundapi/sdk — official JavaScript / TypeScript client for the
 * GreyhoundAPI: greyhound racing data for Great Britain and Australia.
 *
 * Works in Node 18+ and modern browsers (uses the global `fetch`). Zero deps.
 *
 *   import { GreyhoundAPI } from "@greyhoundapi/sdk";
 *
 *   const gapi = new GreyhoundAPI({ apiKey: process.env.GREYHOUNDAPI_KEY });
 *   const card = await gapi.racecards.today({ region: "GB" });
 *   const race = await gapi.races.get(1229082);
 *   console.log(race.data.runners[0].dog_name);
 *
 * Docs: https://greyhoundapi.com/documentation
 */

const DEFAULT_BASE_URL = "https://api.greyhoundapi.com/v1";

/** Error thrown for any non-2xx response (or a transport failure). */
export class GreyhoundAPIError extends Error {
  constructor(message, { status, code, requestId, details } = {}) {
    super(message);
    this.name = "GreyhoundAPIError";
    this.status = status;
    this.code = code;
    this.requestId = requestId;
    this.details = details;
  }
}

export class GreyhoundAPI {
  /**
   * @param {object}        opts
   * @param {string}        opts.apiKey          Your API key (live or sandbox).
   * @param {string}        [opts.baseUrl]       Override the API base URL.
   * @param {typeof fetch}  [opts.fetch]         Custom fetch (for Node < 18 or testing).
   * @param {number}        [opts.timeout=30000] Per-request timeout, ms.
   */
  constructor({ apiKey, baseUrl = DEFAULT_BASE_URL, fetch: fetchImpl, timeout = 30000 } = {}) {
    if (!apiKey) throw new Error("GreyhoundAPI: `apiKey` is required.");
    this.apiKey = apiKey;
    this.baseUrl = String(baseUrl).replace(/\/+$/, "");
    this._fetch = fetchImpl || globalThis.fetch;
    if (!this._fetch) throw new Error("GreyhoundAPI: no global fetch — pass `fetch` (Node < 18).");
    this.timeout = timeout;

    // ---- resource namespaces (mirror the documentation sections) ----
    this.racecards = {
      today:    (params) => this.request("GET", "/racecards/today", { params }),
      upcoming: (params) => this.request("GET", "/racecards/upcoming", { params }),
    };
    this.races = {
      list:    (params)     => this.request("GET", "/races", { params }),
      get:     (raceId)     => this.request("GET", `/races/${enc(raceId)}`),
      runners: (raceId, p)  => this.request("GET", `/races/${enc(raceId)}/runners`, { params: p }),
      result:  (raceId)     => this.request("GET", `/races/${enc(raceId)}/result`),
      status:  (raceId)     => this.request("GET", `/races/${enc(raceId)}/status`),
      market:  (raceId)     => this.request("GET", `/races/${enc(raceId)}/market`),
    };
    this.results = {
      today:  (params) => this.request("GET", "/results/today", { params }),
      search: (params) => this.request("GET", "/results", { params }),
      latest: (params) => this.request("GET", "/results/latest", { params }),
    };
    this.meetings = {
      list:  (params)      => this.request("GET", "/meetings", { params }),
      today: (params)      => this.request("GET", "/meetings/today", { params }),
      get:   (meetingId)   => this.request("GET", `/meetings/${enc(meetingId)}`),
    };
    this.dogs = {
      search:     (params)          => this.request("GET", "/dogs/search", { params }),
      get:        (dogId)           => this.request("GET", `/dogs/${enc(dogId)}`),
      form:       (dogId, p)        => this.request("GET", `/dogs/${enc(dogId)}/form`, { params: p }),
      entries:    (dogId)           => this.request("GET", `/dogs/${enc(dogId)}/entries`),
      prices:     (dogId)           => this.request("GET", `/dogs/${enc(dogId)}/prices`),
      headToHead: (dogId, rivalId)  => this.request("GET", `/dogs/${enc(dogId)}/head-to-head/${enc(rivalId)}`),
    };
    this.trainers = {
      search:  (params)        => this.request("GET", "/trainers/search", { params }),
      get:     (trainerId)     => this.request("GET", `/trainers/${enc(trainerId)}`),
      runners: (trainerId)     => this.request("GET", `/trainers/${enc(trainerId)}/runners`),
      results: (trainerId, p)  => this.request("GET", `/trainers/${enc(trainerId)}/results`, { params: p }),
    };
    this.owners = {
      search: (params)   => this.request("GET", "/owners/search", { params }),
      get:    (ownerId)  => this.request("GET", `/owners/${enc(ownerId)}`),
    };
    this.tracks = {
      list:  (params)       => this.request("GET", "/tracks", { params }),
      get:   (trackId)      => this.request("GET", `/tracks/${enc(trackId)}`),
      races: (trackId, p)   => this.request("GET", `/tracks/${enc(trackId)}/races`, { params: p }),
      stats: (trackId)      => this.request("GET", `/tracks/${enc(trackId)}/stats`),
    };
    this.sires = { progeny: (name, p) => this.request("GET", `/sires/${enc(name)}/progeny`, { params: p }) };
    this.dams  = { progeny: (name, p) => this.request("GET", `/dams/${enc(name)}/progeny`, { params: p }) };
  }

  // ---- platform & reference (top-level) ----
  status()          { return this.request("GET", "/status"); }
  usage()           { return this.request("GET", "/usage"); }
  reference(params) { return this.request("GET", "/reference", { params }); }

  /**
   * Low-level request. Returns the parsed JSON envelope `{ meta, data }`.
   * Throws {@link GreyhoundAPIError} on a non-2xx response.
   * @param {string} method
   * @param {string} path   Path relative to the base URL, e.g. "/races".
   * @param {{ params?: Record<string, any> }} [opts]
   */
  async request(method, path, { params } = {}) {
    const url = new URL(this.baseUrl + path);
    if (params) {
      for (const [k, v] of Object.entries(params)) {
        if (v !== undefined && v !== null && v !== "") url.searchParams.set(k, String(v));
      }
    }
    const ctrl = new AbortController();
    const timer = setTimeout(() => ctrl.abort(), this.timeout);
    let res;
    try {
      res = await this._fetch(url, {
        method,
        headers: { "X-API-Key": this.apiKey, Accept: "application/json" },
        signal: ctrl.signal,
      });
    } catch (e) {
      clearTimeout(timer);
      const msg = e && e.name === "AbortError" ? `Request timed out after ${this.timeout}ms` : `Request failed: ${e && e.message}`;
      throw new GreyhoundAPIError(msg, { code: "network_error" });
    }
    clearTimeout(timer);

    let body = null;
    const text = await res.text();
    if (text) { try { body = JSON.parse(text); } catch { /* leave body null on non-JSON */ } }

    if (!res.ok) {
      const err = (body && body.error) || {};
      throw new GreyhoundAPIError(err.message || `HTTP ${res.status}`, {
        status: res.status,
        code: err.code,
        requestId: body && body.meta ? body.meta.request_id : undefined,
        details: err,
      });
    }
    return body;
  }

  /**
   * Follow `meta.next_cursor` across every page of a list endpoint, yielding
   * each item.
   *
   *   for await (const race of gapi.paginate("/races", { region: "GB" })) { ... }
   *
   * @param {string} path
   * @param {Record<string, any>} [params]
   */
  async *paginate(path, params = {}) {
    let cursor = params.cursor;
    for (;;) {
      const page = await this.request("GET", path, { params: { ...params, cursor } });
      const data = page && page.data;
      const items = Array.isArray(data) ? data : (data && Array.isArray(data.items) ? data.items : []);
      for (const item of items) yield item;
      cursor = page && page.meta ? page.meta.next_cursor : null;
      if (!cursor) break;
    }
  }
}

function enc(v) { return encodeURIComponent(String(v)); }

export default GreyhoundAPI;
