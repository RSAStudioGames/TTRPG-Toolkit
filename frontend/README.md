# Frontend

SvelteKit sources. Compiled via `go generate ./backend/ui/...` from the repo root (runs `npm run build` here and copies output for `go:embed`).

For UI HMR while editing, use `make dev` from the repo root (starts Go API + Vite together). For ship-parity, regenerate and rebuild the Go binary — see the root [README](../README.md).
