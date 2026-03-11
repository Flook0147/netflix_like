package inbound

type AuthPort interface {
	Register(username, password, name string) error
	Login(username, password string) (accessToken string, refreshToken string, err error)
	// RefreshToken(refreshToken string) (string, error)
	// ValidateToken(token string) (string, error)
}
