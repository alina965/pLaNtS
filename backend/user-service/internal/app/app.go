package app

import (
	"context"
	"log"
	"net/http"
	"pLaNtS/internal/api"
	"pLaNtS/internal/auth"
	"pLaNtS/internal/config"
	"pLaNtS/internal/jwt"
	"pLaNtS/internal/repository"
	"pLaNtS/internal/security"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db     *pgxpool.Pool
	server *http.Server
}

func New(config *config.Config) (*App, error) {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, config.DatabaseURL)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewRefreshTokenRepository(db)
	hasher := security.NewBcryptHasher(config.BcryptCost)
	jwtService := jwt.NewJWTService([]byte(config.JwtSecret), config.AccessTokenTTL)
	userService := auth.NewAuthService(hasher, jwtService, tokenRepo, userRepo, config.RefreshTokenTTL)
	authHandler := api.NewAuthHandler(userService)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)
	mux.Handle("GET /me", api.AuthMiddleware(jwtService)(http.HandlerFunc(authHandler.GetMe)))

	return &App{db: db, server: &http.Server{Addr: config.Addr, Handler: mux}}, nil
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
