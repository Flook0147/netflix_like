package inbound

import (
	"github.com/Flook0147/netflix_like/internal/auth/utils"
	"github.com/Flook0147/netflix_like/internal/user/domain"
)

type AuthPort interface {
	Register(username, password, name, email string) error
	Login(username, password string) (accessToken string, refreshToken string, err error)
	RefreshToken(refreshToken string) (string, string, error)
	ValidateToken(token string) (*utils.TokenClaims, error)
	GetUserFromToken(token string) (*domain.User, error)
}
