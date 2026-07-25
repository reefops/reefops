#!/usr/bin/env bash
set -euo pipefail

project_root="$(git rev-parse --show-toplevel)"
config_file="${REEFOPS_CI_SECRETS_CONFIG:-${project_root}/.github/ci-secrets.yaml}"
secret_name="${REEFOPS_CI_SECRET_NAME:?Define REEFOPS_CI_SECRET_NAME}"
requested_version="${REEFOPS_CI_SECRET_VERSION:?Define REEFOPS_CI_SECRET_VERSION}"
requested_action="${REEFOPS_CI_SECRET_ACTION:-sync}"
audit_dir="${REEFOPS_SECRET_SYNC_AUDIT_DIR:-${XDG_STATE_HOME:-${HOME}/.local/state}/reefops/secret-sync}"

secret_count="$(
  SECRET_NAME="${secret_name}" yq -r \
    '[.secrets[] | select(.name == strenv(SECRET_NAME))] | length' \
    "${config_file}"
)"
if [[ "${secret_count}" != "1" ]]; then
  echo "El secreto debe aparecer exactamente una vez en la allowlist." >&2
  exit 1
fi
secret_path="$(
  SECRET_NAME="${secret_name}" yq -er \
    '.secrets[] | select(.name == strenv(SECRET_NAME)) | .openbao_path' \
    "${config_file}"
)"
github_repository="$(
  SECRET_NAME="${secret_name}" yq -er \
    '.secrets[] | select(.name == strenv(SECRET_NAME)) | .github_repository' \
    "${config_file}"
)"
github_environment="$(
  SECRET_NAME="${secret_name}" yq -r \
    '.secrets[] | select(.name == strenv(SECRET_NAME)) | .github_environment // ""' \
    "${config_file}"
)"
active_version="$(
  SECRET_NAME="${secret_name}" yq -er \
    '.secrets[] | select(.name == strenv(SECRET_NAME)) | .active_version' \
    "${config_file}"
)"
if [[ "${requested_version}" != "${active_version}" ]]; then
  echo "La versión solicitada no es la versión activa allowlisted." >&2
  exit 1
fi
if [[ "${requested_action}" != "sync" && "${requested_action}" != "delete" ]]; then
  echo "REEFOPS_CI_SECRET_ACTION debe ser sync o delete." >&2
  exit 1
fi

gh auth status --hostname github.com >/dev/null
bao token lookup >/dev/null

if [[ "${requested_action}" == "sync" ]]; then
  secret_metadata="$(bao kv metadata get -format=json "${secret_path}")"
  current_version="$(jq -r '.data.current_version' <<<"${secret_metadata}")"
  version_destroyed="$(
    jq -r --arg version "${active_version}" \
      '.data.versions[$version].destroyed' <<<"${secret_metadata}"
  )"
  version_deletion_time="$(
    jq -r --arg version "${active_version}" \
      '.data.versions[$version].deletion_time // "missing"' <<<"${secret_metadata}"
  )"
  if [[ "${current_version}" != "${active_version}" ||
    "${version_destroyed}" != "false" ||
    -n "${version_deletion_time}" ]]; then
    echo "La versión allowlisted no es la versión actual y activa en OpenBao." >&2
    exit 1
  fi
fi

operation_id="$(uuidgen | tr '[:upper:]' '[:lower:]')"
correlation_id="${REEFOPS_CORRELATION_ID:-${operation_id}}"
causation_id="${REEFOPS_CAUSATION_ID:-${operation_id}}"
started_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
install -d -m 0700 "${audit_dir}"
audit_file="${audit_dir}/operations.jsonl"
touch "${audit_file}"
chmod 0600 "${audit_file}"

gh_args=(secret set "${secret_name}" --repo "${github_repository}" --app actions)
if [[ -n "${github_environment}" ]]; then
  gh_args+=(--env "${github_environment}")
fi

result="success"
if [[ "${requested_action}" == "delete" ]]; then
  delete_args=(secret delete "${secret_name}" --repo "${github_repository}" --app actions)
  if [[ -n "${github_environment}" ]]; then
    delete_args+=(--env "${github_environment}")
  fi
  if ! gh "${delete_args[@]}"; then
    result="failure"
  fi
  variable_delete_args=(
    variable delete "${secret_name}_OPENBAO_VERSION"
    --repo "${github_repository}"
  )
  variable_list_args=(
    variable list
    --repo "${github_repository}"
    --json name
  )
  if [[ -n "${github_environment}" ]]; then
    variable_delete_args+=(--env "${github_environment}")
    variable_list_args+=(--env "${github_environment}")
  fi
  if gh "${variable_list_args[@]}" \
    --jq ".[] | select(.name == \"${secret_name}_OPENBAO_VERSION\") | .name" |
    rg -q .; then
    if ! gh "${variable_delete_args[@]}"; then
      result="failure"
    fi
  fi
elif ! bao kv get -field=value -version="${active_version}" "${secret_path}" |
  gh "${gh_args[@]}"; then
  result="failure"
else
  variable_args=(
    variable set "${secret_name}_OPENBAO_VERSION"
    --repo "${github_repository}"
    --body "${active_version}"
  )
  if [[ -n "${github_environment}" ]]; then
    variable_args+=(--env "${github_environment}")
  fi
  if ! gh "${variable_args[@]}"; then
    result="failure"
  fi
fi

finished_at="$(date -u '+%Y-%m-%dT%H:%M:%SZ')"
jq -cn \
  --arg operation_id "${operation_id}" \
  --arg actor "$(gh api user --jq .login)" \
  --arg secret_id "${secret_name}" \
  --arg secret_version "${active_version}" \
  --arg action "${requested_action}" \
  --arg repository "${github_repository}" \
  --arg environment "${github_environment}" \
  --arg started_at "${started_at}" \
  --arg finished_at "${finished_at}" \
  --arg result "${result}" \
  --arg correlation_id "${correlation_id}" \
  --arg causation_id "${causation_id}" \
  '{
    operation_id: $operation_id,
    actor: $actor,
    authentication: "openbao-token-and-github-cli",
    authorization: "openbao-policy-and-github-repository-role",
    github_app: "actions",
    secret_id: $secret_id,
    openbao_version: $secret_version,
    action: $action,
    github_repository: $repository,
    github_environment: $environment,
    started_at: $started_at,
    finished_at: $finished_at,
    result: $result,
    error: (if $result == "success" then null else "replication-operation-failed" end),
    correlation_id: $correlation_id,
    causation_id: $causation_id
  }' >>"${audit_file}"

if [[ "${result}" != "success" ]]; then
  exit 1
fi

echo "Réplica CI ${secret_name} sincronizada en ${github_repository}; operación ${operation_id}."
