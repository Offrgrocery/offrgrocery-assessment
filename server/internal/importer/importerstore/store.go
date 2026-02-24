package importerstore

import (
	"context"

	"offgrocery-assessment/internal/ent"
	"offgrocery-assessment/internal/ent/item"
	"offgrocery-assessment/internal/ent/store"
)

type Store interface {
	FindOrCreateStore(ctx context.Context, storeID string, grocer store.Grocer) (*ent.Store, error)
	UpsertItem(ctx context.Context, name string, brand string, price float64, storeID int) error
}

type importerStore struct {
	client *ent.Client
}

func New(client *ent.Client) *importerStore {
	return &importerStore{client: client}
}

func (s *importerStore) FindOrCreateStore(ctx context.Context, storeID string, grocer store.Grocer) (*ent.Store, error) {
	storeRecord, err := s.client.Store.Query().
		Where(store.StoreIDEQ(storeID)).
		Only(ctx)
	if ent.IsNotFound(err) {
		return s.client.Store.Create().
			SetStoreID(storeID).
			SetGrocer(grocer).
			Save(ctx)
	}
	return storeRecord, err
}

func (s *importerStore) UpsertItem(ctx context.Context, name string, brand string, price float64, storeID int) error {
	exists, err := s.client.Item.Query().
		Where(
			item.NameEQ(name),
			item.BrandEQ(brand),
			item.HasStoreWith(store.IDEQ(storeID)),
		).
		Exist(ctx)
	if err != nil {
		return err
	}

	if exists {
		_, err = s.client.Item.Update().
			Where(
				item.NameEQ(name),
				item.BrandEQ(brand),
				item.HasStoreWith(store.IDEQ(storeID)),
			).
			SetPrice(price).
			Save(ctx)
		return err
	}

	_, err = s.client.Item.Create().
		SetName(name).
		SetBrand(brand).
		SetPrice(price).
		SetStoreID(storeID).
		Save(ctx)
	return err
}
