# Repository Guidelines

## Project Structure & Module Organization
- `domain/`: Protocol Buffers shared across services.
- `gamemaster/` (Go): REST/game orchestration; clients in `gamemaster/clients`.
- `library/`, `importer/`, `mediaproxy/` (Go): data APIs and utilities.
- `armybuilder/`, `battlefield/`, `unit/` (Java): Dapr-enabled services; Boot apps for battlefield/unit.
- `ai-agent/` (Python/FastAPI): reacts to game events.
- `ui/` (SvelteKit): web client.
- `k8s/`: manifests; `k8s/mekstrike-base` and `k8s/platform` set up infra.
- Root scripts: `setup.sh`, `deploy.sh`.

## Build, Test, and Development Commands
- Bootstrap platform: `./setup.sh` (Traefik, Dapr, Redis, OTEL, Jaeger; expects a running Minikube profile `mekstrike`).
- Generate protos: handled by `./deploy.sh`; manual example:
  - `protoc --proto_path=. --go_out=. domain/unit/unit.proto`
- Build & deploy all: `./deploy.sh`
- Incremental deploy: `./deploy.sh ui` (replace with `gamemaster`, `armybuilder`, etc.)
- Local dev examples:
  - UI: `cd ui && npm run dev`
  - Go: `go run ./gamemaster`
  - Java (battlefield/unit): `cd <svc> && mvn package && java -jar target/mekstrike-<svc>-exec.jar`
  - Java (armybuilder): `cd armybuilder && mvn exec:java`

## Coding Style & Naming Conventions
- Go: format with `gofmt`; idiomatic package names (lowercase, no underscores). Exported identifiers use CamelCase.
- Java: package `net.mekstrike.*`; classes PascalCase, methods camelCase. Use Maven profiles as defined.
- UI: run `npm run lint` and `npm run format` (ESLint + Prettier). Svelte components PascalCase.
- Python: follow PEP 8; prefer type hints. Consider `black` formatting locally.

## Testing Guidelines
- Current tests are minimal. Prefer adding:
  - Go: `_test.go` beside sources; run `go test ./...`
  - Java: `src/test/java`; run `mvn test`
  - UI: add Vitest when needed; run via npm script
  - Python: `pytest` recommended for new modules

## Commit & Pull Request Guidelines
- Commits: short, imperative subject (e.g., "refactor ui", "fix ai agent").
- PRs: include a clear description, linked issues, and screenshots for UI changes. Note any proto or k8s updates. Verify locally via `./deploy.sh <svc>`.

## Security & Configuration Tips
- Do not hard-code service URLs; use Dapr sidecars/Kubernetes services. Store secrets in Kubernetes, not in code. Keep changes scoped and reproducible.
