package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"go_web/internal/model"
	"go_web/internal/repository"
	"go_web/internal/txmanager"
)

var (
	ErrInvalidName  = errors.New("name is required")
	ErrInvalidEmail = errors.New("email is required")
	ErrInvalidAge   = errors.New("age must be between 1 and 150")
	ErrEmailTaken   = errors.New("email already exists")
)

type CreateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
}

type UpdateUserInput struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required"`
}

type UserService interface {
	CreateUser(ctx context.Context, input CreateUserInput) (*model.User, error)
	ListUsers(ctx context.Context) ([]model.User, error)
	GetUser(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type userService struct {
	repo      repository.UserRepository
	txManager txmanager.Manager
}

func NewUserService(repo repository.UserRepository, txManager txmanager.Manager) UserService {
	return &userService{repo: repo, txManager: txManager}
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserInput) (*model.User, error) {
	user := &model.User{
		ID:    uuid.NewString(),
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(strings.ToLower(input.Email)),
		Age:   input.Age,
	}
	if err := validateUserPayload(user.Name, user.Email, user.Age); err != nil {
		return nil, err
	}

	if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		exists, err := s.repo.EmailExists(txCtx, user.Email, "")
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailTaken
		}
		return s.repo.Create(txCtx, user)
	}); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ListUsers(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, id string, input UpdateUserInput) (*model.User, error) {
	updatedUser := &model.User{
		ID:    id,
		Name:  strings.TrimSpace(input.Name),
		Email: strings.TrimSpace(strings.ToLower(input.Email)),
		Age:   input.Age,
	}
	if err := validateUserPayload(updatedUser.Name, updatedUser.Email, updatedUser.Age); err != nil {
		return nil, err
	}

	//这里就是调用事务处理器开启事务，入参 ctx上下文和匿名函数的实现func(txCtx context.Context) ，匿名函数实现内容主要是sql操作。
	/这里要注意为什么dao层的函数不能直接用ctx，因为ctx可能携带也可能不携带事务DB，这个获取逻辑是写在Manager.WithinTransaction的声明过程中的,这个过程里帮回调函数的入参获取到携带事务的上下文。
	if err := s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := s.repo.GetByID(txCtx, id); err != nil {
			return err
		}

		exists, err := s.repo.EmailExists(txCtx, updatedUser.Email, id)
		if err != nil {
			return err
		}
		if exists {
			return ErrEmailTaken
		}

		return s.repo.Update(txCtx, updatedUser)
	}); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, id)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.txManager.WithinTransaction(ctx, func(txCtx context.Context) error {
		if _, err := s.repo.GetByID(txCtx, id); err != nil {
			return err
		}
		return s.repo.Delete(txCtx, id)
	})
}

func validateUserPayload(name, email string, age int) error {
	if name == "" {
		return ErrInvalidName
	}
	if email == "" {
		return ErrInvalidEmail
	}
	if age < 1 || age > 150 {
		return fmt.Errorf("%w: %d", ErrInvalidAge, age)
	}
	return nil
}
