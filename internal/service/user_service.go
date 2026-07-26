package service

import (
	"context"
	"fmt"

	"userrestapigorm/internal/models"
	"userrestapigorm/internal/repository"

	"golang.org/x/crypto/bcrypt"
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

func (s *UserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user.Password = string(passwordHash)
	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, user *models.User) (*models.User, error) {
	return s.userRepository.UpdateUser(ctx, id, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepository.DeleteUser(ctx, id)
}
