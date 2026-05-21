package main

import (
	"log"

	"github.com/gabriel/ttrpg-toolkit/backend/internal/config"
	"github.com/gabriel/ttrpg-toolkit/backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	addr := cfg.Addr()
	log.Printf("TTRPG Toolkit (Fiber) listening on http://%s", addr)

	app := server.New()
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server: %v", err)
	}
}
