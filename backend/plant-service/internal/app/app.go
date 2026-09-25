package app

import (
	"log"
	"net/http"

	"github.com/alina965/pLaNtS/plant-service/internal/api"
	"github.com/alina965/pLaNtS/plant-service/internal/client"
	"github.com/alina965/pLaNtS/plant-service/internal/config"
	"github.com/alina965/pLaNtS/plant-service/internal/plant"
)

type App struct {
	server *http.Server
}

func New(config *config.Config) (*App, error) {
	perenualClient := client.NewPerenualClient(config.Timeout, config.PerenualKey)
	wikipediaClient := client.NewWikipediaClient(config.Timeout)
	plantService := plant.NewService(perenualClient, wikipediaClient)
	plantsHandler := api.NewPlantsHandler(plantService)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /plants", plantsHandler.GetSpeciesList)
	mux.HandleFunc("GET /plants/{id}", plantsHandler.GetSpeciesDetails)

	return &App{server: &http.Server{Addr: config.Addr, Handler: mux}}, nil
}

func (a *App) Run() error {
	log.Println("Server started on " + a.server.Addr)
	return a.server.ListenAndServe()
}
