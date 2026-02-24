package app

import (
	"context"
	"log/slog"

	"offgrocery-assessment/internal/config"
	"offgrocery-assessment/internal/seed"
)

func NewSeed(cfg config.Config) error {
	ctx := context.Background()

	slog.Info("seed: connecting to database")
	db, err := NewDB(ctx, cfg, ConfigureMySQLParseTime)
	if err != nil {
		return err
	}
	defer db.Close()

	slog.Info("seed: pinging database")
	if err := db.PingContext(ctx); err != nil {
		slog.Error("seed: failed to ping database", "error", err)
		return err
	}

	slog.Info("seed: creating ent client")
	client := NewEntClient(db)

	slog.Info("seed: running auto migration")
	if err := client.Schema.Create(ctx); err != nil {
		slog.Error("seed: failed to run auto migration", "error", err)
		return err
	}

	return seed.Seed(ctx, client)
}
