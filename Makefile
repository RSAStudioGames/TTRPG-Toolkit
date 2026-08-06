.PHONY: dev

# Ensure Postgres, rebuild Go API (:8080), start Vite HMR (:3030). Ctrl+C stops API+Vite.
# Browse http://localhost:3030
dev:
	@sh scripts/dev.sh
