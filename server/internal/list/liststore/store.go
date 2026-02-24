package liststore

import (
	"context"

	"offgrocery-assessment/internal/ent"
	"offgrocery-assessment/internal/ent/list"
	"offgrocery-assessment/internal/ent/user"
)

type Store interface {
	CreateList(ctx context.Context, userID int, name string) (*ent.List, error)
	GetListsByUserID(ctx context.Context, userID int) ([]*ent.List, error)
	GetListByID(ctx context.Context, id int) (*ent.List, error)
	AddItemsToList(ctx context.Context, listID int, itemIDs []int) (*ent.List, error)
}

type store struct {
	client *ent.Client
}

func New(client *ent.Client) *store {
	return &store{client: client}
}

func (s *store) CreateList(ctx context.Context, userID int, name string) (*ent.List, error) {
	return s.client.List.Create().
		SetName(name).
		SetUserID(userID).
		Save(ctx)
}

func (s *store) GetListsByUserID(ctx context.Context, userID int) ([]*ent.List, error) {
	return s.client.List.Query().
		Where(list.HasUserWith(user.ID(userID))).
		All(ctx)
}

func (s *store) GetListByID(ctx context.Context, id int) (*ent.List, error) {
	return s.client.List.Query().
		Where(list.IDEQ(id)).
		WithItems().
		First(ctx)
}

func (s *store) AddItemsToList(ctx context.Context, listID int, itemIDs []int) (*ent.List, error) {
	return s.client.List.UpdateOneID(listID).
		AddItemIDs(itemIDs...).
		Save(ctx)
}
