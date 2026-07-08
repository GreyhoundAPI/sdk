// Package greyhoundapi is the official Go client for the GreyhoundAPI:
// greyhound racing data for Great Britain and Australia (racecards, results,
// sectionals, starting prices and Betfair markets, plus dogs, trainers,
// owners, tracks and breeding).
//
//	gapi := greyhoundapi.New(os.Getenv("GREYHOUNDAPI_KEY"))
//	race, err := gapi.Race(1229082)
//
// Standard library only. Docs: https://greyhoundapi.com/documentation
package greyhoundapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// DefaultBaseURL is the production API base URL.
const DefaultBaseURL = "https://api.greyhoundapi.com/v1"

// Params holds query parameters for a request. Values are stringified; nil or
// empty values are skipped.
type Params map[string]interface{}

// Meta is the response envelope metadata.
type Meta struct {
	RequestID  string  `json:"request_id"`
	DataAsOf   string  `json:"data_as_of"`
	Count      int     `json:"count"`
	NextCursor *string `json:"next_cursor"`
}

// Envelope is the standard response shape. Data is raw JSON — unmarshal it into
// your own type: json.Unmarshal(env.Data, &v).
type Envelope struct {
	Meta Meta            `json:"meta"`
	Data json.RawMessage `json:"data"`
}

// Error is returned for any non-2xx response (or a transport failure).
type Error struct {
	Status    int
	Code      string
	Message   string
	RequestID string
}

func (e *Error) Error() string {
	if e.Code != "" {
		return fmt.Sprintf("greyhoundapi: %d %s: %s", e.Status, e.Code, e.Message)
	}
	return fmt.Sprintf("greyhoundapi: %d: %s", e.Status, e.Message)
}

// Client is a GreyhoundAPI client. The zero value is not usable; use New.
type Client struct {
	APIKey  string
	BaseURL string
	HTTP    *http.Client
}

// New returns a Client with the given API key and sensible defaults. Override
// BaseURL or HTTP afterwards if you need to.
func New(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: DefaultBaseURL,
		HTTP:    &http.Client{Timeout: 30 * time.Second},
	}
}

// Do performs a request and returns the decoded envelope.
func (c *Client) Do(method, path string, params Params) (*Envelope, error) {
	u := strings.TrimRight(c.BaseURL, "/") + path
	if len(params) > 0 {
		q := url.Values{}
		for k, v := range params {
			if v == nil {
				continue
			}
			s := fmt.Sprint(v)
			if s == "" {
				continue
			}
			q.Set(k, s)
		}
		if enc := q.Encode(); enc != "" {
			u += "?" + enc
		}
	}

	req, err := http.NewRequest(method, u, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", c.APIKey)
	req.Header.Set("Accept", "application/json")

	httpClient := c.HTTP
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, &Error{Code: "network_error", Message: err.Error()}
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var e struct {
			Error struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			} `json:"error"`
			Meta struct {
				RequestID string `json:"request_id"`
			} `json:"meta"`
		}
		_ = json.Unmarshal(body, &e)
		msg := e.Error.Message
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return nil, &Error{Status: resp.StatusCode, Code: e.Error.Code, Message: msg, RequestID: e.Meta.RequestID}
	}

	var env Envelope
	if len(body) > 0 {
		if err := json.Unmarshal(body, &env); err != nil {
			return nil, err
		}
	}
	return &env, nil
}

func esc(v interface{}) string { return url.PathEscape(fmt.Sprint(v)) }

// ---- racecards ----
func (c *Client) RacecardsToday(p Params) (*Envelope, error) {
	return c.Do("GET", "/racecards/today", p)
}
func (c *Client) RacecardsUpcoming(p Params) (*Envelope, error) {
	return c.Do("GET", "/racecards/upcoming", p)
}

// ---- races ----
func (c *Client) Races(p Params) (*Envelope, error) { return c.Do("GET", "/races", p) }
func (c *Client) Race(raceID interface{}) (*Envelope, error) {
	return c.Do("GET", "/races/"+esc(raceID), nil)
}
func (c *Client) RaceRunners(raceID interface{}, p Params) (*Envelope, error) {
	return c.Do("GET", "/races/"+esc(raceID)+"/runners", p)
}
func (c *Client) RaceResult(raceID interface{}) (*Envelope, error) {
	return c.Do("GET", "/races/"+esc(raceID)+"/result", nil)
}
func (c *Client) RaceStatus(raceID interface{}) (*Envelope, error) {
	return c.Do("GET", "/races/"+esc(raceID)+"/status", nil)
}
func (c *Client) RaceMarket(raceID interface{}) (*Envelope, error) {
	return c.Do("GET", "/races/"+esc(raceID)+"/market", nil)
}

// ---- results ----
func (c *Client) ResultsToday(p Params) (*Envelope, error)  { return c.Do("GET", "/results/today", p) }
func (c *Client) Results(p Params) (*Envelope, error)       { return c.Do("GET", "/results", p) }
func (c *Client) LatestResults(p Params) (*Envelope, error) { return c.Do("GET", "/results/latest", p) }

