package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/pressly/goose/v3/lock"
	"github.com/reefops/reefops/services/authorizer/internal/migrations"
)

const databaseURLVariable = "REEFOPS_AUTHORIZER_MIGRATOR_DATABASE_URL"

func main() {
	if err := run(); err != nil {
		slog.Error("authorizer audit migration failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	databaseURL := os.Getenv(databaseURLVariable)
	if databaseURL == "" {
		return fmt.Errorf("%s is required", databaseURLVariable)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer func() {
		if closeErr := db.Close(); closeErr != nil {
			slog.Error("close migration database", "error", closeErr)
		}
	}()

	pingCtx, cancelPing := context.WithTimeout(ctx, 10*time.Second)
	defer cancelPing()
	if err := db.PingContext(pingCtx); err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	sessionLocker, err := lock.NewPostgresSessionLocker(lock.WithLockTimeout(2, 30))
	if err != nil {
		return fmt.Errorf("create PostgreSQL session locker: %w", err)
	}
	migrationFiles, err := fs.Sub(migrations.Files, "files")
	if err != nil {
		return fmt.Errorf("open embedded migrations: %w", err)
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		migrationFiles,
		goose.WithSessionLocker(sessionLocker),
	)
	if err != nil {
		return fmt.Errorf("create migration provider: %w", err)
	}

	results, err := provider.Up(ctx)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return fmt.Errorf("migration interrupted: %w", err)
		}
		return fmt.Errorf("apply migrations: %w", err)
	}

	for _, result := range results {
		slog.Info("authorizer audit migration applied", "migration", result.String())
	}
	if len(results) == 0 {
		slog.Info("authorizer audit schema is current")
	}

	return nil
}
