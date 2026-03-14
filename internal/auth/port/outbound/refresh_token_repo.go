package outbound

import "github.com/Flook0147/netflix_like/internal/auth/domain"

type RefreshTokenPort interface {
	SaveRefreshToken(username, refreshToken string) error
	DeleteRefreshToken(refreshToken string) error
	FindRefreshToken(refreshToken string) (*domain.RefreshToken, error)
}
