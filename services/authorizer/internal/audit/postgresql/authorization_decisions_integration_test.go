package auditdb

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const auditTestDatabaseURLVariable = "REEFOPS_AUTHORIZER_AUDIT_TEST_DATABASE_URL"

func TestRecordAuthorizationDecisionIsIdempotentAndEnvironmentScoped(t *testing.T) {
	databaseURL := os.Getenv(auditTestDatabaseURLVariable)
	if databaseURL == "" {
		t.Skipf("%s is not set", auditTestDatabaseURLVariable)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conn, err := pgx.Connect(ctx, databaseURL)
	if err != nil {
		t.Fatalf("connect to audit database: %v", err)
	}
	defer func() {
		if err := conn.Close(context.Background()); err != nil {
			t.Errorf("close audit database: %v", err)
		}
	}()

	queries := New(conn)
	decidedAt := pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true}
	params := RecordAuthorizationDecisionParams{
		DecisionID:            testUUID(11),
		EnvironmentID:         "acceptance",
		RequestID:             testUUID(12),
		AttemptID:             testUUID(13),
		AttemptNumber:         1,
		DecisionPayloadSha256: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		CorrelationID:         testUUID(14),
		CorrelationSource:     "generated",
		CausationID:           testUUID(15),
		DecisionStage:         "identity",
		Result:                "deny",
		ReasonCode:            "missing_identity",
		DecidedAt:             decidedAt,
	}

	first, err := queries.RecordAuthorizationDecision(ctx, params)
	if err != nil {
		t.Fatalf("record first decision: %v", err)
	}
	second, err := queries.RecordAuthorizationDecision(ctx, params)
	if err != nil {
		t.Fatalf("repeat identical decision: %v", err)
	}
	if first.DecisionPayloadSha256 != second.DecisionPayloadSha256 {
		t.Fatal("idempotent replay returned different evidence")
	}

	conflict := params
	conflict.DecisionPayloadSha256 = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	if _, err := queries.RecordAuthorizationDecision(ctx, conflict); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("conflicting replay error = %v, want pgx.ErrNoRows", err)
	}

	if _, err := queries.GetAuthorizationDecision(ctx, GetAuthorizationDecisionParams{
		EnvironmentID: "another-environment",
		DecisionID:    params.DecisionID,
	}); !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("cross-environment read error = %v, want pgx.ErrNoRows", err)
	}

	decisions, err := queries.ListAuthorizationDecisionsByCorrelation(
		ctx,
		ListAuthorizationDecisionsByCorrelationParams{
			EnvironmentID: params.EnvironmentID,
			CorrelationID: params.CorrelationID,
		},
	)
	if err != nil {
		t.Fatalf("list decisions by correlation: %v", err)
	}
	if len(decisions) != 1 {
		t.Fatalf("decision count = %d, want 1", len(decisions))
	}
}

func testUUID(last byte) pgtype.UUID {
	var value [16]byte
	value[0] = 1
	value[6] = 0x70
	value[8] = 0x80
	value[15] = last
	return pgtype.UUID{Bytes: value, Valid: true}
}
