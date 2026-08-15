package main

import (
	"log"

	"hbfittings-front/internal/app"
)

func main() {
	cfg := app.LoadConfig()
	server, err := app.New(cfg)
	if err != nil {
		log.Fatalf("start front api: %v", err)
	}
	log.Printf("hbfittings front api listening on %s", cfg.Addr)
	if err := server.Run(cfg.Addr); err != nil {
		log.Fatalf("run front api: %v", err)
	}
}
