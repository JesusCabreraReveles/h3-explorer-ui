# H3 Explorer UI — Frontend

SvelteKit + TypeScript (strict) + Tailwind CSS v4 + MapLibre GL JS + h3-js.

The frontend is **presentation only** — all authoritative H3 computation happens
in the Go backend. The browser always calls the API same-origin (`/api/*`):

- **dev** — Vite proxies `/api` and `/health` to the backend (`vite.config.ts`)
- **prod** — the adapter-node server proxies them (`src/hooks.server.ts`)

## Architecture

```
src/
├── lib/
│   ├── components/   # presentational UI (MapView, Sidebar, CoordinateSearch, …)
│   ├── services/     # typed API client (the only place that calls fetch)
│   ├── stores/       # reusable Svelte stores + feature orchestration
│   ├── types/        # TypeScript mirror of the backend DTOs
│   ├── utils/        # pure helpers (GeoJSON conversion, formatting)
│   └── config.ts     # constants (map style, defaults)
├── routes/           # +page / +layout (the explorer shell)
├── hooks.server.ts   # production API proxy
├── app.css           # Tailwind v4 + design tokens (dark theme)
└── app.html
```

Business logic stays out of components: they render stores and call store
actions (e.g. `indexPoint`, `changeResolution`), which talk to the API client.

## Develop

```bash
npm install
npm run dev            # http://localhost:5173 (needs the backend on :8080)
```

Point the dev proxy elsewhere with `API_PROXY_TARGET=http://host:port npm run dev`.

## Scripts

| Script            | Purpose                             |
| ----------------- | ----------------------------------- |
| `npm run dev`     | Vite dev server with API proxy      |
| `npm run build`   | Production build (adapter-node)     |
| `npm run preview` | Preview the production build        |
| `npm run check`   | `svelte-check` strict type checking |
| `npm run test`    | Vitest unit tests (services, utils) |
| `npm run lint`    | Prettier check + ESLint             |
| `npm run format`  | Prettier write                      |

## Production

```bash
npm run build
PORT=3000 API_PROXY_TARGET=http://localhost:8080 node build
```

Or run the whole stack via `docker compose up --build` from the repo root.
