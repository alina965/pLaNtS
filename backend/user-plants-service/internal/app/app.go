package app

import (
	"context"
	"log"
	"net/http"

	"github.com/alina965/pLaNtS/user-plants-service/internal/api"
	"github.com/alina965/pLaNtS/user-plants-service/internal/config"
	"github.com/alina965/pLaNtS/user-plants-service/internal/jwt"
	"github.com/alina965/pLaNtS/user-plants-service/internal/repository"
	"github.com/alina965/pLaNtS/user-plants-service/internal/user_plant"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db     *pgxpool.Pool
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	plantRepo := repository.NewUserPlantsRepository(db)
	eventsRepo := repository.NewWateringEventsRepository(db)
	plantService := user_plant.NewUserPlantsService(plantRepo, eventsRepo, db)
	plantsHandler := api.NewUserPlantsHandler(plantService)
	jwtService := jwt.NewService([]byte(cfg.JwtSecret))
	auth := api.AuthMiddleware(jwtService)

	mux := http.NewServeMux()
	mux.Handle("POST /user-plants", auth(http.HandlerFunc(plantsHandler.CreateUserPlant)))
	mux.Handle("GET /user-plants", auth(http.HandlerFunc(plantsHandler.GetUserPlants)))
	mux.Handle("GET /user-plants/{id}", auth(http.HandlerFunc(plantsHandler.GetUserPlantByID)))
	mux.Handle("PATCH /user-plants/{id}", auth(http.HandlerFunc(plantsHandler.UpdateUserPlant)))
	mux.Handle("DELETE /user-plants/{id}", auth(http.HandlerFunc(plantsHandler.DeleteUserPlant)))
	mux.Handle("POST /user-plants/{id}/water", auth(http.HandlerFunc(plantsHandler.MarkWatered)))
	mux.Handle("GET /user-plants/{id}/watering-events", auth(http.HandlerFunc(plantsHandler.GetWateringEvents)))

	return &App{
		db: db,
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
		},
	}, nil
}

func (a *App) Run() error {
	log.Println("Server started on " + a.server.Addr)
	return a.server.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) error {
	_ = a.server.Shutdown(ctx)
	a.db.Close()
	return nil
}
