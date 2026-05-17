.PHONY: install build build-frontend build-backend run dev-frontend dev

BINARY := build/bin/ttrpg-toolkit
STATIC_SRC := build/static
STATIC_DST := backend/ui/static

install:
	cd frontend && npm install
	cd backend && go mod tidy

build-frontend:
	cd frontend && npm run build

build-backend: build-frontend
ifeq ($(OS),Windows_NT)
	if not exist backend\ui\static mkdir backend\ui\static
	xcopy /E /I /Y build\static backend\ui\static
else
	rm -rf $(STATIC_DST)
	mkdir -p $(STATIC_DST)
	cp -r $(STATIC_SRC)/. $(STATIC_DST)/
endif
	cd backend && go build -o ../$(BINARY) ./cmd

build: build-backend

run: build
ifeq ($(OS),Windows_NT)
	$(BINARY).exe
else
	./$(BINARY)
endif

dev-frontend:
	cd frontend && npm run dev

dev: dev-frontend
