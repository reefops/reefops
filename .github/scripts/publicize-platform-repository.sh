#!/usr/bin/env bash
set -euo pipefail

github_owner="${REEFOPS_GITHUB_OWNER:?Define REEFOPS_GITHUB_OWNER con la organización}"
repository_ref="${github_owner}/reefops-platform"

gh auth status --hostname github.com >/dev/null
gh api "orgs/${github_owner}" --silent
repository_url="$(gh repo view "${repository_ref}" --json sshUrl --jq .sshUrl)"

audit_dir="$(mktemp -d)"
trap 'rm -rf "${audit_dir}"' EXIT
git clone --quiet "${repository_url}" "${audit_dir}/repository"

test -f "${audit_dir}/repository/LICENSE"
rg -q 'Apache License' "${audit_dir}/repository/LICENSE"

for forbidden_path in .sops.yaml clusters; do
  if [[ -e "${audit_dir}/repository/${forbidden_path}" ]]; then
    echo "Contenido no publicable presente: ${forbidden_path}" >&2
    exit 1
  fi
done
if git -C "${audit_dir}/repository" log \
  --all \
  --format= \
  --name-only \
  -- .sops.yaml clusters |
  rg -q '.'; then
  echo "El historial contiene paths de configuración no publicables." >&2
  exit 1
fi

gitleaks git --no-banner --redact --exit-code 1 "${audit_dir}/repository"
trivy fs \
  --scanners secret \
  --exit-code 1 \
  --skip-dirs .git \
  "${audit_dir}/repository"

platform_commit="$(git -C "${audit_dir}/repository" rev-parse HEAD)"
gitops_source="$(
  gh api \
    "repos/${github_owner}/reefops-gitops/contents/clusters/local/workloads/platform-source.yaml?ref=main" \
    --jq .content |
    base64 --decode
)"
promoted_commit="$(yq -r '.spec.ref.commit' <<<"${gitops_source}")"
if [[ "${platform_commit}" != "${promoted_commit}" ]]; then
  echo "GitOps no fija el HEAD de plataforma auditado." >&2
  exit 1
fi

current_visibility="$(gh api "repos/${repository_ref}" --jq .visibility)"
if [[ "${current_visibility}" != "private" &&
  "${current_visibility}" != "public" ]]; then
  echo "Visibilidad inesperada: ${current_visibility}" >&2
  exit 1
fi

if [[ "${current_visibility}" == "private" ]]; then
  gh api \
    --method PATCH \
    "repos/${repository_ref}" \
    -f visibility=public \
    --silent
fi

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

echo "${repository_ref} público y con seguridad gratuita habilitada."
