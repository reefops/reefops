# GitHub automation guidance

- Read `docs/plataforma-github.md` before changing workflows or repository
  policy.
- Keep workflow-level `permissions: {}` and grant minimum permissions per job.
- Pin every third-party action to a full commit SHA and retain its version in a
  comment.
- Do not use pull-request secrets for untrusted fork code.
- GitHub Actions may build and publish but must not deploy directly to the
  local Kubernetes cluster.
- Promote releases by pull request to the GitOps repository.
- Never print credentials, SOPS plaintext or aquarium data.
- Validate workflows with `actionlint` and run `task validate-static`.
