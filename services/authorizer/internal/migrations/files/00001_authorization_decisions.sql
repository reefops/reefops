-- +goose Up
-- +goose StatementBegin
SET LOCAL ROLE reefops_authorizer_owner;

CREATE SCHEMA authorizer_audit AUTHORIZATION reefops_authorizer_owner;

REVOKE ALL ON SCHEMA authorizer_audit FROM PUBLIC;

CREATE TABLE authorizer_audit.authorization_decisions (
    decision_id uuid PRIMARY KEY,
    environment_id text NOT NULL CHECK (environment_id ~ '^[a-z][a-z0-9-]{0,62}$'),
    request_id uuid NOT NULL,
    attempt_id uuid NOT NULL UNIQUE,
    attempt_number integer NOT NULL CHECK (attempt_number > 0),
    decision_payload_sha256 text NOT NULL CHECK (decision_payload_sha256 ~ '^sha256:[0-9a-f]{64}$'),
    correlation_id uuid NOT NULL,
    correlation_source text NOT NULL CHECK (
        correlation_source IN ('trusted', 'generated')
    ),
    causation_id uuid NOT NULL,
    decision_stage text NOT NULL CHECK (
        decision_stage IN ('route_resolution', 'identity', 'tenant', 'openfga', 'actor_context')
    ),
    route_id text CHECK (route_id ~ '^[a-z][a-z0-9._-]{0,127}$'),
    route_contract_version text CHECK (length(route_contract_version) BETWEEN 1 AND 64),
    http_method text CHECK (
        http_method IN ('GET', 'HEAD', 'POST', 'PUT', 'PATCH', 'DELETE', 'OPTIONS')
    ),
    route_template text CHECK (
        length(route_template) BETWEEN 1 AND 512 AND left(route_template, 1) = '/'
    ),
    principal_type text CHECK (principal_type IN ('human', 'service_account', 'share_link')),
    subject_id text CHECK (length(subject_id) BETWEEN 1 AND 512),
    actor_id text CHECK (length(actor_id) BETWEEN 1 AND 512),
    delegator_id text CHECK (length(delegator_id) BETWEEN 1 AND 512),
    credential_id_sha256 text CHECK (credential_id_sha256 ~ '^sha256:[0-9a-f]{64}$'),
    oidc_issuer text CHECK (length(oidc_issuer) BETWEEN 1 AND 512),
    oidc_audience text CHECK (length(oidc_audience) BETWEEN 1 AND 256),
    active_organization_id uuid,
    action text CHECK (action ~ '^[a-z][a-z0-9_]{0,62}:[a-z][a-z0-9_]{0,62}$'),
    resource_type text CHECK (resource_type ~ '^[a-z][a-z0-9_]{0,63}$'),
    resource_id text CHECK (length(resource_id) BETWEEN 1 AND 512),
    openfga_store_id text CHECK (length(openfga_store_id) BETWEEN 1 AND 128),
    openfga_model_id text CHECK (length(openfga_model_id) BETWEEN 1 AND 128),
    openfga_model_sha256 text CHECK (openfga_model_sha256 ~ '^sha256:[0-9a-f]{64}$'),
    openfga_user text CHECK (length(openfga_user) BETWEEN 1 AND 512),
    openfga_relation text CHECK (openfga_relation ~ '^[a-z][a-z0-9_]{0,63}$'),
    openfga_object text CHECK (length(openfga_object) BETWEEN 1 AND 512),
    contextual_tuples_sha256 text CHECK (contextual_tuples_sha256 ~ '^sha256:[0-9a-f]{64}$'),
    contextual_tuples_canonical_version text CHECK (
        length(contextual_tuples_canonical_version) BETWEEN 1 AND 32
    ),
    result text NOT NULL CHECK (result IN ('allow', 'deny')),
    reason_code text NOT NULL CHECK (reason_code ~ '^[a-z][a-z0-9_]{0,127}$'),
    openfga_latency_ms integer CHECK (openfga_latency_ms >= 0),
    openfga_attempts integer CHECK (openfga_attempts >= 0),
    actor_context_jti uuid,
    actor_context_kid text CHECK (length(actor_context_kid) BETWEEN 1 AND 128),
    actor_context_audience text CHECK (length(actor_context_audience) BETWEEN 1 AND 256),
    actor_context_sha256 text CHECK (actor_context_sha256 ~ '^sha256:[0-9a-f]{64}$'),
    actor_context_issued_at timestamptz,
    actor_context_expires_at timestamptz,
    decided_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    UNIQUE (environment_id, request_id, attempt_number),
    CHECK (delegator_id IS NULL OR delegator_id <> actor_id),
    CHECK (
        (contextual_tuples_sha256 IS NULL) = (contextual_tuples_canonical_version IS NULL)
    ),
    CHECK (
        (openfga_latency_ms IS NULL) = (openfga_attempts IS NULL)
    ),
    CHECK (
        result = 'deny'
        OR (
            decision_stage = 'actor_context'
            AND route_id IS NOT NULL
            AND route_contract_version IS NOT NULL
            AND http_method IS NOT NULL
            AND route_template IS NOT NULL
            AND principal_type IS NOT NULL
            AND subject_id IS NOT NULL
            AND actor_id IS NOT NULL
            AND active_organization_id IS NOT NULL
            AND action IS NOT NULL
            AND resource_type IS NOT NULL
            AND resource_id IS NOT NULL
            AND openfga_store_id IS NOT NULL
            AND openfga_model_id IS NOT NULL
            AND openfga_model_sha256 IS NOT NULL
            AND openfga_user IS NOT NULL
            AND openfga_relation IS NOT NULL
            AND openfga_object IS NOT NULL
            AND contextual_tuples_sha256 IS NOT NULL
            AND openfga_latency_ms IS NOT NULL
            AND openfga_attempts IS NOT NULL
            AND actor_context_jti IS NOT NULL
            AND actor_context_kid IS NOT NULL
            AND actor_context_audience IS NOT NULL
            AND actor_context_sha256 IS NOT NULL
            AND actor_context_issued_at IS NOT NULL
            AND actor_context_expires_at IS NOT NULL
        )
    ),
    CHECK (
        result = 'allow'
        OR (
            actor_context_jti IS NULL
            AND actor_context_kid IS NULL
            AND actor_context_audience IS NULL
            AND actor_context_sha256 IS NULL
            AND actor_context_issued_at IS NULL
            AND actor_context_expires_at IS NULL
        )
    ),
    CHECK (
        actor_context_expires_at IS NULL
        OR (
            actor_context_issued_at <= decided_at
            AND actor_context_expires_at > actor_context_issued_at
            AND actor_context_expires_at <= actor_context_issued_at + interval '30 seconds'
        )
    )
);

