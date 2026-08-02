package actorcontext

import (
	"crypto/ed25519"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reefops/reefops/services/authorizer/internal/authorization"
)

type Signer struct {
	key ed25519.PrivateKey
	kid string
	ttl time.Duration
	now func() time.Time
}

func New(privateKeyPKCS8PEMB64, publicKeyPEMB64, kid, algorithm string, ttl time.Duration) (*Signer, error) {
	if algorithm != "Ed25519" && algorithm != "EdDSA" {
		return nil, fmt.Errorf("unsupported actor context algorithm %q", algorithm)
	}
	if kid == "" {
		return nil, errors.New("actor context kid is required")
	}
	if ttl <= 0 || ttl > 30*time.Second {
		return nil, errors.New("actor context ttl must be between 1ns and 30s")
	}
	raw, err := base64.StdEncoding.DecodeString(privateKeyPKCS8PEMB64)
	if err != nil {
		return nil, fmt.Errorf("decode private key: %w", err)
	}
	if block, _ := pem.Decode(raw); block != nil {
		raw = block.Bytes
	}
	parsed, err := x509.ParsePKCS8PrivateKey(raw)
	if err != nil {
		return nil, fmt.Errorf("parse private key: %w", err)
	}
	key, ok := parsed.(ed25519.PrivateKey)
	if !ok {
		return nil, errors.New("private key is not Ed25519")
	}
	publicRaw, err := base64.StdEncoding.DecodeString(publicKeyPEMB64)
	if err != nil {
		return nil, fmt.Errorf("decode public key: %w", err)
	}
	if block, _ := pem.Decode(publicRaw); block != nil {
		publicRaw = block.Bytes
	}
	publicParsed, err := x509.ParsePKIXPublicKey(publicRaw)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}
	publicKey, ok := publicParsed.(ed25519.PublicKey)
	if !ok || !key.Public().(ed25519.PublicKey).Equal(publicKey) {
		return nil, errors.New("public key does not match Ed25519 private key")
	}
	return &Signer{key: key, kid: kid, ttl: ttl, now: time.Now}, nil
}

type header struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}
type claims struct {
	JTI            string `json:"jti"`
	Issuer         string `json:"iss"`
	Audience       string `json:"aud"`
	IssuedAt       int64  `json:"iat"`
	NotBefore      int64  `json:"nbf"`
	ExpiresAt      int64  `json:"exp"`
	EnvironmentID  string `json:"environment_id"`
	ActorID        string `json:"actor_id"`
	SubjectID      string `json:"subject_id"`
	DelegatorID    string `json:"delegator_id,omitempty"`
	OrganizationID string `json:"organization_id"`
	Action         string `json:"action"`
	ResourceType   string `json:"resource_type"`
	ResourceID     string `json:"resource_id"`
	DecisionID     string `json:"decision_id"`
	CorrelationID  string `json:"correlation_id"`
	OpenFGAStoreID string `json:"openfga_store_id"`
	OpenFGAModelID string `json:"openfga_model_id"`
}

func (s *Signer) Sign(in authorization.ActorContext) (authorization.SignedActorContext, error) {
	now := s.now().UTC().Truncate(time.Second)
	expires := now.Add(s.ttl)
	jti, err := uuid.NewV7()
	if err != nil {
		return authorization.SignedActorContext{}, err
	}
	h := header{Alg: "EdDSA", Kid: s.kid, Typ: "reefops-actor-context+jwt"}
	c := claims{JTI: jti.String(), Issuer: "reefops-authorizer", Audience: in.Audience, IssuedAt: now.Unix(), NotBefore: now.Unix(), ExpiresAt: expires.Unix(), EnvironmentID: in.EnvironmentID, ActorID: in.ActorID, SubjectID: in.SubjectID, DelegatorID: in.DelegatorID, OrganizationID: in.OrganizationID, Action: in.Action, ResourceType: in.ResourceType, ResourceID: in.ResourceID, DecisionID: in.DecisionID, CorrelationID: in.CorrelationID, OpenFGAStoreID: in.OpenFGAStoreID, OpenFGAModelID: in.OpenFGAModelID}
	headerJSON, err := json.Marshal(h)
	if err != nil {
		return authorization.SignedActorContext{}, err
	}
	claimsJSON, err := json.Marshal(c)
	if err != nil {
		return authorization.SignedActorContext{}, err
	}
	enc := base64.RawURLEncoding
	input := enc.EncodeToString(headerJSON) + "." + enc.EncodeToString(claimsJSON)
	compact := input + "." + enc.EncodeToString(ed25519.Sign(s.key, []byte(input)))
	digest := sha256.Sum256([]byte(compact))
	return authorization.SignedActorContext{Compact: compact, JTI: jti.String(), KID: s.kid, Audience: in.Audience, SHA256: "sha256:" + hex.EncodeToString(digest[:]), IssuedAt: now, ExpiresAt: expires}, nil
}
