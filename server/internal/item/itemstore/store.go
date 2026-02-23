package itemstore

import (
	"context"

	"offgrocery-assessment/internal/ent"

	"entgo.io/ent/dialect/sql"
)

type Store interface {
	SearchWithLimit(ctx context.Context, query string, limit int) ([]*ent.Item, error)
}

type store struct {
	client *ent.Client
}

func New(client *ent.Client) *store {
	return &store{client: client}
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
		Limit(limit).
		All(ctx)
}
