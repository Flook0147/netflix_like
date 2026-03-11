package app

import (
	"github.com/Flook0147/netflix_like/internal/user/domain"
	outbound "github.com/Flook0147/netflix_like/internal/user/port/outbound"
)

type UserService struct {
	repo outbound.UserRepository
}

func NewUserService(repo outbound.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(username, password, name string) error {
	return s.repo.CreateUser(username, password, name)
}

func (s *UserService) GetUser(username string) (*domain.User, error) {
	return s.repo.GetUser(username)
}
