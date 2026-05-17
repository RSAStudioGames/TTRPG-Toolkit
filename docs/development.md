# Development Guide

## First-time setup

```bash
cp .env.example .env
make install
make build
```

## Environment variables

Copy `.env.example` to `.env` at the repo root. The Go backend loads this file on startup.

| Variable | Default | Purpose |
|----------|---------|---------|
| `TTRPG_SERVER_HOST` | `0.0.0.0` | HTTP bind address |
| `TTRPG_SERVER_PORT` | `8080` | HTTP port |

OS environment variables override values from `.env` (useful in Docker/Kubernetes).

## Development modes

### Frontend only (HMR)

```bash
make dev-frontend
```

Open http://localhost:5173 for hot reload while editing Svelte files.

### Full stack (production-like)

```bash
make build
make run
```

Open http://localhost:8080. Only the Go binary runs; no Node process.

### Concurrent dev

Run `make dev-frontend` in one terminal and `make run` in another to test API routes against the built UI.

## Docker

```bash
cp .env.example .env
docker compose up --build
```

## Makefile targets

| Target | Description |
|--------|-------------|
| `install` | Install frontend deps and tidy Go modules |
| `build-frontend` | Build SvelteKit static assets to `build/static/` |
| `build-backend` | Copy static assets and compile Go binary |
| `build` | Both build steps |
| `run` | Build and run the Go server |
| `dev-frontend` | SvelteKit dev server with HMR |
