package listservice

import (
	"context"

	"offgrocery-assessment/internal/ent"
	"offgrocery-assessment/internal/list/liststore"
)

type Service interface {
	CreateList(ctx context.Context, userID int, name string) (*ent.List, error)
	GetListsByUserID(ctx context.Context, userID int) ([]*ent.List, error)
	GetListByID(ctx context.Context, id int) (*ent.List, error)
	AddItemsToList(ctx context.Context, listID int, itemIDs []int) (*ent.List, error)
	DeleteList(ctx context.Context, id int) error
	RemoveItemsFromList(ctx context.Context, listID int, itemIDs []int) (*ent.List, error)
}

type service struct {
	store liststore.Store
}

func New(store liststore.Store) *service {
	return &service{store: store}
}

func (s *service) CreateList(ctx context.Context, userID int, name string) (*ent.List, error) {
	return s.store.CreateList(ctx, userID, name)
}

func (s *service) GetListsByUserID(ctx context.Context, userID int) ([]*ent.List, error) {
	return s.store.GetListsByUserID(ctx, userID)
}

func (s *service) GetListByID(ctx context.Context, id int) (*ent.List, error) {
	return s.store.GetListByID(ctx, id)
}

func (s *service) AddItemsToList(ctx context.Context, listID int, itemIDs []int) (*ent.List, error) {
	return s.store.AddItemsToList(ctx, listID, itemIDs)
}

func (s *service) DeleteList(ctx context.Context, id int) error {
	return s.store.DeleteList(ctx, id)
}

func (s *service) RemoveItemsFromList(ctx context.Context, listID int, itemIDs []int) (*ent.List, error) {
	return s.store.RemoveItemsFromList(ctx, listID, itemIDs)
}
