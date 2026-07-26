package repository

import (
	"context"

	"userrestapigorm/internal/models"

	"gorm.io/gorm"
)

// UserRepository contains persistence operations for users.
type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetUsers returns all users while preserving the caller's cancellation and
// deadline through to the database query.
func (r *UserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User

	err := r.db.WithContext(ctx).Order("id DESC").Find(&users).Error
	return users, err
}
