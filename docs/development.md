# Development Guide

## Architecture

```
Browser → Caddy (:80) → Go Fiber (:8080) → embedded SvelteKit static + /api/*
```

- **Fiber** serves the embedded build and API. This is the only application process you compile.
- **Caddy** is the public entrypoint in Docker Compose (TLS-ready, reverse proxy). It does not replace the Go binary.

## Runtime

No Node at runtime. No `npm run dev`.

## Environment

`.env` at repo root (from `.env.example`):

| Variable | Default | Purpose |
|----------|---------|---------|
| `TTRPG_SERVER_HOST` | `0.0.0.0` | Fiber bind address |
| `TTRPG_SERVER_PORT` | `8080` | Fiber port (internal in Docker) |
| `CADDY_PORT` | `80` | Host port mapped to Caddy |

## Build

```bash
go generate ./backend/ui/...    # when frontend changed
go build -C backend -o ttrpg-toolkit ./cmd
```

`go generate` runs `tools/buildfrontend` (`npm run build` + copy into `backend/ui/static/` for `go:embed`).

Restart the binary after rebuilding.

## Docker Compose

1. **frontend** stage: `npm run build`
2. **backend** stage: embed static, `go build` Fiber binary
3. **app** service: runs binary on 8080 (not published to host)
4. **caddy** service: `deploy/Caddyfile` → `reverse_proxy app:8080`

## Direct vs Caddy

| Mode | How |
|------|-----|
| Direct | Run `./ttrpg-toolkit` (repo root) → http://localhost:8080 |
| Caddy (local) | Binary on 8080 + `caddy run --config deploy/Caddyfile.local` → http://localhost |
| Caddy (Docker) | `docker compose up` → http://localhost |
