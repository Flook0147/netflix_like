package app

import (
	"github.com/Flook0147/netflix_like/internal/user/domain"
	outbound "github.com/Flook0147/netflix_like/internal/user/port/outbound"
)

type UserService struct {
	userRepo outbound.UserRepository
}

func NewUserService(userRepo outbound.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(username, password, name, email string) error {
	return s.userRepo.CreateUser(username, password, name, email)
}

func (s *UserService) GetUser(username string) (*domain.User, error) {
	return s.userRepo.GetUser(username)
}
