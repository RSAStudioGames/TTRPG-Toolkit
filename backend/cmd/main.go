package main

import (
	"log"
	"net/http"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/config"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	addr := cfg.Addr()
	log.Printf("TTRPG Toolkit listening on http://%s", addr)

	handler := server.New()
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server: %v", err)
	}
}
