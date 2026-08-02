package openfga

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/reefops/reefops/services/authorizer/internal/authorization"
)

func TestCheckUsesPinnedModelAndContextualTuple(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/stores/store/check" || r.Header.Get("authorization") != "Bearer secret" {
			t.Errorf("unexpected request: %s %q", r.URL.Path, r.Header.Get("authorization"))
		}
		var body checkRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Error(err)
		}
		if body.AuthorizationModelID != "model" || len(body.ContextualTuples.TupleKeys) != 1 {
			t.Errorf("unexpected body: %#v", body)
		}
		_, _ = w.Write([]byte(`{"allowed":true}`))
	}))
	defer server.Close()
	c, err := New(server.URL, "store", "model", "secret", time.Second)
	if err != nil {
		t.Fatal(err)
	}
	allowed, err := c.Check(context.Background(), authorization.Check{User: "user:u", Relation: "can_view", Object: "resource:r", ContextualTuples: []authorization.Tuple{{User: "user:u", Relation: "active_member", Object: "organization:o"}}})
	if err != nil || !allowed {
		t.Fatalf("allowed=%v err=%v", allowed, err)
	}
}

func TestCheckFailsClosedOnServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { http.Error(w, "no", http.StatusServiceUnavailable) }))
	defer server.Close()
	c, _ := New(server.URL, "store", "model", "secret", time.Second)
	if _, err := c.Check(context.Background(), authorization.Check{}); err == nil {
		t.Fatal("expected error")
	}
}

func TestCheckRejectsIncompleteResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(`{}`)) }))
	defer server.Close()
	c, _ := New(server.URL, "store", "model", "secret", time.Second)
	if _, err := c.Check(context.Background(), authorization.Check{}); err == nil {
		t.Fatal("expected incomplete response error")
	}
}
