package service

import (
	"context"

	"userrestapigorm/internal/models"
)

// UserRepository defines the persistence operation required by UserService.
// Keeping this interface at the consumer makes the service straightforward to
// test with a repository mock.
type UserRepository interface {
	GetUsers(ctx context.Context) ([]models.User, error)
}

// UserService coordinates user-related application operations.
type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepository UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepository.GetUsers(ctx)
}
