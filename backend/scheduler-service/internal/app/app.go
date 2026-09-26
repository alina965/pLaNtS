package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/alina965/pLaNtS/scheduler-service/internal/client"
	"github.com/alina965/pLaNtS/scheduler-service/internal/config"
	"github.com/alina965/pLaNtS/scheduler-service/internal/scheduler"
	"github.com/robfig/cron/v3"
)

type App struct {
	cron   *cron.Cron
	server *http.Server
}

func New(cfg *config.Config) (*App, error) {
	userPlantsClient := client.NewUserPlantsClient(cfg.UserPlantsURL, cfg.Timeout, cfg.InternalAPIKey)
	schedulerService := scheduler.NewService(cfg.Location, userPlantsClient)

	c, err := schedulerService.Setup()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return &App{
		cron: c,
		server: &http.Server{
			Addr:    cfg.Addr,
			Handler: mux,
		},
	}, nil
}

func (a *App) Run() error {
	a.cron.Start()
	log.Printf("scheduler started; health on %s", a.server.Addr)

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
	cronCtx := a.cron.Stop()
	select {
	case <-cronCtx.Done():
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return a.server.Shutdown(shutdownCtx)
}
