// Type definitions for @greyhoundapi/sdk

export interface GreyhoundAPIOptions {
  /** Your API key (live or sandbox). */
  apiKey: string;
  /** Override the API base URL. Defaults to https://api.greyhoundapi.com/v1 */
  baseUrl?: string;
  /** Custom fetch implementation (for Node < 18 or testing). */
  fetch?: typeof fetch;
  /** Per-request timeout in milliseconds. Defaults to 30000. */
  timeout?: number;
}

export interface Meta {
  request_id?: string;
  data_as_of?: string;
  count?: number;
  next_cursor?: string | null;
  [key: string]: unknown;
}

/** The standard response envelope. `data` shape varies per endpoint. */
export interface Envelope<T = any> {
  meta: Meta;
  data: T;
}

export type Params = Record<string, string | number | boolean | null | undefined>;

export class GreyhoundAPIError extends Error {
  name: "GreyhoundAPIError";
  status?: number;
  code?: string;
  requestId?: string;
  details?: any;
}

export class GreyhoundAPI {
  constructor(options: GreyhoundAPIOptions);

  apiKey: string;
  baseUrl: string;
  timeout: number;

  racecards: {
    today(params?: Params): Promise<Envelope>;
    upcoming(params?: Params): Promise<Envelope>;
  };
  races: {
    list(params?: Params): Promise<Envelope>;
    get(raceId: number | string): Promise<Envelope>;
    runners(raceId: number | string, params?: Params): Promise<Envelope>;
    result(raceId: number | string): Promise<Envelope>;
    status(raceId: number | string): Promise<Envelope>;
    market(raceId: number | string): Promise<Envelope>;
  };
  results: {
    today(params?: Params): Promise<Envelope>;
    search(params?: Params): Promise<Envelope>;
    latest(params?: Params): Promise<Envelope>;
  };
  meetings: {
    list(params?: Params): Promise<Envelope>;
    today(params?: Params): Promise<Envelope>;
    get(meetingId: number | string): Promise<Envelope>;
  };
  dogs: {
    search(params?: Params): Promise<Envelope>;
    get(dogId: number | string): Promise<Envelope>;
    form(dogId: number | string, params?: Params): Promise<Envelope>;
    entries(dogId: number | string): Promise<Envelope>;
    prices(dogId: number | string): Promise<Envelope>;
    headToHead(dogId: number | string, rivalId: number | string): Promise<Envelope>;
  };
  trainers: {
    search(params?: Params): Promise<Envelope>;
    get(trainerId: number | string): Promise<Envelope>;
    runners(trainerId: number | string): Promise<Envelope>;
    results(trainerId: number | string, params?: Params): Promise<Envelope>;
  };
  owners: {
    search(params?: Params): Promise<Envelope>;
    get(ownerId: number | string): Promise<Envelope>;
  };
  tracks: {
    list(params?: Params): Promise<Envelope>;
    get(trackId: number | string): Promise<Envelope>;
    races(trackId: number | string, params?: Params): Promise<Envelope>;
    stats(trackId: number | string): Promise<Envelope>;
  };
  sires: { progeny(name: string, params?: Params): Promise<Envelope> };
  dams: { progeny(name: string, params?: Params): Promise<Envelope> };

  status(): Promise<Envelope>;
  usage(): Promise<Envelope>;
  reference(params?: Params): Promise<Envelope>;

  request(method: string, path: string, opts?: { params?: Params }): Promise<Envelope>;
  paginate(path: string, params?: Params): AsyncGenerator<any, void, unknown>;
}

export default GreyhoundAPI;
