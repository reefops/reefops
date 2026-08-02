package authorization

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reefops/reefops/services/authorizer/internal/contract"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type Config struct{ EnvironmentID, StoreID, ModelID, ModelSHA256, ActorContextAudience string }
type Service struct {
	config  Config
	checker Checker
	auditor Auditor
	signer  Signer
	now     func() time.Time
}

func NewService(config Config, checker Checker, auditor Auditor, signer Signer) (*Service, error) {
	if config.EnvironmentID == "" || config.StoreID == "" || config.ModelID == "" || config.ModelSHA256 == "" || config.ActorContextAudience == "" {
		return nil, fmt.Errorf("authorizer configuration is incomplete")
	}
	if checker == nil || auditor == nil || signer == nil {
		return nil, fmt.Errorf("authorizer dependencies are required")
	}
	return &Service{config: config, checker: checker, auditor: auditor, signer: signer, now: time.Now}, nil
}

func (s *Service) Authorize(ctx context.Context, in Input) Result {
	decisionID, err := uuid.NewV7()
	if err != nil {
		return Result{ReasonCode: ReasonInvalidRequest}
	}
	d := Decision{ID: decisionID.String(), EnvironmentID: s.config.EnvironmentID, Input: in, Result: "deny", ReasonCode: ReasonInvalidRequest, Stage: "route_resolution", OpenFGAStoreID: s.config.StoreID, OpenFGAModelID: s.config.ModelID, OpenFGAModelSHA256: s.config.ModelSHA256, DecidedAt: s.now().UTC()}
	deny := func(reason, stage string) Result {
		d.ReasonCode = reason
		d.Stage = stage
		d.DecidedAt = s.now().UTC()
		if err := s.record(ctx, d); err != nil {
			return Result{ReasonCode: ReasonAuditFailure, DecisionID: d.ID, CorrelationID: in.CorrelationID}
		}
		return Result{ReasonCode: reason, DecisionID: d.ID, CorrelationID: in.CorrelationID}
	}
	if !validInput(in) {
		return deny(ReasonInvalidRequest, "identity")
	}
	if missingIdentity(in) {
		return deny(ReasonMissingIdentity, "identity")
	}
	if in.IdentityMalformed || !validIdentity(in) {
		return deny(ReasonInvalidIdentity, "identity")
	}
	route, err := contract.Resolve(in.RouteID, in.Method, in.Path)
	if err != nil {
		reason := ReasonMalformedPath
		if errors.Is(err, contract.ErrUnknownRoute) {
			reason = ReasonUnknownRoute
		} else if errors.Is(err, contract.ErrMethodMismatch) {
			reason = ReasonMethodMismatch
		}
		return deny(reason, "route_resolution")
	}
	d.Route = Route{ID: route.ID, ContractVersion: route.ContractVersion, Method: route.Method, Template: route.Template, Action: route.Action, ResourceType: route.ResourceType, ResourceID: route.ResourceID}
	d.Stage = "openfga"
	d.OpenFGAUser = "user:" + in.SubjectID
	d.OpenFGARelation = "can_view"
	d.OpenFGAObject = "resource:" + route.ResourceID
	tuple := Tuple{User: d.OpenFGAUser, Relation: "active_member", Object: "organization:" + in.OrganizationID}
	d.ContextualTuplesSHA256 = tupleDigest([]Tuple{tuple})
	started := s.now()
	checkContext, checkSpan := otel.Tracer("reefops-authorizer").Start(ctx, "openfga.check")
	allowed, err := s.checker.Check(checkContext, Check{User: d.OpenFGAUser, Relation: d.OpenFGARelation, Object: d.OpenFGAObject, ContextualTuples: []Tuple{tuple}})
	checkSpan.SetAttributes(attribute.Bool("reefops.openfga.allowed", allowed))
	if err != nil {
		checkSpan.SetStatus(codes.Error, "unavailable")
	}
	checkSpan.End()
	d.OpenFGALatency = s.now().Sub(started)
	d.OpenFGAAttempts = 1
	if err != nil {
		return deny(ReasonOpenFGAUnavailable, "openfga")
	}
	if !allowed {
		return deny(ReasonOpenFGADenied, "openfga")
	}
	d.Stage = "actor_context"
	_, signSpan := otel.Tracer("reefops-authorizer").Start(ctx, "actor_context.sign")
	signed, err := s.signer.Sign(ActorContext{Audience: s.config.ActorContextAudience, EnvironmentID: s.config.EnvironmentID, ActorID: in.ActorID, SubjectID: in.SubjectID, DelegatorID: in.DelegatorID, OrganizationID: in.OrganizationID, Action: route.Action, ResourceType: route.ResourceType, ResourceID: route.ResourceID, DecisionID: d.ID, CorrelationID: in.CorrelationID, OpenFGAStoreID: s.config.StoreID, OpenFGAModelID: s.config.ModelID})
	if err != nil {
		signSpan.SetStatus(codes.Error, "signing failed")
	}
	signSpan.End()
	if err != nil {
		return deny(ReasonActorContextFailure, "actor_context")
	}
	d.Result = "allow"
	d.ReasonCode = ReasonAllowed
	d.ActorContext = signed.Compact
	d.ActorContextJTI = signed.JTI
	d.ActorContextKID = signed.KID
	d.ActorContextAudience = signed.Audience
	d.ActorContextSHA256 = signed.SHA256
	d.ActorContextIssuedAt = signed.IssuedAt
	d.ActorContextExpiresAt = signed.ExpiresAt
	d.DecidedAt = s.now().UTC()
	if err := s.record(ctx, d); err != nil {
		return Result{ReasonCode: ReasonAuditFailure, DecisionID: d.ID, CorrelationID: in.CorrelationID}
	}
	return Result{Allowed: true, ReasonCode: ReasonAllowed, DecisionID: d.ID, CorrelationID: in.CorrelationID, ActorContext: signed.Compact}
}

