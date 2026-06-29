# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.1.1] - 2026-06-28

### Added
- `SECURITY.md` with a private vulnerability disclosure policy.
- `CONTRIBUTING.md` and `CODE_OF_CONDUCT.md`.
- Issue and pull-request templates under `.github/`.
- Dependabot configuration for Go, npm, GitHub Actions and Docker.
- Root-level "Running the tests" section in the README.

### Changed
- Enabled branch protection on `master` (1 approving review, strict status
  checks, no force-push or deletion).

## [1.1.0] - 2026-06-28

### Fixed
- Regenerated `package-lock.json` to include missing `@emnapi` entries.

## [1.0.0] - 2026-06-27

### Added
- **Backend (Go, Clean Architecture):** API-first H3 service with full endpoint
  surface — coordinate indexing, inspection, `gridDisk`/`gridRing`/`gridPath`,
  parent/children, neighbors, `polygonToCells`, `cellsToMultiPolygon` — guarded
  by safety limits and an embedded OpenAPI 3.1 contract.
- **Frontend (SvelteKit + MapLibre):** interactive inspector (coordinate/index
  search, clickable parent · children · neighbors, map overlays, resolution
  explorer) and an H3 Playground (`gridDisk`/`gridRing` with k animation,
  `gridPath`, `polygonToCells`, index labels, boundary toggles).
- **Export:** GeoJSON FeatureCollection, boundary CSV, and H3 index lists.
- **Keyboard shortcuts** across the Explorer and Playground.
- **Tooling:** Docker / docker-compose stack, CI (build, race tests, lint) and a
  release workflow.

[Unreleased]: https://github.com/JesusCabreraReveles/h3-explorer-ui/compare/v1.1.1...HEAD
[1.1.1]: https://github.com/JesusCabreraReveles/h3-explorer-ui/compare/v1.1.0...v1.1.1
[1.1.0]: https://github.com/JesusCabreraReveles/h3-explorer-ui/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/JesusCabreraReveles/h3-explorer-ui/releases/tag/v1.0.0
