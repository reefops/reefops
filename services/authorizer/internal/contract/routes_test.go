package contract

import "testing"

func TestResolveSyntheticResourceView(t *testing.T) {
	r, err := Resolve(SyntheticResourceView, "GET", "/_acceptance/authorization/resources/tank-1")
	if err != nil {
		t.Fatal(err)
	}
	if r.ResourceID != "tank-1" || r.Action != "resource:view" {
		t.Fatalf("unexpected route: %#v", r)
	}
}

func TestResolveRejectsUnlistedOrMalformedRequests(t *testing.T) {
	tests := [][3]string{{"unknown", "GET", "/_acceptance/authorization/resources/a"}, {SyntheticResourceView, "POST", "/_acceptance/authorization/resources/a"}, {SyntheticResourceView, "GET", "/_acceptance/authorization/resources/a/b"}, {SyntheticResourceView, "GET", "/_acceptance/authorization/resources/a?x=1"}}
	for _, tt := range tests {
		if _, err := Resolve(tt[0], tt[1], tt[2]); err == nil {
			t.Fatalf("expected rejection for %v", tt)
		}
	}
}
