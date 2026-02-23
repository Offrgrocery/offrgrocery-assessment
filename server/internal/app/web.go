package app

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"offgrocery-assessment/internal/auth/authhandler"
	"offgrocery-assessment/internal/auth/authservice"
	"offgrocery-assessment/internal/auth/authstore"
	"offgrocery-assessment/internal/config"
)

func NewWeb(cfg config.Config) error {
	ctx := context.Background()

	slog.Info("web: connecting to database")
	db, err := NewDB(ctx, cfg, ConfigureMySQLParseTime)
	if err != nil {
		return err
	}
	defer db.Close()

	slog.Info("web: pinging database")
	if err := db.PingContext(ctx); err != nil {
		slog.Error("web: failed to ping database", "error", err)
		return err
	}

	slog.Info("web: creating ent client")
	client := NewEntClient(db)

	slog.Info("web: running auto migration")
	if err := client.Schema.Create(ctx); err != nil {
		slog.Error("web: failed to run auto migration", "error", err)
		return err
	}

	store := authstore.New(client)
	service := authservice.New(store)
	handler := authhandler.New(service)

	r := chi.NewRouter()
	r.Mount("/auth", handler.Routes())

	slog.Info("web: starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("web: server failed", "error", err)
		return err
	}

	return nil
}
