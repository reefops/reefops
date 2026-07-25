#!/usr/bin/env bash
set -euo pipefail

project_root="$(git rev-parse --show-toplevel)"
github_owner="${REEFOPS_GITHUB_OWNER:?Define REEFOPS_GITHUB_OWNER con la organización}"
repository_ref="${github_owner}/reefops"

if [[ -n "$(git -C "${project_root}" status --porcelain)" ]]; then
  echo "El worktree debe estar limpio antes de publicar." >&2
  exit 1
fi
if gh api "repos/${repository_ref}/commits" --silent >/dev/null 2>&1; then
  echo "El repositorio remoto ya contiene commits; se rechaza el primer push." >&2
  exit 1
fi

test -f "${project_root}/LICENSE"
rg -q 'Apache License' "${project_root}/LICENSE"
gitleaks git --no-banner --redact --exit-code 1 "${project_root}"
trivy fs --scanners secret --exit-code 1 --skip-dirs .git "${project_root}"

repository_ssh_url="$(gh repo view "${repository_ref}" --json sshUrl --jq .sshUrl)"
if git -C "${project_root}" remote get-url origin >/dev/null 2>&1; then
  if [[ "$(git -C "${project_root}" remote get-url origin)" != "${repository_ssh_url}" ]]; then
    echo "El remote origin existente no coincide con ${repository_ref}." >&2
    exit 1
  fi
else
  git -C "${project_root}" remote add origin "${repository_ssh_url}"
fi

git -C "${project_root}" push --set-upstream origin main
echo "Producto publicado en ${repository_ref}:main."
