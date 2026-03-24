package user

import (
	"context"
	"errors"
	"strings"
)

var ErrUserNotFound = errors.New("user not found")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int64) (*User, error) {
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, ErrUserNotFound
	}

	return item, nil
}

func (s *Service) Create(ctx context.Context, input CreateUserRequest) (*User, error) {
	input.Name = strings.TrimSpace(input.Name)
	input.Email = strings.TrimSpace(input.Email)
	return s.repo.Create(ctx, input)
}
