package authstore

import (
	"context"
	"offgrocery-assessment/internal/ent"
)

type Store interface {
	CreateUser(ctx context.Context, email, name string) (*ent.User, error)
}

type store struct {
	client *ent.Client
}

func New(client *ent.Client) *store {
	return &store{
		client: client,
	}
}

func (s *store) CreateUser(ctx context.Context, email, name string) (*ent.User, error) {
	return s.client.User.Create().SetEmail(email).SetName(name).Save(ctx)
}
