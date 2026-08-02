package auditdb

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/reefops/reefops/services/authorizer/internal/authorization"
)

type Repository struct{ queries *Queries }

func NewRepository(db DBTX) *Repository { return &Repository{queries: New(db)} }

func (r *Repository) Record(ctx context.Context, d authorization.Decision) error {
	payloadDigest, err := decisionPayloadDigest(d)
	if err != nil {
		return err
	}
	params := RecordAuthorizationDecisionParams{
		DecisionID: uuidValue(d.ID), EnvironmentID: d.EnvironmentID, RequestID: uuidValue(d.Input.RequestID), AttemptID: uuidValue(d.Input.AttemptID), AttemptNumber: d.Input.AttemptNumber,
		DecisionPayloadSha256: payloadDigest, CorrelationID: uuidValue(d.Input.CorrelationID), CorrelationSource: d.Input.CorrelationSource, CausationID: uuidValue(d.Input.CausationID), DecisionStage: d.Stage,
		RouteID: strptr(d.Route.ID), RouteContractVersion: strptr(d.Route.ContractVersion), HttpMethod: strptr(d.Route.Method), RouteTemplate: strptr(d.Route.Template), PrincipalType: strptrIf(d.Input.SubjectID, "human"), SubjectID: bounded(d.Input.SubjectID, 512), ActorID: bounded(d.Input.ActorID, 512), DelegatorID: delegator(d.Input.DelegatorID, d.Input.ActorID), CredentialIDSha256: digest(d.Input.CredentialIDSHA256), OidcIssuer: bounded(d.Input.Issuer, 512), OidcAudience: bounded(d.Input.Audience, 256), ActiveOrganizationID: uuidValue(d.Input.OrganizationID), Action: strptr(d.Route.Action), ResourceType: strptr(d.Route.ResourceType), ResourceID: strptr(d.Route.ResourceID),
		OpenfgaStoreID: strptr(d.OpenFGAStoreID), OpenfgaModelID: strptr(d.OpenFGAModelID), OpenfgaModelSha256: strptr(d.OpenFGAModelSHA256), OpenfgaUser: strptr(d.OpenFGAUser), OpenfgaRelation: strptr(d.OpenFGARelation), OpenfgaObject: strptr(d.OpenFGAObject), ContextualTuplesSha256: strptr(d.ContextualTuplesSHA256), ContextualTuplesCanonicalVersion: strptrIf(d.ContextualTuplesSHA256, "v1"), Result: d.Result, ReasonCode: d.ReasonCode,
		OpenfgaLatencyMs: durationMillis(d.OpenFGALatency, d.OpenFGAAttempts), OpenfgaAttempts: intptr(d.OpenFGAAttempts), ActorContextJti: uuidValue(d.ActorContextJTI), ActorContextKid: strptr(d.ActorContextKID), ActorContextAudience: strptr(d.ActorContextAudience), ActorContextSha256: strptr(d.ActorContextSHA256), ActorContextIssuedAt: timeValue(d.ActorContextIssuedAt), ActorContextExpiresAt: timeValue(d.ActorContextExpiresAt), DecidedAt: timeValue(d.DecidedAt),
	}
	if _, err = r.queries.RecordAuthorizationDecision(ctx, params); err != nil {
		return fmt.Errorf("record authorization decision: %w", err)
	}
	return nil
}

func decisionPayloadDigest(d authorization.Decision) (string, error) {
	payload := d
	payload.ActorContext = ""
	b, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(append([]byte("reefops.authorization-decision/v1\n"), b...))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func uuidValue(s string) pgtype.UUID {
	id, err := uuid.Parse(s)
	if err != nil {
		return pgtype.UUID{}
	}
	var b [16]byte
	copy(b[:], id[:])
	return pgtype.UUID{Bytes: b, Valid: true}
}
func strptr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func strptrIf(condition, value string) *string {
	if condition == "" {
		return nil
	}
	return &value
}
func bounded(value string, max int) *string {
	if value == "" || len(value) > max {
		return nil
	}
	return &value
}
func delegator(value, actor string) *string {
	if value == actor {
		return nil
	}
	return bounded(value, 512)
}
func digest(value string) *string {
	if len(value) != 71 || value[:7] != "sha256:" {
		return nil
	}
	for _, c := range value[7:] {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return nil
		}
	}
	return &value
}
func intptr(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}
func durationMillis(v time.Duration, attempts int32) *int32 {
	if attempts == 0 {
		return nil
	}
	n := int32(v.Milliseconds())
	return &n
}
func timeValue(v time.Time) pgtype.Timestamptz {
	return pgtype.Timestamptz{Time: v, Valid: !v.IsZero()}
}
