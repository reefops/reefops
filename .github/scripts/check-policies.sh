#!/usr/bin/env bash
set -euo pipefail

project_root="$(git rev-parse --show-toplevel)"
workflow_root="${project_root}/.github/workflows"

yq -e '.version == 2' "${project_root}/.github/dependabot.yml" >/dev/null
yq -e '.blank_issues_enabled != null' \
  "${project_root}/.github/ISSUE_TEMPLATE/config.yml" >/dev/null
yq -e \
  '.organization.base_permission == "none" and
   .repositories.reefops.visibility == "public" and
   .repositories.reefops-platform.visibility == "public" and
   .repositories.reefops-gitops.visibility == "private" and
   .branches.default == "main" and
   .branches.strategy == "github-flow" and
   .branches.merge_method == "squash" and
   .connectivity.local_inbound_from_github == false and
   .connectivity.flux_repository_access == "read-only" and
   .secrets.runtime_authority == "openbao-local" and
   .secrets.github_is_replica_for_ci_only == true' \
  "${project_root}/.github/repository-governance.yaml" >/dev/null
yq -e \
  '[.secrets[] |
    select(
      (.name | test("^[A-Z][A-Z0-9_]*$") | not) or
      (.openbao_path | test("^ci/") | not) or
      (.github_repository | test("^[A-Za-z0-9_.-]+/[A-Za-z0-9_.-]+$") | not) or
      (.active_version | type != "!!int") or
      (.active_version < 1) or
      (.github_environment != null and (.github_environment | type != "!!str"))
    )
  ] | length == 0' \
  "${project_root}/.github/ci-secrets.yaml" >/dev/null

while IFS= read -r workflow_file; do
  yq -e \
    '(.permissions | type) == "!!map" and (.permissions | length) == 0' \
    "${workflow_file}" >/dev/null
  yq -e \
    '[.jobs | to_entries | .[].value | select(.permissions == null)] | length == 0' \
    "${workflow_file}" >/dev/null
done < <(find "${workflow_root}" -type f \( -name '*.yml' -o -name '*.yaml' \) -print)

while IFS=$'\t' read -r secret_name active_version; do
  while IFS= read -r workflow_file; do
    while IFS= read -r job_id; do
      job_condition="$(
        JOB_ID="${job_id}" yq -r '.jobs[strenv(JOB_ID)].if // ""' \
          "${workflow_file}"
      )"
      required_condition="vars.${secret_name}_OPENBAO_VERSION == '${active_version}'"
      if [[ "${job_condition}" != *"${required_condition}"* ]]; then
        echo "El job ${job_id} consume ${secret_name} sin condición ${required_condition}." >&2
        exit 1
      fi
    done < <(
      SECRET_REF="secrets.${secret_name}" yq -r \
        '.jobs | to_entries | .[] |
         select((.value | to_json) | contains(strenv(SECRET_REF))) |
         .key' "${workflow_file}"
    )
  done < <(find "${workflow_root}" -type f \( -name '*.yml' -o -name '*.yaml' \) -print)
done < <(
  yq -r '.secrets[] | [.name, .active_version] | @tsv' \
    "${project_root}/.github/ci-secrets.yaml"
)

while IFS= read -r use_line; do
  action_ref="${use_line##*@}"
  action_ref="${action_ref%% *}"
  if [[ ! "${action_ref}" =~ ^[0-9a-f]{40}$ ]]; then
    echo "Action sin fijar por SHA completo: ${use_line}" >&2
    exit 1
  fi
done < <(rg --no-filename '^[[:space:]]*uses:' "${workflow_root}")

if rg -n \
  'kubectl[[:space:]]+(apply|delete|patch|replace)|helm[[:space:]]+(install|upgrade|uninstall)|flux[[:space:]]+(bootstrap|reconcile)|task[[:space:]]+(bootstrap|reconcile|gitops-seed|platform-seed|ci-secret-sync)|bootstrap-flux\.sh|seed-(gitops|platform)-repository\.sh|sync-secret-from-openbao\.sh|KUBECONFIG|runs-on:[[:space:]]*self-hosted|pull_request_target:' \
  "${workflow_root}"; then
  echo "Workflow incompatible con pull-only o con la separación de secretos." >&2
  exit 1
fi

if rg -n \
  '(^|[^0-9])(10|192\.168|172\.(1[6-9]|2[0-9]|3[01]))\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}([^0-9]|$)|/Volumes/|smb://|//[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}/' \
  "${project_root}/README.md" "${project_root}/docs"; then
  echo "La documentación pública contiene topología operativa local." >&2
  exit 1
fi

echo "Políticas GitHub validadas."
