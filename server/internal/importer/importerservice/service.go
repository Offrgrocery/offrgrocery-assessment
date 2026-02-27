package importerservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"offgrocery-assessment/internal/ent/store"
	"offgrocery-assessment/internal/importer/importerstore"
)

// =====================================================================
// 1. STANDARD MODELS & INTERFACES
// =====================================================================

// NormalizedProduct is our clean, internal representation of any product.
type NormalizedProduct struct {
	Name         string
	Manufacturer string
	Price        float64
	SKU          string
}

// StoreParser defines the contract for parsing any store's JSON feed.
type StoreParser interface {
	StoreType() store.Grocer
	Parse(file *os.File) (storeLocationID string, products []NormalizedProduct, err error)
}

// =====================================================================
// 2. STORE ADAPTERS (A, B, C)
// =====================================================================

// --- Store A ---
type StoreAData struct {
	StoreLocationID string `json:"store_location_id"`
	Products        []struct {
		ProductName  string  `json:"product_name"`
		Manufacturer string  `json:"manufacturer"`
		RetailPrice  float64 `json:"retail_price"`
		SKU          string  `json:"sku"`
	} `json:"products"`
}

type StoreAParser struct{}

func (p *StoreAParser) StoreType() store.Grocer { return store.GrocerStoreA }
func (p *StoreAParser) Parse(file *os.File) (string, []NormalizedProduct, error) {
	var data StoreAData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return "", nil, err
	}

	var products []NormalizedProduct
	for _, prod := range data.Products {
		products = append(products, NormalizedProduct{
			Name:         prod.ProductName,
			Manufacturer: prod.Manufacturer,
			Price:        prod.RetailPrice,
			SKU:          prod.SKU,
		})
	}
	return data.StoreLocationID, products, nil
}

// --- Store B ---
type StoreBData struct {
	Location struct {
		ID string `json:"id"`
	} `json:"location"`
	Inventory []struct {
		Item struct {
			Label     string `json:"label"`
			BrandName string `json:"brand_name"`
		} `json:"item"`
		Pricing struct {
			CurrentPrice float64 `json:"current_price"`
		} `json:"pricing"`
		Barcode string `json:"barcode"`
	} `json:"inventory"`
}

type StoreBParser struct{}

// Assuming GrocerStoreB exists in their enum based on GrocerStoreA
func (p *StoreBParser) StoreType() store.Grocer { return store.GrocerStoreB }
func (p *StoreBParser) Parse(file *os.File) (string, []NormalizedProduct, error) {
	var data StoreBData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return "", nil, err
	}

	var products []NormalizedProduct
	for _, inv := range data.Inventory {
		products = append(products, NormalizedProduct{
			Name:         inv.Item.Label,
			Manufacturer: inv.Item.BrandName,
			Price:        inv.Pricing.CurrentPrice,
			SKU:          inv.Barcode,
		})
	}
	return data.Location.ID, products, nil
}

// --- Store C ---
type StoreCData struct {
	StoreCode string `json:"store_code"`
	Catalogue []struct {
		DisplayName string  `json:"display_name"`
		Producer    string  `json:"producer"`
		Cost        float64 `json:"cost"`
		ProductID   string  `json:"product_id"`
	} `json:"catalogue"`
}

type StoreCParser struct{}

// Assuming GrocerStoreC exists in their enum
func (p *StoreCParser) StoreType() store.Grocer { return store.GrocerStoreC }
func (p *StoreCParser) Parse(file *os.File) (string, []NormalizedProduct, error) {
	var data StoreCData
	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return "", nil, err
	}

	var products []NormalizedProduct
	for _, cat := range data.Catalogue {
		products = append(products, NormalizedProduct{
			Name:         cat.DisplayName,
			Manufacturer: cat.Producer,
			Price:        cat.Cost,
			SKU:          cat.ProductID,
		})
	}
	return data.StoreCode, products, nil
}

// =====================================================================
// 3. FACTORY
// =====================================================================

// getParserForFile determines which adapter to use based on the filename.
func getParserForFile(filePath string) (StoreParser, error) {
	lowerPath := strings.ToLower(filePath)
	if strings.Contains(lowerPath, "store_a") {
		return &StoreAParser{}, nil
	} else if strings.Contains(lowerPath, "store_b") {
		return &StoreBParser{}, nil
	} else if strings.Contains(lowerPath, "store_c") {
		return &StoreCParser{}, nil
	}
	return nil, fmt.Errorf("unsupported store format for file: %s", filePath)
}

// =====================================================================
// 4. CORE SERVICE
// =====================================================================

type Service interface {
	Import(ctx context.Context, filePath string) error
}

type service struct {
	store importerstore.Store
}

func New(store importerstore.Store) *service {
	return &service{store: store}
}

// Import reads a Store JSON file, normalizes it, and imports products into the database.
func (s *service) Import(ctx context.Context, filePath string) error {
	slog.Info("importer: starting import", "file", filePath)

	// 1. Route to the correct parser
	parser, err := getParserForFile(filePath)
	if err != nil {
		return err
	}

	// 2. Open file (Fixing the OOM memory issue by using a decoder inside Parse)
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// 3. Parse and Normalize the data
	storeLocationID, products, err := parser.Parse(file)
	if err != nil {
		return fmt.Errorf("parsing store data: %w", err)
	}

	slog.Info("importer: parsed data", "store", storeLocationID, "products", len(products))

	// 4. Database Operations
	storeRecord, err := s.store.FindOrCreateStore(ctx, storeLocationID, parser.StoreType())
	if err != nil {
		return fmt.Errorf("finding or creating store: %w", err)
	}

	var errorCount int
	for _, p := range products {
		// Context check for graceful cancellation
		if ctx.Err() != nil {
			return fmt.Errorf("import aborted: %w", ctx.Err())
		}

		err := s.store.UpsertItem(ctx, p.Name, p.Manufacturer, p.Price, storeRecord.ID)
		if err != nil {
			slog.Error("failed to upsert item", "product", p.Name, "error", err)
			errorCount++
			continue // Don't fail the whole batch for one bad item
		}
	}

	slog.Info("importer: import complete",
		"store", storeLocationID,
		"total_products", len(products),
		"failed_upserts", errorCount,
	)

	return nil
}
