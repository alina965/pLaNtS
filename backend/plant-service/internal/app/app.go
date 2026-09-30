package app

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/alina965/pLaNtS/plant-service/internal/api"
	cache2 "github.com/alina965/pLaNtS/plant-service/internal/cache"
	"github.com/alina965/pLaNtS/plant-service/internal/client"
	"github.com/alina965/pLaNtS/plant-service/internal/config"
	"github.com/alina965/pLaNtS/plant-service/internal/plant"
	"github.com/redis/go-redis/v9"
)

type App struct {
	server      *http.Server
	redisClient *redis.Client
}

func New(config *config.Config) (*App, error) {
	redisClient := redis.NewClient(&redis.Options{Addr: config.RedisAddr})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Printf("redis ping failed: %v", err)
	} else {
		log.Println("redis connected")
	}

	cache := cache2.NewRedisCache(redisClient)

	perenualClient := client.NewPerenualClient(config.Timeout, config.PerenualKey, cache, config.CacheTTLList, config.CacheTTLDetails)
	wikipediaClient := client.NewWikipediaClient(config.Timeout)
	plantService := plant.NewService(perenualClient, wikipediaClient)
	plantsHandler := api.NewPlantsHandler(plantService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /plants", plantsHandler.GetSpeciesList)
	mux.HandleFunc("GET /plants/{id}", plantsHandler.GetSpeciesDetails)

	return &App{server: &http.Server{Addr: config.Addr, Handler: mux}, redisClient: redisClient}, nil
}

func (a *App) Run() error {
	log.Println("Server started on " + a.server.Addr)
	return a.server.ListenAndServe()
}
