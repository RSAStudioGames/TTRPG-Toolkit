I'm using AI to develop most of this because I don't know how to program all that much.
This is supposed to be a tool to help keep things organized for your TTRPG game. You can enter in all the game rules and such and this toolkit can then use that info for data-driven tools based off the specific ruleset and such. I'm still planning things out as I go along.

There's too many tools out there that have too narrow of a scope. I'm trying to combine all those different aspects together into one full-fledged campaign manager that has support for any game, including custom stuff.

---

# TTRPG Toolkit

A tabletop RPG companion app: **SvelteKit** UI embedded in a **Go (Fiber)** binary, with **Caddy** as the public HTTP edge in Docker.

## Stack

| Layer | Role |
|-------|------|
| **SvelteKit + TypeScript** | UI — strict TS, `lib/api` + `lib/types`, embedded in the Go binary |
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
├── tools/build/         # full prep tool (`go run -C tools ./build`)
├── prep.ps1 / prep.sh   # shortcuts to the prep tool
├── ttrpg-toolkit      # compiled binary (repo root, gitignored)
├── .env               # local config (from .env.example)
└── docker-compose.yml
```

## Prerequisites

- [Go](https://go.dev/) 1.22+
- [Node.js](https://nodejs.org/) (for `npm run build` during `go generate` only)
- [Docker](https://www.docker.com/) (optional, for Compose + Caddy)

## Quick start

One command prepares everything (deps, frontend build, embed, Go binary at repo root):

```bash
go run -C tools ./build
```

Windows: `.\prep.ps1` · macOS/Linux: `./prep.sh`

Run from the **repository root**:

| OS | Command |
|----|---------|
| Windows | `ttrpg-toolkit.exe` |
| macOS / Linux | `./ttrpg-toolkit` |

Open **http://localhost:8080** (or the port set in `.env`).

## Rebuilding

After any code change, run the prep tool again:

```bash
go run -C tools ./build
```

Manual steps (if you prefer): `go generate ./backend/ui/...` when only the UI changed, then `go build -C backend -o ttrpg-toolkit ./cmd`.

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
- [TypeScript guide](docs/typescript.md) — API client, types
- [Rulebook reference](docs/rulebook-reference.md) — naming, API envelope, theming
