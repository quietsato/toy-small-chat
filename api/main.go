package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/quietsato/toy-small-chat/api/internal/config"
	"github.com/quietsato/toy-small-chat/api/internal/di"
	"github.com/quietsato/toy-small-chat/api/internal/instrument"
	instrumentdb "github.com/quietsato/toy-small-chat/api/internal/instrument/db"
	instrumenthttp "github.com/quietsato/toy-small-chat/api/internal/instrument/http"
	"github.com/quietsato/toy-small-chat/api/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, os.Interrupt)
	defer stop()

	cfg := config.Load()

	// Initialize OpenTelemetry tracer
	shutdownInstr := instrument.Init(ctx, slog.LevelInfo, cfg.OtlpEndpoint)
	defer shutdownInstr()

	// Initialize database connection pool with tracing
	pool, err := instrumentdb.NewPool(ctx, cfg.Database.URL())
	if err != nil {
		slog.Error("failed to connect to database", slog.Any("err", err))
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		slog.Error("failed to ping database", slog.Any("err", err))
		return
	}
	slog.Info("successfully connected to database")

	// Create router and wrap with HTTP tracing
	router := server.New(di.New(pool, cfg.JWTSecretKey))
	handler := instrumenthttp.NewHandler(router, "toy-small-chat")

	srv := http.Server{
		Addr:    ":8080",
		Handler: handler,
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			slog.Error("failed to serve", slog.Any("err", err))
		}
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.ErrorContext(ctx, "failed to shutdown", slog.Any("err", err))
		return
	}

	slog.InfoContext(ctx, "server stopped gracefully")
}
