#!/usr/bin/env bash
set -euo pipefail

github_owner="${REEFOPS_GITHUB_OWNER:?Define REEFOPS_GITHUB_OWNER con la organización}"
repository_ref="${github_owner}/reefops"

gh auth status --hostname github.com >/dev/null
repository_url="$(gh repo view "${repository_ref}" --json sshUrl --jq .sshUrl)"

audit_dir="$(mktemp -d)"
trap 'rm -rf "${audit_dir}"' EXIT
git clone --quiet "${repository_url}" "${audit_dir}/repository"

test -f "${audit_dir}/repository/LICENSE"
rg -q 'Apache License' "${audit_dir}/repository/LICENSE"

for forbidden_path in .sops.yaml infra clusters; do
  if [[ -e "${audit_dir}/repository/${forbidden_path}" ]]; then
    echo "Contenido no publicable presente: ${forbidden_path}" >&2
    exit 1
  fi
done
if git -C "${audit_dir}/repository" log \
  --all \
  --format= \
  --name-only \
  -- .sops.yaml infra clusters |
  rg -q '.'; then
  echo "El historial contiene paths de infraestructura no publicables." >&2
  exit 1
fi

commit_count="$(git -C "${audit_dir}/repository" rev-list --count origin/main)"
if [[ "${commit_count}" != "1" ]]; then
  echo "main debe contener un único commit raíz." >&2
  exit 1
fi
while IFS= read -r remote_ref; do
  if [[ "${remote_ref}" == "origin/HEAD" ||
    "${remote_ref}" == "origin/main" ]]; then
    continue
  fi
  if ! git -C "${audit_dir}/repository" \
    merge-base --is-ancestor origin/main "${remote_ref}"; then
    echo "La rama ${remote_ref} no desciende del root público." >&2
    exit 1
  fi
done < <(
  git -C "${audit_dir}/repository" \
    for-each-ref refs/remotes/origin --format='%(refname:short)'
)

gitleaks git --no-banner --redact --exit-code 1 "${audit_dir}/repository"
trivy fs \
  --scanners secret \
  --exit-code 1 \
  --skip-dirs .git \
  "${audit_dir}/repository"

gh api \
  --method PATCH \
  "repos/${repository_ref}" \
  -f visibility=public \
  --silent
gh api \
  --method PUT \
  "repos/${repository_ref}/vulnerability-alerts" \
  --silent
gh api \
  --method PUT \
  "repos/${repository_ref}/automated-security-fixes" \
  --silent
gh api \
  --method PATCH \
  "repos/${repository_ref}" \
  -F 'security_and_analysis[secret_scanning][status]=enabled' \
  -F 'security_and_analysis[secret_scanning_push_protection][status]=enabled' \
  -F 'security_and_analysis[secret_scanning_validity_checks][status]=enabled' \
  -F 'security_and_analysis[dependabot_security_updates][status]=enabled' \
  --silent

echo "${repository_ref} publicado desde un único commit raíz auditado."
