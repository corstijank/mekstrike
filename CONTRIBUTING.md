# Contributing

Thanks for helping improve Mekstrike! This guide is a quick entry point. For full contributor details, see AGENTS.md.

## Getting Started
- Read: AGENTS.md (Repository Guidelines)
- Quickstart: follow README.md (minikube profile `mekstrike`, then `./setup.sh` and `./deploy.sh`).
- Services live under language-specific folders (`gamemaster` Go, `battlefield`/`unit` Java, `ai-agent` Python, `ui` SvelteKit).

## Development Workflow
- Branching: create a short-lived feature branch from `main`.
- Commits: imperative, concise (e.g., "refactor ui", "fix ai agent").
- Protos: when editing files in `domain/`, regenerate via `./deploy.sh` (or `protoc` commands in AGENTS.md).
- Local checks:
  - UI: `cd ui && npm run lint && npm run format`.
  - Go: `gofmt -w .` and `go build ./...` in the changed module.
  - Java: `mvn -q -DskipTests package` in the changed service.
  - Python: follow PEP 8; consider `black .` locally.
- Kubernetes: if you change config, update `k8s/` manifests and note it in your PR.

## Pull Requests
- Describe the change, motivation, and scope.
- Link related issues. Add screenshots/GIFs for UI changes.
- Call out proto or k8s updates explicitly.
- Verify with `./deploy.sh <svc>` for affected services.

Questions or blockers? Add notes in the PR for reviewers.
