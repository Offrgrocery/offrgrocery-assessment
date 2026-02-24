package seed

import (
	"context"
	"fmt"
	"log/slog"

	"offgrocery-assessment/internal/ent"
)

// Seed populates the database with dev data. Assumes a clean database
// with items already imported via the import command.
func Seed(ctx context.Context, client *ent.Client) error {
	slog.Info("seed: creating users")

	alex, err := client.User.Create().SetName("Alex").SetEmail("alex@example.com").Save(ctx)
	if err != nil {
		return fmt.Errorf("seed: creating user alex: %w", err)
	}

	_, err = client.User.Create().SetName("Colin").SetEmail("colin@example.com").Save(ctx)
	if err != nil {
		return fmt.Errorf("seed: creating user colin: %w", err)
	}

	_, err = client.User.Create().SetName("Houman").SetEmail("houman@example.com").Save(ctx)
	if err != nil {
		return fmt.Errorf("seed: creating user houman: %w", err)
	}

	ray, err := client.User.Create().SetName("Ray").SetEmail("ray@example.com").Save(ctx)
	if err != nil {
		return fmt.Errorf("seed: creating user ray: %w", err)
	}

	slog.Info("seed: created users", "count", 4)

	// Create a list for Ray with some imported items.
	slog.Info("seed: creating list for ray")

	items, err := client.Item.Query().Limit(5).All(ctx)
	if err != nil {
		return fmt.Errorf("seed: querying items: %w", err)
	}

	if len(items) == 0 {
		slog.Warn("seed: no items found, run the import command first")
		return nil
	}

	itemIDs := make([]int, len(items))
	for i, item := range items {
		itemIDs[i] = item.ID
		slog.Info("seed: adding item to ray's list", "id", item.ID, "name", item.Name, "brand", item.Brand, "price", item.Price)
	}

	_, err = client.List.Create().
		SetName("Weekly Groceries").
		SetUser(ray).
		AddItemIDs(itemIDs...).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("seed: creating list for ray: %w", err)
	}

	slog.Info("seed: created list for ray", "items", len(itemIDs))

	// Create an empty list for Alex.
	_, err = client.List.Create().
		SetName("My Grocery List").
		SetUser(alex).
		Save(ctx)
	if err != nil {
		return fmt.Errorf("seed: creating list for alex: %w", err)
	}

	slog.Info("seed: created list for alex")

	slog.Info("seed: complete")

	return nil
}
