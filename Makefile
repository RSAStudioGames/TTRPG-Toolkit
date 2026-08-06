.PHONY: dev

# Rebuild Go API + start Vite HMR in one terminal (Ctrl+C stops both).
# Browse the Vite URL (usually http://localhost:5173), not :8080.
dev:
	@sh scripts/dev.sh
