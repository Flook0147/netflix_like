package inbound

import "github.com/Flook0147/netflix_like/internal/user/domain"

type UserService interface {
	CreateUser(username, password, name string) error
	GetUser(username string) (*domain.User, error)
}
