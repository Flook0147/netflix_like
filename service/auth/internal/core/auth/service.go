package auth

import (
	"context"
	"errors"

	"github.com/Flook0147/netflix_like/service/auth/internal/domain"
	"github.com/Flook0147/netflix_like/service/auth/internal/port"
	"github.com/google/uuid"
)

var (
	ErrInvalidCredential     = errors.New("invalid credential")
	ErrUsernameAlreadyExists = errors.New("username already exists")
	ErrEmailAlreadyExists    = errors.New("email already exists")
)

type AuthService struct {
	// Add fields as necessary, e.g., database connection, config, etc.
	authProvider   port.AuthProvider
	userRepo       port.UserRepository
	passwordHasher port.PasswordHasher
}

func NewAuthService(authProvider port.AuthProvider, userRepo port.UserRepository, passwordHasher port.PasswordHasher) *AuthService {
	return &AuthService{
		authProvider:   authProvider,
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
	}
}

func (a *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	// Implement login logic here
	user, err := a.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return "", err
	}

	if err := a.passwordHasher.Compare(user.Password, password); err != nil {
		return "", ErrInvalidCredential
	}

	return a.authProvider.IssueToken(user.ID)
}

func (a *AuthService) Register(ctx context.Context, email, username, password string) error {
	// Implement registration logic here
	hashedPassword, err := a.passwordHasher.Hash(password)
	if err != nil {
		return err
	}

	if _, err := a.userRepo.FindByUsername(ctx, username); err == nil {
		return ErrUsernameAlreadyExists
	}

	if _, err := a.userRepo.FindByEmail(ctx, email); err == nil {
		return ErrEmailAlreadyExists
	}

	userId := uuid.New().String()

	user := &domain.User{
		ID:       userId,
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}
	return a.userRepo.Create(ctx, user)
}
