package usecase

import (
	"context"
	"errors"
	"time"

	"go-clean-arch/internal/user/domain/entity"
	"go-clean-arch/internal/user/domain/repository"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/hash"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserUseCase struct {
	userRepo repository.UserRepository
	logger   *zap.Logger
}

func NewUserUseCase(userRepo repository.UserRepository, logger *zap.Logger) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (uc *UserUseCase) CreateUser(ctx context.Context, req *entity.CreateUserRequest) (*entity.UserResponse, error) {
	// Check if user already exists
	existingUser, _ := uc.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("user already exists")
	}

	// Hash password
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		uc.logger.Error("Failed to hash password", zap.Error(err))
		return nil, errors.New("internal server error")
	}

	// Create user entity
	user := &entity.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		Name:      req.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Save to repository
	if err := uc.userRepo.Create(ctx, user); err != nil {
		uc.logger.Error("Failed to create user", zap.Error(err))
		return nil, errors.New("failed to create user")
	}

	return &entity.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (uc *UserUseCase) Login(ctx context.Context, req *entity.LoginRequest) (string, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	if !hash.CheckPasswordHash(req.Password, user.Password) {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token, err := auth.GenerateToken(user.ID.String())
	if err != nil {
		uc.logger.Error("Failed to generate token", zap.Error(err))
		return "", errors.New("internal server error")
	}

	return token, nil
}

func (uc *UserUseCase) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.UserResponse, error) {
	user, err := uc.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &entity.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
	}, nil
}
