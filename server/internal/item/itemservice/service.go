package itemservice

import (
	"context"

	"offgrocery-assessment/internal/ent"
	"offgrocery-assessment/internal/item/itemstore"
)

type Service interface {
	GetItemByID(ctx context.Context, id int) (*ent.Item, error)
	SearchWithLimit(ctx context.Context, query string, limit int) ([]*ent.Item, error)
}

type service struct {
	store itemstore.Store
}

func New(store itemstore.Store) *service {
	return &service{store: store}
}

func (s *service) GetItemByID(ctx context.Context, id int) (*ent.Item, error) {
	return s.store.GetItemByID(ctx, id)
}

func (s *service) SearchWithLimit(ctx context.Context, query string, limit int) ([]*ent.Item, error) {
	return s.store.SearchWithLimit(ctx, query, limit)
}
