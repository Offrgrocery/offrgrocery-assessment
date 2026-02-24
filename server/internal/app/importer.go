package app

import (
	"context"
	"log/slog"

	"offgrocery-assessment/internal/config"
	"offgrocery-assessment/internal/importer/importerservice"
	"offgrocery-assessment/internal/importer/importerstore"
)

func NewImporter(cfg config.Config, filePath string) error {
	ctx := context.Background()

	slog.Info("importer: connecting to database")
	db, err := NewDB(ctx, cfg, ConfigureMySQLParseTime)
	if err != nil {
		return err
	}
	defer db.Close()

	slog.Info("importer: pinging database")
	if err := db.PingContext(ctx); err != nil {
		slog.Error("seed: failed to ping database", "error", err)
		return err
	}

	slog.Info("importer: creating ent client")
	client := NewEntClient(db)

	slog.Info("importer: running auto migration")
	if err := client.Schema.Create(ctx); err != nil {
		slog.Error("importer: failed to run auto migration", "error", err)
		return err
	}

	store := importerstore.New(client)
	service := importerservice.New(store)

	return service.Import(ctx, filePath)
}
