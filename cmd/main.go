package main

import (
	"face-detection/internal/config"
	"face-detection/internal/db"
	"face-detection/internal/handler"
	"github.com/Kagami/go-face"
	"log"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}
	database := db.NewDB(cfg)
	if err = database.Open(); err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	recognizer, err := face.NewRecognizer("internal/facedetect/models")
	if err != nil {
		log.Fatalf("Ошибка при инициализации recognizer: %v", err)
	}
	defer recognizer.Close()
	h := handler.NewHandler(recognizer, database)
	srv := handler.NewServer(h)
	if err := srv.Run(cfg.Host, cfg.Port); err != nil {
		log.Fatalf("failed to run server: %v\n", err)
	}
}
