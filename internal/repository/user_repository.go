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
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, id uint, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uint) error
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

func (r *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *userRepository) UpdateUser(ctx context.Context, id uint, user *models.User) (*models.User, error) {
	result := r.db.WithContext(ctx).
		Model(&models.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"name":      user.Name,
			"email":     user.Email,
			"gender":    user.Gender,
			"birthdate": user.Birthdate,
			"is_active": user.IsActive,
		})
	if result.Error != nil {
		return nil, fmt.Errorf("update user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("user not found: %w", gorm.ErrRecordNotFound)
	}

	return r.GetUserByID(ctx, id)
}

func (r *userRepository) DeleteUser(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("delete user: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found: %w", gorm.ErrRecordNotFound)
	}

	return nil
}
