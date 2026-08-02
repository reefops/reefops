package grpcauthz

import (
	"context"
	"log/slog"
	"strconv"
	"strings"
	"time"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	typev3 "github.com/envoyproxy/go-control-plane/envoy/type/v3"
	"github.com/google/uuid"
	"github.com/reefops/reefops/services/authorizer/internal/authorization"
	"github.com/reefops/reefops/services/authorizer/internal/observability"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

const authNamespace = "reefops.authentication"
const routeNamespace = "reefops.authorization"

type Authorizer interface {
	Authorize(context.Context, authorization.Input) authorization.Result
}
type Server struct {
	authv3.UnimplementedAuthorizationServer
	authorizer Authorizer
	metrics    *observability.Metrics
}

func New(authorizer Authorizer, metrics ...*observability.Metrics) *Server {
	s := &Server{authorizer: authorizer}
	if len(metrics) > 0 {
		s.metrics = metrics[0]
	}
	return s
}

func (s *Server) Check(ctx context.Context, req *authv3.CheckRequest) (*authv3.CheckResponse, error) {
	started := time.Now()
	in := extract(req)
	result := s.authorizer.Authorize(ctx, in)
	decisionResult := "deny"
	if result.Allowed {
		decisionResult = "allow"
	}
	if s.metrics != nil {
		s.metrics.Observe(decisionResult, result.ReasonCode, time.Since(started))
	}
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute.String("reefops.authorization.result", decisionResult), attribute.String("reefops.authorization.reason", result.ReasonCode))
	spanContext := span.SpanContext()
	slog.InfoContext(ctx, "authorization check completed", "result", decisionResult, "reason", result.ReasonCode, "decision_id", result.DecisionID, "correlation_id", result.CorrelationID, "trace_id", spanContext.TraceID().String(), "span_id", spanContext.SpanID().String(), "duration_ms", time.Since(started).Milliseconds())
	if result.Allowed {
		return &authv3.CheckResponse{Status: &status.Status{Code: int32(codes.OK)}, HttpResponse: &authv3.CheckResponse_OkResponse{OkResponse: &authv3.OkHttpResponse{Headers: []*corev3.HeaderValueOption{header("x-reefops-actor-context", result.ActorContext), header("x-reefops-authorization-decision-id", result.DecisionID), header("x-correlation-id", result.CorrelationID)}}}}, nil
	}
	return &authv3.CheckResponse{Status: &status.Status{Code: int32(codes.PermissionDenied), Message: result.ReasonCode}, HttpResponse: &authv3.CheckResponse_DeniedResponse{DeniedResponse: &authv3.DeniedHttpResponse{Status: &typev3.HttpStatus{Code: typev3.StatusCode_Forbidden}, Body: result.ReasonCode, Headers: []*corev3.HeaderValueOption{header("x-reefops-authorization-decision-id", result.DecisionID), header("x-correlation-id", result.CorrelationID)}}}}, nil
}
func header(k, v string) *corev3.HeaderValueOption {
	return &corev3.HeaderValueOption{Header: &corev3.HeaderValue{Key: k, Value: v}, AppendAction: corev3.HeaderValueOption_OVERWRITE_IF_EXISTS_OR_ADD}
}

func extract(req *authv3.CheckRequest) authorization.Input {
	in := authorization.Input{AttemptNumber: 1, CorrelationSource: "generated"}
	attrs := req.GetAttributes()
	httpReq := attrs.GetRequest().GetHttp()
	in.Method = httpReq.GetMethod()
	in.Path = httpReq.GetPath()
	in.RequestID, in.Invalid = requiredUUID(httpReq.GetId())
	in.AttemptID = in.RequestID
	if v := httpReq.GetHeaders()["x-attempt-id"]; v != "" {
		var bad bool
		in.AttemptID, bad = requiredUUID(v)
		in.Invalid = in.Invalid || bad
	}
	if v := httpReq.GetHeaders()["x-attempt-number"]; v != "" {
		n, err := strconv.ParseInt(v, 10, 32)
		if err != nil || n < 1 {
			in.Invalid = true
		} else {
			in.AttemptNumber = int32(n)
		}
	}
	in.CausationID = in.RequestID
	if v := httpReq.GetHeaders()["x-causation-id"]; v != "" {
		var bad bool
		in.CausationID, bad = requiredUUID(v)
		in.Invalid = in.Invalid || bad
	}
	if v := httpReq.GetHeaders()["x-correlation-id"]; v != "" {
		if id, err := uuid.Parse(v); err == nil {
			in.CorrelationID = id.String()
			in.CorrelationSource = "trusted"
		}
	}
	if in.CorrelationID == "" {
		in.CorrelationID = newUUID()
	}
	routeMetadata := attrs.GetRouteMetadataContext().GetFilterMetadata()
	identityMetadata := attrs.GetMetadataContext().GetFilterMetadata()
	if routeMetadata[authNamespace] != nil || identityMetadata[routeNamespace] != nil {
		in.Invalid = true
	}
	routeValues, bad := readFields(routeMetadata[routeNamespace], map[string]bool{"route_id": true})
	in.Invalid = in.Invalid || bad
	in.RouteID = routeValues["route_id"]
	identityValues, bad := readFields(identityMetadata[authNamespace], map[string]bool{"subject_id": true, "actor_id": true, "delegator_id": true, "credential_id_sha256": true, "issuer": true, "audience": true, "active_organization_id": true})
	in.IdentityMalformed = bad
	in.SubjectID = identityValues["subject_id"]
	in.ActorID = identityValues["actor_id"]
	in.DelegatorID = identityValues["delegator_id"]
	in.CredentialIDSHA256 = identityValues["credential_id_sha256"]
	in.Issuer = identityValues["issuer"]
	in.Audience = identityValues["audience"]
	in.OrganizationID = identityValues["active_organization_id"]
	return in
}
func requiredUUID(v string) (string, bool) {
	id, err := uuid.Parse(strings.TrimSpace(v))
	if err != nil {
		return newUUID(), true
	}
	return id.String(), false
}
func newUUID() string {
	id, err := uuid.NewV7()
	if err != nil {
		return uuid.NewString()
	}
	return id.String()
}
func readFields(s *structpb.Struct, allowed map[string]bool) (map[string]string, bool) {
	result := map[string]string{}
	if s == nil {
		return result, false
	}
	bad := false
	for key, value := range s.GetFields() {
		if !allowed[key] {
			bad = true
			continue
		}
		typed, ok := value.GetKind().(*structpb.Value_StringValue)
		if !ok {
			bad = true
			continue
		}
		result[key] = typed.StringValue
	}
	return result, bad
}
