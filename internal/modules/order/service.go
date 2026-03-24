package order

import (
	"context"
	"errors"
	"strings"
)

var ErrOrderNotFound = errors.New("order not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context, userID *int64) ([]Order, error) {
	return s.repo.List(ctx, userID)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*Order, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrOrderNotFound
	}

	return item, nil
}

func (s *Service) Create(ctx context.Context, input CreateOrderRequest) (*Order, error) {
	input.Item = strings.TrimSpace(input.Item)
	return s.repo.Create(ctx, input)
}
