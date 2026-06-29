# Contributing to H3 Explorer UI

Thanks for your interest in contributing! This document explains how to set up
your environment, the conventions we follow, and how to get a change merged.

By participating in this project you agree to abide by our
[Code of Conduct](./CODE_OF_CONDUCT.md).

## Ways to contribute

- 🐛 **Report bugs** — open a [bug report](https://github.com/JesusCabreraReveles/h3-explorer-ui/issues/new/choose).
- 💡 **Propose features** — open a feature request and describe the use case.
- 📝 **Improve docs** — README, OpenAPI, code comments.
- 🔧 **Send code** — see the workflow below.

## Development setup

### Prerequisites

- **Go 1.25+** (developed on 1.26) and a C compiler (`uber/h3-go` uses cgo).
- **Node 20.19+ / 22.12+**.
- **Docker** (optional, for the containerised path).

### Run locally

```bash
# backend
cd backend
go run ./cmd/server          # serves :8080

# frontend
cd frontend
npm install
npm run dev                  # serves :5173, proxies /api → :8080
```

## Architecture rules

This is a monorepo with a strict separation of concerns. Please respect it:

- **All H3 business logic lives in the backend `internal/service/h3` package.**
  The frontend is pure presentation — it must never re-implement H3 logic.
- Backend dependencies point **inward** (`domain ← service ← handler ← main`).
  The domain has no knowledge of HTTP or `h3-go`.
- Frontend business logic belongs in stores/services, not in UI components.

## Before you open a pull request

Run the full check suite locally and make sure it is green:

**Backend**

```bash
cd backend
go test -race -cover ./...
go vet ./...
golangci-lint run
gofmt -l .                   # should print nothing
```

**Frontend**

```bash
cd frontend
npm run check                # svelte-check (strict)
npm run test                 # vitest
npm run lint                 # prettier + eslint
npm run build                # production build
```

## Commit & PR conventions

- Use [Conventional Commits](https://www.conventionalcommits.org) for messages,
  e.g. `feat(api): add cells-to-multi-polygon endpoint` or `fix: handle k=0`.
- Keep PRs focused; one logical change per PR.
- Update `CHANGELOG.md` (under `[Unreleased]`) for user-facing changes.
- Add or update tests for any behavioural change.
- Fill in the pull-request template and link the related issue.

## Versioning

This project follows [Semantic Versioning](https://semver.org). Maintainers cut
releases by tagging `vMAJOR.MINOR.PATCH`, which triggers the release workflow.

## Questions

Open a [discussion or issue](https://github.com/JesusCabreraReveles/h3-explorer-ui/issues)
— we're happy to help.
