"""
greyhoundapi — official Python client for the GreyhoundAPI: greyhound racing
data for Great Britain and Australia.

    from greyhoundapi import GreyhoundAPI

    gapi = GreyhoundAPI(api_key="gapi_live_...")
    card = gapi.racecards_today(region="GB")
    race = gapi.race(1229082)
    print(race["data"]["runners"][0]["dog_name"])

Zero dependencies — standard library only. Docs: https://greyhoundapi.com/documentation
"""
from __future__ import annotations

import json
import urllib.error
import urllib.parse
import urllib.request
from typing import Any, Dict, Iterator, Optional

__version__ = "1.0.0"
__all__ = ["GreyhoundAPI", "GreyhoundAPIError"]

DEFAULT_BASE_URL = "https://api.greyhoundapi.com/v1"


class GreyhoundAPIError(Exception):
    """Raised for any non-2xx response, or a transport failure."""

    def __init__(self, message: str, status: Optional[int] = None, code: Optional[str] = None,
                 request_id: Optional[str] = None, details: Optional[dict] = None):
        super().__init__(message)
        self.status = status
        self.code = code
        self.request_id = request_id
        self.details = details or {}


class GreyhoundAPI:
    """Client for the GreyhoundAPI.

    Args:
        api_key:  Your API key (live or sandbox).
        base_url: Override the API base URL.
        timeout:  Per-request timeout in seconds.
    """

    def __init__(self, api_key: str, base_url: str = DEFAULT_BASE_URL, timeout: float = 30.0):
        if not api_key:
            raise ValueError("GreyhoundAPI: api_key is required.")
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout

    # ------------------------------------------------------------------ #
    # Low-level
    # ------------------------------------------------------------------ #
    def _request(self, method: str, path: str, params: Optional[Dict[str, Any]] = None) -> Dict[str, Any]:
        url = self.base_url + path
        if params:
            clean = {k: v for k, v in params.items() if v is not None and v != ""}
            if clean:
                url += "?" + urllib.parse.urlencode(clean)
        req = urllib.request.Request(url, method=method)
        req.add_header("X-API-Key", self.api_key)
        req.add_header("Accept", "application/json")
        req.add_header("User-Agent", "greyhoundapi-python/" + __version__)
        try:
            with urllib.request.urlopen(req, timeout=self.timeout) as resp:
                raw = resp.read().decode("utf-8")
                return json.loads(raw) if raw else {}
        except urllib.error.HTTPError as exc:
            raw = exc.read().decode("utf-8", "replace")
            body: Any = {}
            try:
                body = json.loads(raw)
            except ValueError:
                pass
            err = body.get("error", {}) if isinstance(body, dict) else {}
            meta = body.get("meta", {}) if isinstance(body, dict) else {}
            message = err.get("message")
            if not message:
                snippet = " ".join(raw.split())[:200]
                message = "HTTP {}{}".format(exc.code, (": " + snippet) if snippet else "")
            raise GreyhoundAPIError(
                message,
                status=exc.code,
                code=err.get("code"),
                request_id=meta.get("request_id"),
                details=err,
            ) from None
        except urllib.error.URLError as exc:
            raise GreyhoundAPIError("Request failed: {}".format(exc.reason), code="network_error") from None

    def _get(self, path: str, **params: Any) -> Dict[str, Any]:
        return self._request("GET", path, params)

    @staticmethod
    def _enc(value: Any) -> str:
        return urllib.parse.quote(str(value), safe="")

    # ------------------------------------------------------------------ #
    # Racecards
    # ------------------------------------------------------------------ #
    def racecards_today(self, **params: Any) -> Dict[str, Any]:
        return self._get("/racecards/today", **params)

    def racecards_upcoming(self, **params: Any) -> Dict[str, Any]:
        return self._get("/racecards/upcoming", **params)

    # ------------------------------------------------------------------ #
    # Races
    # ------------------------------------------------------------------ #
    def races(self, **params: Any) -> Dict[str, Any]:
        return self._get("/races", **params)

    def race(self, race_id: Any) -> Dict[str, Any]:
        return self._get("/races/{}".format(self._enc(race_id)))

    def race_runners(self, race_id: Any, **params: Any) -> Dict[str, Any]:
        return self._get("/races/{}/runners".format(self._enc(race_id)), **params)

    def race_result(self, race_id: Any) -> Dict[str, Any]:
        return self._get("/races/{}/result".format(self._enc(race_id)))

    def race_status(self, race_id: Any) -> Dict[str, Any]:
        return self._get("/races/{}/status".format(self._enc(race_id)))

    def race_market(self, race_id: Any) -> Dict[str, Any]:
        return self._get("/races/{}/market".format(self._enc(race_id)))

    # ------------------------------------------------------------------ #
    # Results
    # ------------------------------------------------------------------ #
    def results_today(self, **params: Any) -> Dict[str, Any]:
        return self._get("/results/today", **params)

    def results(self, **params: Any) -> Dict[str, Any]:
        return self._get("/results", **params)

    def latest_results(self, **params: Any) -> Dict[str, Any]:
        return self._get("/results/latest", **params)

    # ------------------------------------------------------------------ #
    # Meetings
    # ------------------------------------------------------------------ #
    def meetings(self, **params: Any) -> Dict[str, Any]:
        return self._get("/meetings", **params)

    def meetings_today(self, **params: Any) -> Dict[str, Any]:
        return self._get("/meetings/today", **params)

    def meeting(self, meeting_id: Any) -> Dict[str, Any]:
        return self._get("/meetings/{}".format(self._enc(meeting_id)))

    # ------------------------------------------------------------------ #
    # Dogs
    # ------------------------------------------------------------------ #
    def search_dogs(self, **params: Any) -> Dict[str, Any]:
        return self._get("/dogs/search", **params)

    def dog(self, dog_id: Any) -> Dict[str, Any]:
        return self._get("/dogs/{}".format(self._enc(dog_id)))

    def dog_form(self, dog_id: Any, **params: Any) -> Dict[str, Any]:
        return self._get("/dogs/{}/form".format(self._enc(dog_id)), **params)

    def dog_entries(self, dog_id: Any) -> Dict[str, Any]:
        return self._get("/dogs/{}/entries".format(self._enc(dog_id)))

    def dog_prices(self, dog_id: Any) -> Dict[str, Any]:
        return self._get("/dogs/{}/prices".format(self._enc(dog_id)))

    def dog_head_to_head(self, dog_id: Any, rival_id: Any) -> Dict[str, Any]:
        return self._get("/dogs/{}/head-to-head/{}".format(self._enc(dog_id), self._enc(rival_id)))

    # ------------------------------------------------------------------ #
    # Trainers
    # ------------------------------------------------------------------ #
    def search_trainers(self, **params: Any) -> Dict[str, Any]:
        return self._get("/trainers/search", **params)

    def trainer(self, trainer_id: Any) -> Dict[str, Any]:
        return self._get("/trainers/{}".format(self._enc(trainer_id)))

    def trainer_runners(self, trainer_id: Any) -> Dict[str, Any]:
        return self._get("/trainers/{}/runners".format(self._enc(trainer_id)))

    def trainer_results(self, trainer_id: Any, **params: Any) -> Dict[str, Any]:
        return self._get("/trainers/{}/results".format(self._enc(trainer_id)), **params)

    # ------------------------------------------------------------------ #
    # Owners
    # ------------------------------------------------------------------ #
    def search_owners(self, **params: Any) -> Dict[str, Any]:
        return self._get("/owners/search", **params)

    def owner(self, owner_id: Any) -> Dict[str, Any]:
        return self._get("/owners/{}".format(self._enc(owner_id)))

    # ------------------------------------------------------------------ #
    # Tracks
    # ------------------------------------------------------------------ #
    def tracks(self, **params: Any) -> Dict[str, Any]:
        return self._get("/tracks", **params)

    def track(self, track_id: Any) -> Dict[str, Any]:
        return self._get("/tracks/{}".format(self._enc(track_id)))

    def track_races(self, track_id: Any, **params: Any) -> Dict[str, Any]:
        return self._get("/tracks/{}/races".format(self._enc(track_id)), **params)

    def track_stats(self, track_id: Any) -> Dict[str, Any]:
        return self._get("/tracks/{}/stats".format(self._enc(track_id)))

    # ------------------------------------------------------------------ #
    # Breeding
    # ------------------------------------------------------------------ #
    def sire_progeny(self, name: str, **params: Any) -> Dict[str, Any]:
        return self._get("/sires/{}/progeny".format(self._enc(name)), **params)

    def dam_progeny(self, name: str, **params: Any) -> Dict[str, Any]:
        return self._get("/dams/{}/progeny".format(self._enc(name)), **params)

    # ------------------------------------------------------------------ #
    # Platform & reference
    # ------------------------------------------------------------------ #
    def status(self) -> Dict[str, Any]:
        return self._get("/status")

    def usage(self) -> Dict[str, Any]:
        return self._get("/usage")

    def reference(self, **params: Any) -> Dict[str, Any]:
        return self._get("/reference", **params)

    # ------------------------------------------------------------------ #
    # Pagination helper
    # ------------------------------------------------------------------ #
    def paginate(self, path: str, **params: Any) -> Iterator[Any]:
        """Follow ``meta.next_cursor`` across every page, yielding each item.

            for race in gapi.paginate("/races", region="GB"):
                ...
        """
        cursor = params.pop("cursor", None)
        while True:
            page = self._get(path, cursor=cursor, **params)
            data = page.get("data")
            if isinstance(data, list):
                items = data
            elif isinstance(data, dict):
                items = data.get("items", [])
            else:
                items = []
            for item in items:
                yield item
            cursor = (page.get("meta") or {}).get("next_cursor")
            if not cursor:
                break
