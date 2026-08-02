package grpcauthz

import (
	"context"
	"testing"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	authv3 "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"github.com/google/uuid"
	"github.com/reefops/reefops/services/authorizer/internal/authorization"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/structpb"
)

type fakeAuthorizer struct {
	input  authorization.Input
	result authorization.Result
}

func (f *fakeAuthorizer) Authorize(_ context.Context, in authorization.Input) authorization.Result {
	f.input = in
	return f.result
}

func TestCheckExtractsTrustedMetadataAndReturnsAllowHeaders(t *testing.T) {
	requestID := uuid.NewString()
	organizationID := uuid.NewString()
	correlationID := uuid.NewString()
	fake := &fakeAuthorizer{result: authorization.Result{Allowed: true, ActorContext: "signed", DecisionID: uuid.NewString(), CorrelationID: correlationID}}
	server := New(fake)
	response, err := server.Check(context.Background(), &authv3.CheckRequest{Attributes: &authv3.AttributeContext{
		Request:              &authv3.AttributeContext_Request{Http: &authv3.AttributeContext_HttpRequest{Id: requestID, Method: "GET", Path: "/_acceptance/authorization/resources/r1", Headers: map[string]string{"x-correlation-id": correlationID}}},
		RouteMetadataContext: &corev3.Metadata{FilterMetadata: map[string]*structpb.Struct{routeNamespace: fields(map[string]string{"route_id": "acceptance.synthetic.resource.view"})}},
		MetadataContext:      &corev3.Metadata{FilterMetadata: map[string]*structpb.Struct{authNamespace: fields(map[string]string{"subject_id": "s", "actor_id": "a", "issuer": "i", "audience": "aud", "active_organization_id": organizationID})}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if codes.Code(response.GetStatus().GetCode()) != codes.OK || response.GetOkResponse() == nil {
		t.Fatalf("unexpected response: %#v", response)
	}
	if fake.input.Invalid || fake.input.RouteID == "" || fake.input.OrganizationID != organizationID {
		t.Fatalf("unexpected extracted input: %#v", fake.input)
	}
	if len(response.GetOkResponse().GetHeaders()) != 3 {
		t.Fatalf("unexpected headers: %#v", response.GetOkResponse().GetHeaders())
	}
}

func TestCheckNeverLeaksActorContextOnDeny(t *testing.T) {
	fake := &fakeAuthorizer{result: authorization.Result{ReasonCode: authorization.ReasonAuditFailure}}
	response, err := New(fake).Check(context.Background(), &authv3.CheckRequest{})
	if err != nil {
		t.Fatal(err)
	}
	if codes.Code(response.GetStatus().GetCode()) != codes.PermissionDenied || response.GetDeniedResponse() == nil || response.GetOkResponse() != nil {
		t.Fatalf("unexpected response: %#v", response)
	}
}
func TestExtractRejectsUnknownOrWrongTypedTrustedMetadata(t *testing.T) {
	requestID := uuid.NewString()
	route := fields(map[string]string{"route_id": "acceptance.synthetic.resource.view", "unknown": "value"})
	identity := fields(map[string]string{"subject_id": "s"})
	identity.Fields["actor_id"] = structpb.NewNumberValue(42)
	in := extract(&authv3.CheckRequest{Attributes: &authv3.AttributeContext{Request: &authv3.AttributeContext_Request{Http: &authv3.AttributeContext_HttpRequest{Id: requestID}}, RouteMetadataContext: &corev3.Metadata{FilterMetadata: map[string]*structpb.Struct{routeNamespace: route}}, MetadataContext: &corev3.Metadata{FilterMetadata: map[string]*structpb.Struct{authNamespace: identity}}}})
	if !in.Invalid || !in.IdentityMalformed {
		t.Fatalf("malformed metadata accepted: %#v", in)
	}
}
func fields(values map[string]string) *structpb.Struct {
	result := map[string]*structpb.Value{}
	for k, v := range values {
		result[k] = structpb.NewStringValue(v)
	}
	return &structpb.Struct{Fields: result}
}
