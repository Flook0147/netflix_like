package app

import (
	"github.com/Flook0147/netflix_like/internal/user/domain"
	outbound "github.com/Flook0147/netflix_like/internal/user/port/outbound"
)

type UserService struct {
	userRepo   outbound.UserRepository
	token_port outbound.TokenPort
}

func NewUserService(userRepo outbound.UserRepository) *UserService {
	return &UserService{userRepo: userRepo, token_port: nil}
}

func (s *UserService) SetTokenPort(token_port outbound.TokenPort) {
	s.token_port = token_port
}

func (s *UserService) CreateUser(username, password, name, email string) error {
	return s.userRepo.CreateUser(username, password, name, email)
}

func (s *UserService) GetUser(username string) (*domain.User, error) {
	return s.userRepo.GetUser(username)
}

func (s *UserService) ValidateToken(token string) (string, error) {
	return s.token_port.ValidateToken(token)
}

func (s *UserService) GetUserFromToken(token string) (*domain.User, error) {
	username, err := s.token_port.ValidateToken(token)
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetUser(username)
}
