package repository

import (
	"context"
	"errors"

	"go_web/internal/model"
	"go_web/internal/txmanager"

	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	List(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id string) error
	EmailExists(ctx context.Context, email string, excludeID string) (bool, error)
}

type userRepository struct {
	txManager txmanager.Manager
}

func NewUserRepository(txManager txmanager.Manager) UserRepository {
	return &userRepository{txManager: txManager}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	return r.txManager.DB(ctx).Create(user).Error
}

func (r *userRepository) List(ctx context.Context) ([]model.User, error) {
	var users []model.User
	err := r.txManager.DB(ctx).Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.txManager.DB(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	result := r.txManager.DB(ctx).Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]any{
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
	})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := r.txManager.DB(ctx).Delete(&model.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string, excludeID string) (bool, error) {
	query := r.txManager.DB(ctx).Model(&model.User{}).Where("email = ?", email)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
