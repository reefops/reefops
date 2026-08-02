package authorization

import (
	"context"
	"time"
)

const (
	ReasonAllowed             = "allowed"
	ReasonUnknownRoute        = "unknown_route"
	ReasonMethodMismatch      = "method_mismatch"
	ReasonMalformedPath       = "malformed_path"
	ReasonMissingIdentity     = "missing_identity"
	ReasonInvalidIdentity     = "invalid_identity"
	ReasonInvalidRequest      = "invalid_request_context"
	ReasonOpenFGADenied       = "openfga_denied"
	ReasonOpenFGAUnavailable  = "openfga_unavailable"
	ReasonActorContextFailure = "signing_failed"
	ReasonAuditFailure        = "audit_unavailable"
)

type Input struct {
	Invalid, IdentityMalformed                           bool
	RequestID, AttemptID, CorrelationID, CausationID     string
	CorrelationSource                                    string
	AttemptNumber                                        int32
	RouteID, Method, Path                                string
	SubjectID, ActorID, DelegatorID                      string
	CredentialIDSHA256, Issuer, Audience, OrganizationID string
}

type Route struct {
	ID, ContractVersion, Method, Template string
	Action, ResourceType, ResourceID      string
}

type Decision struct {
	ID, EnvironmentID                                                                        string
	Input                                                                                    Input
	Route                                                                                    Route
	Result, ReasonCode, Stage                                                                string
	OpenFGAUser, OpenFGARelation, OpenFGAObject                                              string
	OpenFGAStoreID, OpenFGAModelID, OpenFGAModelSHA256                                       string
	ContextualTuplesSHA256                                                                   string
	OpenFGALatency                                                                           time.Duration
	OpenFGAAttempts                                                                          int32
	ActorContext, ActorContextJTI, ActorContextKID, ActorContextAudience, ActorContextSHA256 string
	ActorContextIssuedAt, ActorContextExpiresAt                                              time.Time
	DecidedAt                                                                                time.Time
}

type Check struct {
	User, Relation, Object string
	ContextualTuples       []Tuple
}

type Tuple struct{ User, Relation, Object string }

type Checker interface {
	Check(context.Context, Check) (bool, error)
}
type Auditor interface {
	Record(context.Context, Decision) error
}
type Signer interface {
	Sign(ActorContext) (SignedActorContext, error)
}

type ActorContext struct {
	ID, Audience, EnvironmentID                                 string
	ActorID, SubjectID, DelegatorID, OrganizationID             string
	Action, ResourceType, ResourceID, DecisionID, CorrelationID string
	OpenFGAStoreID, OpenFGAModelID                              string
}

type SignedActorContext struct {
	Compact, JTI, KID, Audience, SHA256 string
	IssuedAt, ExpiresAt                 time.Time
}

type Result struct {
	Allowed                                             bool
	ReasonCode, DecisionID, CorrelationID, ActorContext string
}