func (s *Service) record(requestContext context.Context, decision Decision) error {
	auditContext, cancel := context.WithTimeout(context.WithoutCancel(requestContext), 2*time.Second)
	defer cancel()
	auditContext, span := otel.Tracer("reefops-authorizer").Start(auditContext, "audit.persist", trace.WithAttributes(attribute.String("reefops.authorization.result", decision.Result)))
	defer span.End()
	err := s.auditor.Record(auditContext, decision)
	if err != nil {
		span.SetStatus(codes.Error, "persistence failed")
	}
	return err
}

func validInput(in Input) bool {
	if in.Invalid {
		return false
	}
	if _, err := uuid.Parse(in.RequestID); err != nil {
		return false
	}
	if _, err := uuid.Parse(in.AttemptID); err != nil {
		return false
	}
	if _, err := uuid.Parse(in.CorrelationID); err != nil {
		return false
	}
	if _, err := uuid.Parse(in.CausationID); err != nil {
		return false
	}
	return in.AttemptNumber > 0
}

func missingIdentity(in Input) bool {
	return in.SubjectID == "" || in.ActorID == "" || in.Issuer == "" || in.Audience == "" || in.OrganizationID == ""
}
func validIdentity(in Input) bool {
	if _, err := uuid.Parse(in.OrganizationID); err != nil {
		return false
	}
	return validIdentifier(in.SubjectID) && validIdentifier(in.ActorID) && validOptionalIdentifier(in.DelegatorID) && in.DelegatorID != in.ActorID && validText(in.Issuer, 512) && validText(in.Audience, 256) && validOptionalDigest(in.CredentialIDSHA256)
}
func validIdentifier(value string) bool {
	if value == "" || len(value) > 256 {
		return false
	}
	for _, c := range value {
		if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '.' || c == '_' || c == '-' || c == '@') {
			return false
		}
	}
	return true
}
func validOptionalIdentifier(value string) bool { return value == "" || validIdentifier(value) }
func validText(value string, max int) bool {
	if value == "" || len(value) > max {
		return false
	}
	for _, c := range value {
		if c < 0x20 || c == 0x7f {
			return false
		}
	}
	return true
}
func validOptionalDigest(value string) bool {
	if value == "" {
		return true
	}
	if len(value) != 71 || value[:7] != "sha256:" {
		return false
	}
	for _, c := range value[7:] {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}

func tupleDigest(tuples []Tuple) string {
	b, _ := json.Marshal(tuples)
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:])
}
