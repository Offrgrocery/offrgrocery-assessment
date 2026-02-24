package itemstore

import (
	"context"

	"offgrocery-assessment/internal/ent"
	"offgrocery-assessment/internal/ent/item"

	"entgo.io/ent/dialect/sql"
)

type Store interface {
	GetItemByID(ctx context.Context, id int) (*ent.Item, error)
	SearchWithLimit(ctx context.Context, query string, limit int) ([]*ent.Item, error)
}

type store struct {
	client *ent.Client
}

func New(client *ent.Client) *store {
	return &store{client: client}
}

func (s *store) GetItemByID(ctx context.Context, id int) (*ent.Item, error) {
	return s.client.Item.Query().
		Where(item.IDEQ(id)).
		WithStore().
		First(ctx)
}

func (s *store) SearchWithLimit(ctx context.Context, query string, limit int) ([]*ent.Item, error) {
	return s.client.Item.Query().
		Where(func(sel *sql.Selector) {
			sel.Where(
				sql.ExprP(
					"MATCH(name, brand) AGAINST(? IN BOOLEAN MODE)",
					"+"+query+"*",
				),
			)
		}).
		WithStore().
		Limit(limit).
		All(ctx)
}
