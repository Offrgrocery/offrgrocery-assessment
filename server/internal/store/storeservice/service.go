package storeservice

import (
	"context"

	"offgrocery-assessment/internal/ent"
	"offgrocery-assessment/internal/store/storestore"
)

type Service interface {
	GetStoreByID(ctx context.Context, id int) (*ent.Store, error)
}

type service struct {
	store storestore.Store
}

func New(store storestore.Store) *service {
	return &service{store: store}
}

func (s *service) GetStoreByID(ctx context.Context, id int) (*ent.Store, error) {
	return s.store.GetStoreByID(ctx, id)
}