CREATE INDEX authorization_decisions_correlation_idx
    ON authorizer_audit.authorization_decisions (environment_id, correlation_id, decided_at);

CREATE INDEX authorization_decisions_subject_idx
    ON authorizer_audit.authorization_decisions (environment_id, subject_id, decided_at)
    WHERE subject_id IS NOT NULL;

CREATE FUNCTION authorizer_audit.reject_authorization_decision_mutation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'authorization decisions are append-only'
        USING ERRCODE = '55000';
END;
$$;

CREATE TRIGGER authorization_decisions_reject_mutation
BEFORE UPDATE OR DELETE ON authorizer_audit.authorization_decisions
FOR EACH ROW EXECUTE FUNCTION authorizer_audit.reject_authorization_decision_mutation();

CREATE TRIGGER authorization_decisions_reject_truncate
BEFORE TRUNCATE ON authorizer_audit.authorization_decisions
FOR EACH STATEMENT EXECUTE FUNCTION authorizer_audit.reject_authorization_decision_mutation();

REVOKE ALL ON ALL TABLES IN SCHEMA authorizer_audit FROM PUBLIC;
REVOKE ALL ON ALL FUNCTIONS IN SCHEMA authorizer_audit FROM PUBLIC;

GRANT USAGE ON SCHEMA authorizer_audit TO reefops_authorizer_runtime;
GRANT SELECT, INSERT ON authorizer_audit.authorization_decisions TO reefops_authorizer_runtime;

RESET ROLE;
-- +goose StatementEnd
