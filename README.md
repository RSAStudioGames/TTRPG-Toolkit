# TTRPG Toolkit

A tabletop RPG companion app: **SvelteKit** UI embedded in a **Go (Fiber)** binary, with **Caddy** as the public HTTP edge in Docker.

## Stack

| Layer | Role |
|-------|------|
| **SvelteKit** | UI — compiled to static files and embedded in the Go binary |
| **Go (Fiber)** | Serves `/`, embedded SPA, and `/api/*` on `TTRPG_SERVER_PORT` (default `8080`) |
| **Caddy** | Reverse proxy on port `80` → Fiber (Docker Compose only) |

```
Browser → Caddy (:80)  →  Fiber (:8080)  →  embedded SvelteKit + API
          (Docker)         (ttrpg-toolkit)
```

## Project layout

```
TTRPG-Toolkit/
├── frontend/          # SvelteKit source
├── backend/           # Go Fiber app (cmd/, internal/, ui/static embed)
├── deploy/            # Caddyfile configs
├── build/static/      # SvelteKit build output (generated)
├── tools/buildfrontend/  # go generate helper
├── ttrpg-toolkit      # compiled binary (repo root, gitignored)
├── .env               # local config (from .env.example)
└── docker-compose.yml
```

## Prerequisites

- [Go](https://go.dev/) 1.22+
- [Node.js](https://nodejs.org/) (for `npm run build` during `go generate` only)
- [Docker](https://www.docker.com/) (optional, for Compose + Caddy)

## Quick start

```bash
cp .env.example .env
cd frontend && npm ci
cd ../backend && go mod download
cd ..

go generate ./backend/ui/...
go build -C backend -o ttrpg-toolkit ./cmd
```

Run from the **repository root**:

| OS | Command |
|----|---------|
| Windows | `ttrpg-toolkit.exe` |
| macOS / Linux | `./ttrpg-toolkit` |

Open **http://localhost:8080** (or the port set in `.env`).

## Rebuilding

| You changed | Before `go build` |
|-------------|-------------------|
| Svelte / `frontend/` | `go generate ./backend/ui/...` |
| Go / `backend/` only | `go build -C backend -o ttrpg-toolkit ./cmd` |

`go generate` runs `npm run build` and copies assets into `backend/ui/static/` for `go:embed`. It is a **compile step**, not a dev server.

**Do not use:** `npm run dev`, Makefile, or any Node HTTP server at runtime.

## Environment

Copy [`.env.example`](.env.example) to `.env` at the repo root.

| Variable | Default | Purpose |
|----------|---------|---------|
| `TTRPG_SERVER_HOST` | `0.0.0.0` | Fiber bind address |
| `TTRPG_SERVER_PORT` | `8080` | Fiber listen port |
| `CADDY_PORT` | `80` | Host port for Caddy (Docker Compose) |

## Docker (Fiber + Caddy)

```bash
cp .env.example .env
docker compose up --build
```

Open **http://localhost** — Caddy proxies to the Fiber container on port `8080`. The Go app is not published directly to the host.

Optional: run Caddy locally in front of a binary already listening on `:8080`:

```bash
caddy run --config deploy/Caddyfile.local
```

## Documentation

- [Development guide](docs/development.md) — architecture, rebuild workflow, Docker stages
- [Rulebook reference](docs/rulebook-reference.md) — naming, API envelope, theming
