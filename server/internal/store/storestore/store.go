package storestore

import (
	"context"

	"offgrocery-assessment/internal/ent"
	"offgrocery-assessment/internal/ent/store"
)

type Store interface {
	GetStoreByID(ctx context.Context, id int) (*ent.Store, error)
}

type storeStore struct {
	client *ent.Client
}

func New(client *ent.Client) *storeStore {
	return &storeStore{client: client}
}

func (s *storeStore) GetStoreByID(ctx context.Context, id int) (*ent.Store, error) {
	return s.client.Store.Query().
		Where(store.IDEQ(id)).
		WithItems().
		First(ctx)
}
