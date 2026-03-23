package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/I-Van-Radkov/vesta-gkh/internal/config"
	v1 "github.com/I-Van-Radkov/vesta-gkh/internal/controller/http/v1"
	postgres "github.com/I-Van-Radkov/vesta-gkh/pkg/db"
)

type App struct {
	httpServer *v1.Server
	postgresDB *postgres.Database
}

func RunApp(cfg *config.Config) {
	ctx := context.Background()

	db, err := postgres.New(cfg.PostgresConfig)
	if err != nil {
		panic(fmt.Errorf("failed to connect to database: %w", err))
	}

	server := v1.NewServer(cfg)
	err = server.RegisterHandlers(ctx, cfg, db.Pool)
	if err != nil {
		panic(fmt.Errorf("failed to register handlers: %w", err))
	}

	app := &App{
		httpServer: server,
		postgresDB: db,
	}

	app.MustRun(ctx, cfg.GHTimeout)
}

func (a *App) MustRun(ctx context.Context, timeout time.Duration) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.httpServer.Start(); !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	graceSh := make(chan os.Signal, 1)
	signal.Notify(graceSh, os.Interrupt, syscall.SIGTERM)
	<-graceSh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := a.httpServer.Stop(shutdownCtx); err != nil {
		panic(err)
	}

	a.postgresDB.Close()

	wg.Wait()
}
