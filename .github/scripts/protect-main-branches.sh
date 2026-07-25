#!/usr/bin/env bash
set -euo pipefail

project_root="$(git rev-parse --show-toplevel)"
github_owner="${REEFOPS_GITHUB_OWNER:?Define REEFOPS_GITHUB_OWNER con la organización}"
protection_file="${project_root}/.github/branch-protection.json"

gh auth status --hostname github.com >/dev/null
gh api "orgs/${github_owner}" --silent

for repository in reefops reefops-platform reefops-gitops; do
  repository_ref="${github_owner}/${repository}"
  gh api "repos/${repository_ref}/branches/main" --silent
  visibility="$(gh api "repos/${repository_ref}" --jq .visibility)"
  if [[ "${visibility}" == "private" ]]; then
    echo "Se omite ${repository_ref}: el plan actual no protege ramas privadas."
    continue
  fi
  gh api \
    --method PUT \
    "repos/${repository_ref}/branches/main/protection" \
    --input "${protection_file}" \
    --silent
done

echo "Protección pull-request-only aplicada a los repositorios compatibles."
