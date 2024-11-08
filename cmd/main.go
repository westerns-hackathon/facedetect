package main

import (
	"face-detection/internal/config"
	"face-detection/internal/handler"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	h := handler.NewHandler()
	srv := handler.NewServer(h)
	if err := srv.Run(cfg.Host, cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}
}
