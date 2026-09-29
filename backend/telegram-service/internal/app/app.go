package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alina965/pLaNtS/telegram-service/internal/api"
	"github.com/alina965/pLaNtS/telegram-service/internal/client"
	"github.com/alina965/pLaNtS/telegram-service/internal/config"
	"github.com/alina965/pLaNtS/telegram-service/internal/jwt"
	"github.com/alina965/pLaNtS/telegram-service/internal/repository"
	"github.com/alina965/pLaNtS/telegram-service/internal/telegram"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	db         *pgxpool.Pool
	server     *http.Server
	cancelPoll context.CancelFunc
}

func New(cfg *config.Config) (*App, error) {
	ctx := context.Background()
	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	repo := repository.NewTelegramLinksRepository(db)
	bot := client.NewBotClient(cfg.Token, cfg.Timeout)
	service := telegram.NewService(repo, bot, cfg.BotUsername)
	handler := api.NewTelegramHandler(service)
	jwtService := jwt.NewService([]byte(cfg.JwtSecret))

	auth := api.AuthMiddleware(jwtService)
	internal := api.InternalAPIMiddleware(cfg.InternalAPIKey)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("POST /telegram/link", auth(http.HandlerFunc(handler.CreateLink)))
	mux.Handle("POST /notify", internal(http.HandlerFunc(handler.Notify)))

	pollCtx, cancelPoll := context.WithCancel(context.Background())
	go service.StartPolling(pollCtx)

	return &App{
		db: db,
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
		},
		cancelPoll: cancelPoll,
	}, nil
}

func (a *App) Run() error {
	log.Printf("telegram-service started on %s", a.server.Addr)

	errCh := make(chan error, 1)
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case <-stop:
		log.Println("shutdown signal received")
		return nil
	}
}

func (a *App) Shutdown(ctx context.Context) error {
	if a.cancelPoll != nil {
		a.cancelPoll()
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(shutdownCtx)
	a.db.Close()
	return err
}
