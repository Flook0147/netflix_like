package outbound

import "github.com/Flook0147/netflix_like/internal/user/domain"

type UserRepository interface {
	CreateUser(username, password, name, email string) error
	GetUser(username string) (*domain.User, error)
}
