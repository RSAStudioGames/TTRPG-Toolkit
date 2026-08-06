#!/usr/bin/env sh
# One-terminal UI HMR preview: ensure Postgres, rebuild Go API, then Vite on :3030
# (proxies /api → :8080). Not for production — use prep + ./ttrpg-toolkit for ship-parity.
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

echo "==> Starting Postgres (docker compose -f docker-compose.db.yml)..."
if ! command -v docker >/dev/null 2>&1; then
	echo "docker is required for make dev (Postgres on :43021)." >&2
	exit 1
fi
docker compose -f docker-compose.db.yml up -d

echo "==> Waiting for Postgres..."
i=0
while [ "$i" -lt 60 ]; do
	if docker compose -f docker-compose.db.yml exec -T postgres \
		pg_isready -U ttrpg -d ttrpg_toolkit >/dev/null 2>&1; then
		break
	fi
	i=$((i + 1))
	sleep 1
done
if [ "$i" -ge 60 ]; then
	echo "Postgres did not become ready on :43021." >&2
	exit 1
fi

echo "==> Building Go binary..."
go build -C backend -o ../ttrpg-toolkit ./cmd

echo "==> Starting API (./ttrpg-toolkit on :8080)..."
./ttrpg-toolkit &
API_PID=$!

# Fail fast if the binary exits immediately.
sleep 0.4
if ! kill -0 "${API_PID}" 2>/dev/null; then
	status=0
	wait "${API_PID}" 2>/dev/null || status=$?
	echo "API failed to start (exit ${status}). Check the log above (DB, port bind, etc.)." >&2
	API_PID=""
	exit 1
fi

if [ ! -d frontend/node_modules ]; then
	echo "==> Installing frontend deps..."
	npm ci --prefix frontend
fi

echo "==> Starting Vite HMR on http://localhost:3030 (/api → Go :8080)."
cd frontend
# 130/143 = SIGINT/SIGTERM (Ctrl+C) — treat as clean exit for make.
npm run dev || {
	code=$?
	if [ "$code" -eq 130 ] || [ "$code" -eq 143 ]; then
		exit 0
	fi
	exit "$code"
}
