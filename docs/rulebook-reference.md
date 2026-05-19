# TTRPG Dashboard Developer Rulebook

Canonical reference for contributors and AI agents. Cursor rules in `.cursor/rules/ttrpg-*.mdc` summarize enforceable points; this document is the full rulebook.

## 1. Architecture & deployment

- Monorepo: `frontend/` + `backend/`; atomic commits, separate execution.
- **No NPM at runtime:** Users and Docker run only the Go binary. UI is `npm run build` → embed → serve on `/`; API on `/api/*`.
- **Config:** Root `.env` (from `.env.example`). Frontend runtime values via `GET /api/config` only (no `config.yaml` in this repo).
- **Edge:** Caddy reverse-proxies in Docker Compose; Fiber listens on `TTRPG_SERVER_PORT` (default `8080`).

## 2. Directory topology

```
TTRPG-Toolkit/
├── frontend/              # SvelteKit
│   └── src/
│       ├── lib/
│       │   ├── components/
│       │   ├── stores/
│       │   ├── utils/
│       │   ├── api/
│       │   └── types/
│       └── routes/        # gm/, player/, wiki/ as features land
├── backend/
│   ├── cmd/
│   ├── internal/          # api, server, config; + services, repository, ws
│   ├── ui/static/         # go:embed (generated)
│   └── pkg/               # public shared libs (e.g. formula AST)
├── build/static/          # adapter-static output (generated)
├── docs/
├── deploy/
├── tools/
├── Dockerfile
└── docker-compose.yml
```

## 3. Naming & domain lexicon

See `.cursor/rules/ttrpg-naming.mdc` for the full lexicon and component patterns (`CampaignCreateModal`, `TE_*`, etc.).

## 4. Frontend standards

- State: `let` for local UI; stores for shared data.
- API: `$lib/api` only; types in `$lib/types`.
- UI copy: player/GM language; Create / Delete / Save; paired stats `Current / Max`.
- Theme: CSS variables (`--accent-primary`, `--accent-player`, `--accent-gm`); paper white + role accents on layouts.

## 5. Backend standards

- Flow: Handler → Service → Repository.
- Formulas: AST parser only; no `eval`; resolve `@field_id` in Service.
- Errors: wrap with `fmt.Errorf("...: %w", err)`; handlers return envelope.

## 6. API

- REST for CRUD; RPC-style POST for actions (sessions, dice, combat, sheet/card creation).
- Envelope: `status`, `data`, `message`, `errors`.
- WebSocket: one connection; messages `{ "type", "payload" }`.

## 7. Tooling

- **Ship build:** `go generate ./backend/ui/...` then `go build -C backend -o ttrpg-toolkit ./cmd` (or `.\prep.ps1` when `tools/build` exists).
- **UI dev (optional):** Go on `:8080` + `npm run dev` in `frontend/` (Vite proxies `/api`).
- **Commits:** Conventional Commits (`feat(combat): add next-turn endpoint`).
- **Deps:** No new package for &lt;20 lines of stdlib.
