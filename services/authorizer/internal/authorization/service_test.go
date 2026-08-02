package authorization

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type fakeChecker struct {
	allowed bool
	err     error
}

func (f fakeChecker) Check(context.Context, Check) (bool, error) { return f.allowed, f.err }

type fakeSigner struct {
	signed SignedActorContext
	err    error
}

func (f fakeSigner) Sign(ActorContext) (SignedActorContext, error) { return f.signed, f.err }

type fakeAuditor struct {
	decisions []Decision
	err       error
}

func (f *fakeAuditor) Record(_ context.Context, d Decision) error {
	f.decisions = append(f.decisions, d)
	return f.err
}
func validTestInput() Input {
	return Input{RequestID: uuid.NewString(), AttemptID: uuid.NewString(), CorrelationID: uuid.NewString(), CorrelationSource: "trusted", CausationID: uuid.NewString(), AttemptNumber: 1, RouteID: "acceptance.synthetic.resource.view", Method: "GET", Path: "/_acceptance/authorization/resources/r1", SubjectID: "subject", ActorID: "actor", Issuer: "issuer", Audience: "audience", OrganizationID: uuid.NewString()}
}
func config() Config {
	return Config{EnvironmentID: "dev", StoreID: "store", ModelID: "model", ModelSHA256: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", ActorContextAudience: "reefops-services"}
}

func TestAllowIsReturnedOnlyAfterAudit(t *testing.T) {
	audit := &fakeAuditor{}
	now := time.Now()
	service, _ := NewService(config(), fakeChecker{allowed: true}, audit, fakeSigner{signed: SignedActorContext{Compact: "jws", JTI: uuid.NewString(), KID: "kid", Audience: "reefops-services", SHA256: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa", IssuedAt: now, ExpiresAt: now.Add(20 * time.Second)}})
	result := service.Authorize(context.Background(), validTestInput())
	if !result.Allowed || len(audit.decisions) != 1 || audit.decisions[0].Result != "allow" {
		t.Fatalf("result=%#v audit=%#v", result, audit.decisions)
	}
}
func TestAuditFailureChangesAllowToDeny(t *testing.T) {
	audit := &fakeAuditor{err: errors.New("database unavailable")}
	now := time.Now()
	service, _ := NewService(config(), fakeChecker{allowed: true}, audit, fakeSigner{signed: SignedActorContext{Compact: "jws", JTI: uuid.NewString(), IssuedAt: now, ExpiresAt: now.Add(time.Second)}})
	result := service.Authorize(context.Background(), validTestInput())
	if result.Allowed || result.ReasonCode != ReasonAuditFailure || result.ActorContext != "" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
func TestOpenFGAFailureIsAuditedAndDenied(t *testing.T) {
	audit := &fakeAuditor{}
	service, _ := NewService(config(), fakeChecker{err: errors.New("down")}, audit, fakeSigner{})
	result := service.Authorize(context.Background(), validTestInput())
	if result.Allowed || result.ReasonCode != ReasonOpenFGAUnavailable || len(audit.decisions) != 1 {
		t.Fatalf("result=%#v audit=%#v", result, audit.decisions)
	}
}

func TestMalformedIdentityIsAuditedBeforeOpenFGA(t *testing.T) {
	audit := &fakeAuditor{}
	service, _ := NewService(config(), fakeChecker{allowed: true}, audit, fakeSigner{})
	in := validTestInput()
	in.SubjectID = "bad:subject"
	result := service.Authorize(context.Background(), in)
	if result.ReasonCode != ReasonInvalidIdentity || len(audit.decisions) != 1 || audit.decisions[0].OpenFGAAttempts != 0 {
		t.Fatalf("result=%#v audit=%#v", result, audit.decisions)
	}
}
