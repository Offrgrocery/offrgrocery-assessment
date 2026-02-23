package authservice

import (
	"context"
	"offgrocery-assessment/internal/auth/authstore"
)

type Service interface {
	CreateUser(ctx context.Context, email, name string) (int, error)
}

type service struct {
	store authstore.Store
}

func New(store authstore.Store) *service {
	return &service{
		store: store,
	}
}

func (s *service) CreateUser(ctx context.Context, email, name string) (int, error) {
	user, err := s.store.CreateUser(ctx, email, name)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}
