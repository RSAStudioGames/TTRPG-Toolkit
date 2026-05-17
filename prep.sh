#!/usr/bin/env sh
# Prepares frontend embed + ttrpg-toolkit binary at repo root. No npm dev server.
set -e
cd "$(dirname "$0")"
exec go run -C tools ./build
