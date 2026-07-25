#!/usr/bin/env bash
set -euo pipefail

github_owner="${REEFOPS_GITHUB_OWNER:?Define REEFOPS_GITHUB_OWNER con la organización}"

gh auth status --hostname github.com >/dev/null
gh api "orgs/${github_owner}" --silent
gh api \
  --method PATCH \
  "orgs/${github_owner}" \
  -f default_repository_permission=none \
  --silent

for repository in reefops reefops-platform reefops-gitops; do
  repository_ref="${github_owner}/${repository}"

  if ! gh repo view "${repository_ref}" >/dev/null 2>&1; then
    gh repo create "${repository_ref}" \
      --private \
      --disable-wiki \
      --description "ReefOps ${repository}"
  fi

  gh api \
    --method PATCH \
    "repos/${repository_ref}" \
    -F has_issues=true \
    -F has_projects=true \
    -F has_wiki=false \
    -F allow_squash_merge=true \
    -F allow_merge_commit=false \
    -F allow_rebase_merge=false \
    -F delete_branch_on_merge=true \
    -F allow_update_branch=true \
    --silent

  gh api \
    --method PUT \
    "repos/${repository_ref}/actions/permissions/workflow" \
    -f default_workflow_permissions=read \
    -F can_approve_pull_request_reviews=false \
    --silent
done

echo "Repositorios configurados sin modificar su visibilidad en ${github_owner}."
