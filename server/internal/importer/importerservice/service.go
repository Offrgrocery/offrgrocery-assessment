package importerservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	"offgrocery-assessment/internal/ent/store"
	"offgrocery-assessment/internal/importer/importerstore"
)

// StoreAProduct maps to a single product in Store A's JSON data feed.
type StoreAProduct struct {
	ProductName  string  `json:"product_name"`
	Manufacturer string  `json:"manufacturer"`
	RetailPrice  float64 `json:"retail_price"`
	SKU          string  `json:"sku"`
	Category     string  `json:"category"`
	WeightGrams  int     `json:"weight_grams"`
}

// StoreAData is the top-level JSON structure for Store A's data feed.
type StoreAData struct {
	StoreLocationID string          `json:"store_location_id"`
	Products        []StoreAProduct `json:"products"`
}

type Service interface {
	Import(ctx context.Context, filePath string) error
}

type service struct {
	store importerstore.Store
}

func New(store importerstore.Store) *service {
	return &service{store: store}
}

// Import reads a Store A JSON file and imports the products into the database.
func (s *service) Import(ctx context.Context, filePath string) error {
	slog.Info("importer: starting import", "file", filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	var storeData StoreAData
	if err := json.Unmarshal(data, &storeData); err != nil {
		return fmt.Errorf("parsing json: %w", err)
	}

	slog.Info("importer: parsed data", "store", storeData.StoreLocationID, "products", len(storeData.Products))

	storeRecord, err := s.store.FindOrCreateStore(ctx, storeData.StoreLocationID, store.GrocerStoreA)
	if err != nil {
		return fmt.Errorf("finding or creating store: %w", err)
	}

	for _, p := range storeData.Products {
		err := s.store.UpsertItem(ctx, p.ProductName, p.Manufacturer, p.RetailPrice, storeRecord.ID)
		if err != nil {
			return fmt.Errorf("upserting item %q: %w", p.ProductName, err)
		}
	}

	slog.Info("importer: import complete", "store", storeData.StoreLocationID, "total_products", len(storeData.Products))

	return nil
}
