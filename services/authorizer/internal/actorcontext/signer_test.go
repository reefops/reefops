package actorcontext

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"strings"
	"testing"
	"time"

	"github.com/reefops/reefops/services/authorizer/internal/authorization"
)

func TestSignProducesVerifiableCompactJWS(t *testing.T) {
	pub, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	der, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		t.Fatal(err)
	}
	publicDER, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	s, err := New(base64.StdEncoding.EncodeToString(der), base64.StdEncoding.EncodeToString(publicDER), "kid-1", "Ed25519", 20*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	s.now = func() time.Time { return time.Unix(1000, 0) }
	got, err := s.Sign(authorization.ActorContext{Audience: "reefops-services", EnvironmentID: "dev", ActorID: "a", SubjectID: "s", OrganizationID: "o", DecisionID: "d", CorrelationID: "c"})
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(got.Compact, ".")
	if len(parts) != 3 {
		t.Fatalf("not compact JWS: %q", got.Compact)
	}
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		t.Fatal(err)
	}
	if !ed25519.Verify(pub, []byte(parts[0]+"."+parts[1]), sig) {
		t.Fatal("signature is invalid")
	}
	if got.ExpiresAt.Sub(got.IssuedAt) != 20*time.Second {
		t.Fatalf("unexpected ttl: %s", got.ExpiresAt.Sub(got.IssuedAt))
	}
}
