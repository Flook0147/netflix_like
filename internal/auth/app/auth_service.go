package app

import (
	"fmt"

	"github.com/Flook0147/netflix_like/internal/auth/port/outbound"
	"github.com/Flook0147/netflix_like/internal/auth/utils"
	"github.com/Flook0147/netflix_like/internal/user/domain"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userPort         outbound.UserPort
	refreshTokenPort outbound.RefreshTokenPort
}

func NewAuthService(userPort outbound.UserPort, refreshTokenPort outbound.RefreshTokenPort) *AuthService {
	return &AuthService{
		userPort:         userPort,
		refreshTokenPort: refreshTokenPort,
	}
}

func (s *AuthService) Register(username, password, name, email string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.userPort.CreateUser(username, string(hashedPassword), name, email)
	if err != nil {
		return err
	}
	return nil
}

func (s *AuthService) Login(username, password string) (string, string, error) {

	user, err := s.userPort.GetUser(username)
	if err != nil {
		return "", "", err
	}

	// compare password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(user.Username)
	if err != nil {
		return "", "", err
	}

	// generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.Username)
	if err != nil {
		return "", "", err
	}

	err = s.refreshTokenPort.SaveRefreshToken(user.Username, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) ValidateToken(token string) (string, error) {
	username, err := utils.ValidateToken(token)
	if err != nil {
		return "", err
	}

	return username, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {

	_, err := s.refreshTokenPort.FindRefreshToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid refresh token")
	}

	newAccessToken, newRefreshToken, err := utils.RefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	s.refreshTokenPort.DeleteRefreshToken(refreshToken)

	username, _ := utils.ValidateRefreshToken(newRefreshToken)
	s.refreshTokenPort.SaveRefreshToken(username, newRefreshToken)

	return newAccessToken, newRefreshToken, nil
}

func (s *AuthService) GetUserFromToken(token string) (*domain.User, error) {
	username, err := utils.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	return s.userPort.GetUser(username)
}
