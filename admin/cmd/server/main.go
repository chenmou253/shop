package main

import (
	"log"

	"rbac-admin/internal/app"
)

func main() {
	cfg := app.LoadConfig()

	server, err := app.New(cfg)
	if err != nil {
		log.Fatalf("start app: %v", err)
	}

	log.Printf("rbac-admin listening on %s", cfg.Addr)
	if err := server.Run(cfg.Addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
