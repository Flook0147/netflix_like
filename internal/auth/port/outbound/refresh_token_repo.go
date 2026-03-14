package outbound

type RefreshTokenPort interface {
	SaveRefreshToken(username, refreshToken string) error
	DeleteRefreshToken(refreshToken string) error
}
