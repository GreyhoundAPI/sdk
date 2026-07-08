package greyhoundapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		q := r.URL.Query()
		switch {
		case r.URL.Path == "/v1/races" && q.Get("cursor") == "":
			w.Write([]byte(`{"meta":{"next_cursor":"p2"},"data":[{"race_id":1},{"race_id":2}]}`))
		case r.URL.Path == "/v1/races":
			w.Write([]byte(`{"meta":{"next_cursor":null},"data":[{"race_id":3}]}`))
		case r.URL.Path == "/v1/fail":
			w.WriteHeader(404)
			w.Write([]byte(`{"meta":{"request_id":"req_e"},"error":{"code":"not_found","message":"nope"}}`))
		default:
			d := map[string]string{"path": r.URL.Path, "auth": r.Header.Get("X-API-Key"), "query": r.URL.RawQuery}
			b, _ := json.Marshal(map[string]interface{}{"meta": map[string]string{"request_id": "req_x"}, "data": d})
			w.Write(b)
		}
	}))
	defer srv.Close()

	c := New("KEY123")
	c.BaseURL = srv.URL + "/v1"
	dec := func(env *Envelope) map[string]string {
		var m map[string]string
		json.Unmarshal(env.Data, &m)
		return m
	}

	env, err := c.Race(1229082)
	if err != nil {
		t.Fatal(err)
	}
	if p := dec(env)["path"]; p != "/v1/races/1229082" {
		t.Fatalf("race path %q", p)
	}
	env, err = c.RacecardsToday(Params{"region": "GB"})
	if err != nil {
		t.Fatal(err)
	}
	m := dec(env)
	if !strings.Contains(m["query"], "region=GB") || m["auth"] != "KEY123" {
		t.Fatalf("query/auth %v", m)
	}
	env, _ = c.DogHeadToHead(655044, 654595)
	if p := dec(env)["path"]; p != "/v1/dogs/655044/head-to-head/654595" {
		t.Fatalf("h2h %q", p)
	}
	var ids []int
	err = c.Paginate("/races", Params{"region": "GB"}, func(item json.RawMessage) error {
		var r struct {
			RaceID int `json:"race_id"`
		}
		json.Unmarshal(item, &r)
		ids = append(ids, r.RaceID)
		return nil
	})
	if err != nil || len(ids) != 3 || ids[0] != 1 || ids[2] != 3 {
		t.Fatalf("paginate %v err %v", ids, err)
	}
	_, err = c.Do("GET", "/fail", nil)
	ae, ok := err.(*Error)
	if !ok || ae.Status != 404 || ae.Code != "not_found" {
		t.Fatalf("err %v", err)
	}
	t.Log("Go SDK OK — encoding, query+auth, pagination, errors verified")
}
