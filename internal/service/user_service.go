package service

import (
	"context"

	"userrestapigorm/internal/models"
	"userrestapigorm/internal/repository"
)

type UserService struct {
	userRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) GetUsers(ctx context.Context) ([]models.User, error) {
	return s.userRepository.GetUsers(ctx)
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepository.GetUserByID(ctx, id)
}
