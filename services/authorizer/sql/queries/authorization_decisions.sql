-- name: RecordAuthorizationDecision :one
WITH inserted AS (
    INSERT INTO authorizer_audit.authorization_decisions (
        decision_id, environment_id, request_id, attempt_id, attempt_number,
        decision_payload_sha256, correlation_id, correlation_source, causation_id,
        decision_stage, route_id, route_contract_version, http_method, route_template,
        principal_type, subject_id, actor_id, delegator_id, credential_id_sha256,
        oidc_issuer, oidc_audience, active_organization_id, action, resource_type,
        resource_id, openfga_store_id, openfga_model_id, openfga_model_sha256,
        openfga_user, openfga_relation, openfga_object, contextual_tuples_sha256,
        contextual_tuples_canonical_version, result, reason_code, openfga_latency_ms,
        openfga_attempts, actor_context_jti, actor_context_kid,
        actor_context_audience, actor_context_sha256, actor_context_issued_at,
        actor_context_expires_at, decided_at
    ) VALUES (
        sqlc.arg(decision_id), sqlc.arg(environment_id), sqlc.arg(request_id),
        sqlc.arg(attempt_id), sqlc.arg(attempt_number), sqlc.arg(decision_payload_sha256),
        sqlc.arg(correlation_id), sqlc.arg(correlation_source), sqlc.arg(causation_id),
        sqlc.arg(decision_stage), sqlc.narg(route_id), sqlc.narg(route_contract_version),
        sqlc.narg(http_method), sqlc.narg(route_template), sqlc.narg(principal_type),
        sqlc.narg(subject_id), sqlc.narg(actor_id), sqlc.narg(delegator_id),
        sqlc.narg(credential_id_sha256), sqlc.narg(oidc_issuer), sqlc.narg(oidc_audience),
        sqlc.narg(active_organization_id), sqlc.narg(action), sqlc.narg(resource_type),
        sqlc.narg(resource_id), sqlc.narg(openfga_store_id), sqlc.narg(openfga_model_id),
        sqlc.narg(openfga_model_sha256), sqlc.narg(openfga_user),
        sqlc.narg(openfga_relation), sqlc.narg(openfga_object),
        sqlc.narg(contextual_tuples_sha256), sqlc.narg(contextual_tuples_canonical_version),
        sqlc.arg(result), sqlc.arg(reason_code), sqlc.narg(openfga_latency_ms),
        sqlc.narg(openfga_attempts), sqlc.narg(actor_context_jti),
        sqlc.narg(actor_context_kid), sqlc.narg(actor_context_audience),
        sqlc.narg(actor_context_sha256), sqlc.narg(actor_context_issued_at),
        sqlc.narg(actor_context_expires_at), sqlc.arg(decided_at)
    )
    ON CONFLICT DO NOTHING
    RETURNING *
)
SELECT * FROM inserted
UNION ALL
SELECT *
FROM authorizer_audit.authorization_decisions
WHERE environment_id = sqlc.arg(environment_id)
  AND decision_id = sqlc.arg(decision_id)
  AND decision_payload_sha256 = sqlc.arg(decision_payload_sha256)
LIMIT 1;

-- name: GetAuthorizationDecision :one
SELECT *
FROM authorizer_audit.authorization_decisions
WHERE environment_id = sqlc.arg(environment_id)
  AND decision_id = sqlc.arg(decision_id);

-- name: ListAuthorizationDecisionsByCorrelation :many
SELECT *
FROM authorizer_audit.authorization_decisions
WHERE environment_id = sqlc.arg(environment_id)
  AND correlation_id = sqlc.arg(correlation_id)
ORDER BY decided_at, decision_id;
