#!/usr/bin/env sh
# One-terminal UI HMR preview: rebuild Go API, then Vite (proxies /api → :8080).
# Not for production — use prep + ./ttrpg-toolkit for ship-parity runs.
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

API_PID=""

cleanup() {
	if [ -n "${API_PID}" ]; then
		kill "${API_PID}" 2>/dev/null || true
		wait "${API_PID}" 2>/dev/null || true
	fi
}
trap cleanup EXIT INT TERM

echo "==> Building Go binary..."
go build -C backend -o ../ttrpg-toolkit ./cmd

echo "==> Starting API (./ttrpg-toolkit)..."
./ttrpg-toolkit &
API_PID=$!

# Fail fast if the binary exits immediately (e.g. port in use).
sleep 0.4
if ! kill -0 "${API_PID}" 2>/dev/null; then
	wait "${API_PID}" 2>/dev/null || true
	echo "API failed to start. Is port ${TTRPG_SERVER_PORT:-8080} already in use?" >&2
	API_PID=""
	exit 1
fi

if [ ! -d frontend/node_modules ]; then
	echo "==> Installing frontend deps..."
	npm ci --prefix frontend
fi

echo "==> Starting Vite HMR (use the URL Vite prints; /api proxies to the Go API)."
cd frontend
npm run dev
