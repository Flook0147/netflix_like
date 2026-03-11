package outbound

import "github.com/Flook0147/netflix_like/internal/user/domain"

type User struct {
	Username string
	Password string
	Name     string
}

type UserPort interface {
	CreateUser(username, password, name string) error
	GetUser(username string) (*domain.User, error)
}
