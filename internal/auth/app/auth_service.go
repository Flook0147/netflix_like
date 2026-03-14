package app

import (
	"fmt"
	"os"

	"github.com/Flook0147/netflix_like/internal/auth/port/outbound"
	"github.com/Flook0147/netflix_like/internal/auth/utils"
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

	return accessToken, refreshToken, nil
}

func (s *AuthService) ValidateToken(token string) (string, error) {
	fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))
	username, err := utils.ValidateToken(token)
	if err != nil {
		return "", err
	}

	return username, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (string, string, error) {
	token, refreshToken, err := utils.RefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}