// ---- meetings ----
func (c *Client) Meetings(p Params) (*Envelope, error)      { return c.Do("GET", "/meetings", p) }
func (c *Client) MeetingsToday(p Params) (*Envelope, error) { return c.Do("GET", "/meetings/today", p) }
func (c *Client) Meeting(meetingID interface{}) (*Envelope, error) {
	return c.Do("GET", "/meetings/"+esc(meetingID), nil)
}

// ---- dogs ----
func (c *Client) SearchDogs(p Params) (*Envelope, error) { return c.Do("GET", "/dogs/search", p) }
func (c *Client) Dog(dogID interface{}) (*Envelope, error) {
	return c.Do("GET", "/dogs/"+esc(dogID), nil)
}
func (c *Client) DogForm(dogID interface{}, p Params) (*Envelope, error) {
	return c.Do("GET", "/dogs/"+esc(dogID)+"/form", p)
}
func (c *Client) DogEntries(dogID interface{}) (*Envelope, error) {
	return c.Do("GET", "/dogs/"+esc(dogID)+"/entries", nil)
}
func (c *Client) DogPrices(dogID interface{}) (*Envelope, error) {
	return c.Do("GET", "/dogs/"+esc(dogID)+"/prices", nil)
}
func (c *Client) DogHeadToHead(dogID, rivalID interface{}) (*Envelope, error) {
	return c.Do("GET", "/dogs/"+esc(dogID)+"/head-to-head/"+esc(rivalID), nil)
}

// ---- trainers ----
func (c *Client) SearchTrainers(p Params) (*Envelope, error) {
	return c.Do("GET", "/trainers/search", p)
}
func (c *Client) Trainer(trainerID interface{}) (*Envelope, error) {
	return c.Do("GET", "/trainers/"+esc(trainerID), nil)
}
func (c *Client) TrainerRunners(trainerID interface{}) (*Envelope, error) {
	return c.Do("GET", "/trainers/"+esc(trainerID)+"/runners", nil)
}
func (c *Client) TrainerResults(trainerID interface{}, p Params) (*Envelope, error) {
	return c.Do("GET", "/trainers/"+esc(trainerID)+"/results", p)
}

// ---- owners ----
func (c *Client) SearchOwners(p Params) (*Envelope, error) { return c.Do("GET", "/owners/search", p) }
func (c *Client) Owner(ownerID interface{}) (*Envelope, error) {
	return c.Do("GET", "/owners/"+esc(ownerID), nil)
}

// ---- tracks ----
func (c *Client) Tracks(p Params) (*Envelope, error) { return c.Do("GET", "/tracks", p) }
func (c *Client) Track(trackID interface{}) (*Envelope, error) {
	return c.Do("GET", "/tracks/"+esc(trackID), nil)
}
func (c *Client) TrackRaces(trackID interface{}, p Params) (*Envelope, error) {
	return c.Do("GET", "/tracks/"+esc(trackID)+"/races", p)
}
func (c *Client) TrackStats(trackID interface{}) (*Envelope, error) {
	return c.Do("GET", "/tracks/"+esc(trackID)+"/stats", nil)
}

// ---- breeding ----
func (c *Client) SireProgeny(name string, p Params) (*Envelope, error) {
	return c.Do("GET", "/sires/"+esc(name)+"/progeny", p)
}
func (c *Client) DamProgeny(name string, p Params) (*Envelope, error) {
	return c.Do("GET", "/dams/"+esc(name)+"/progeny", p)
}

// ---- platform & reference ----
func (c *Client) Status() (*Envelope, error)            { return c.Do("GET", "/status", nil) }
func (c *Client) Usage() (*Envelope, error)             { return c.Do("GET", "/usage", nil) }
func (c *Client) Reference(p Params) (*Envelope, error) { return c.Do("GET", "/reference", p) }

// Paginate follows meta.next_cursor across every page, calling fn for each
// item. Return a non-nil error from fn to stop early (it is returned as-is).
//
//	err := gapi.Paginate("/races", greyhoundapi.Params{"region": "GB"}, func(item json.RawMessage) error {
//	    var r struct{ RaceID int `json:"race_id"` }
//	    json.Unmarshal(item, &r)
//	    fmt.Println(r.RaceID)
//	    return nil
//	})
func (c *Client) Paginate(path string, params Params, fn func(item json.RawMessage) error) error {
	var cursor interface{}
	if params != nil {
		if v, ok := params["cursor"]; ok {
			cursor = v
		}
	}
	for {
		p := Params{}
		for k, v := range params {
			p[k] = v
		}
		p["cursor"] = cursor

		env, err := c.Do("GET", path, p)
		if err != nil {
			return err
		}

		var items []json.RawMessage
		if len(env.Data) > 0 {
			if err := json.Unmarshal(env.Data, &items); err != nil {
				var obj struct {
					Items []json.RawMessage `json:"items"`
				}
				_ = json.Unmarshal(env.Data, &obj)
				items = obj.Items
			}
		}
		for _, it := range items {
			if err := fn(it); err != nil {
				return err
			}
		}
		if env.Meta.NextCursor == nil || *env.Meta.NextCursor == "" {
			return nil
		}
		cursor = *env.Meta.NextCursor
	}
}
