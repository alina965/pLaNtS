package main

import (
	"context"
	"log"

	"github.com/alina965/pLaNtS/user-service/internal/app"
	"github.com/alina965/pLaNtS/user-service/internal/config"
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
	defer func() {
		if err := application.Shutdown(context.Background()); err != nil {
			log.Printf("shutdown error: %v", err)
		}
	}()
	if err := application.Run(); err != nil {
		log.Fatal(err)
	}
}
