package repository

import (
	"context"
	"fmt"

	"userrestapigorm/internal/models"

	"gorm.io/gorm"
)

// 1. Interface definition
type UserRepository interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
}

// 2. Concrete struct implementation
type userRepository struct {
	db *gorm.DB
}

// 3. Constructor
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// 4. Method implementation
func (r *userRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.
		WithContext(ctx).
		Order("id DESC").
		Find(&users).
		Error

	if err != nil {
		return nil, fmt.Errorf("find users: %w", err)
	}

	return users, nil
}

func (r *userRepository) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	err := r.db.
		WithContext(ctx).
		First(&user, id).
		Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return nil, fmt.Errorf("find user by ID: %w", err)
	}

	return &user, nil
}
