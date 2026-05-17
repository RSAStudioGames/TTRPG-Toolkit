# Rulebook Reference (Summary)

This project follows the TTRPG Dashboard Developer Rulebook. Key points for contributors:

## Architecture

- Monorepo: `frontend/` (SvelteKit) + `backend/` (Go Fiber)
- Production: Fiber binary embeds static SPA on `/`, API on `/api/*`; Caddy reverse-proxies public HTTP
- No `npm` in production runtime

## Configuration

- Root `.env` for backend variables (see `.env.example`)
- No `config.yaml`
- Frontend receives runtime config via `GET /api/config` only

## Naming

- **Components:** `[Feature][Type].svelte` (e.g. `JoinGameButton.svelte`)
- **Routes:** kebab-case; **Go files:** snake_case
- **Domain terms:** Campaign, Character, NPC, Enemy, Wiki (UI: "Game Reference")

## API envelope

```json
{
  "status": "success" | "error",
  "data": { },
  "message": "Human readable message",
  "errors": [ ]
}
```

## Theming

Use CSS variables (`--accent-primary`, `--accent-player`, `--accent-gm`). Do not hardcode accent hex in components.
