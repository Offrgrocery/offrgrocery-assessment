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
	"offgrocery-assessment/internal/item/itemhandler"
	"offgrocery-assessment/internal/item/itemservice"
	"offgrocery-assessment/internal/item/itemstore"
	"offgrocery-assessment/internal/list/listhandler"
	"offgrocery-assessment/internal/list/listservice"
	"offgrocery-assessment/internal/list/liststore"
	"offgrocery-assessment/internal/store/storehandler"
	"offgrocery-assessment/internal/store/storeservice"
	"offgrocery-assessment/internal/store/storestore"
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

	authStore := authstore.New(client)
	authService := authservice.New(authStore)
	authHandler := authhandler.New(authService)

	listStore := liststore.New(client)
	listService := listservice.New(listStore)
	listHandler := listhandler.New(listService)

	itemStore := itemstore.New(client)
	itemService := itemservice.New(itemStore)
	itemHandler := itemhandler.New(itemService)

	stStore := storestore.New(client)
	stService := storeservice.New(stStore)
	stHandler := storehandler.New(stService)

	r := chi.NewRouter()
	r.Mount("/auth", authHandler.Routes())
	r.Mount("/lists", listHandler.Routes())
	r.Mount("/items", itemHandler.Routes())
	r.Mount("/stores", stHandler.Routes())

	slog.Info("web: starting server", "port", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("web: server failed", "error", err)
		return err
	}

	return nil
}
