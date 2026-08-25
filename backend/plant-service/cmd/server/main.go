package main

import (
	"log"

	"github.com/alina965/pLaNtS/plant-service/internal/app"
	"github.com/alina965/pLaNtS/plant-service/internal/config"
)

func main() {
	cfg, err := config.New(".env")
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
