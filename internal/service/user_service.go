package service

import (
	"context"

	"userrestapigorm/internal/models"
	"userrestapigorm/internal/repository"
)

// UserService coordinates user-related application operations.
type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepository.GetUsers(ctx)
}
