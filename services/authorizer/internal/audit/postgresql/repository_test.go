package auditdb

import (
	"github.com/reefops/reefops/services/authorizer/internal/authorization"
	"testing"
)

func TestDecisionPayloadDigestCoversTraceabilityEvidenceButNotRawJWS(t *testing.T) {
	base := authorization.Decision{ID: "d", Input: authorization.Input{CorrelationID: "c"}, ActorContext: "secret-jws", ActorContextSHA256: "sha256:a"}
	first, err := decisionPayloadDigest(base)
	if err != nil {
		t.Fatal(err)
	}
	changed := base
	changed.Input.CorrelationID = "different"
	second, _ := decisionPayloadDigest(changed)
	if first == second {
		t.Fatal("correlation evidence did not affect digest")
	}
	changed = base
	changed.ActorContext = "another-secret-jws"
	third, _ := decisionPayloadDigest(changed)
	if first != third {
		t.Fatal("raw JWS must not affect canonical decision digest")
	}
}
