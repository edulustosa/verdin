package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/edulustosa/verdin/config"
	"github.com/edulustosa/verdin/internal/api"
	"github.com/edulustosa/verdin/internal/api/router"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(
		ctx,
		os.Interrupt,
		os.Kill,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)
	defer cancel()

	if err := run(ctx); err != nil {
		slog.Error(err.Error())
		cancel()

		os.Exit(1)
	}

	slog.Info("To infinity and beyond!")
}

func run(ctx context.Context) error {
	env, err := config.LoadEnv(".")
	if err != nil {
		return err
	}

	pool, err := pgxpool.New(ctx, env.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	if err := runMigrations(ctx, pool); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	r := router.NewServer(&api.API{
		Database: pool,
		JWTKey:   env.JWTSecret,
	})
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", env.Port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 12 * time.Second,
	}
	defer func() {
		const timeout = 30 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown server", "error", err)
		}
	}()

	errChan := make(chan error, 1)
	go func() {
		slog.Info("starting server", "port", env.Port)
		errChan <- srv.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			return fmt.Errorf("server error: %w", err)
		}
	}

	return nil
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		os.DirFS("./internal/database/migrations"),
	)
	if err != nil {
		return fmt.Errorf("failed to create goose provider: %w", err)
	}

	_, err = provider.Up(ctx)
	return err
}
