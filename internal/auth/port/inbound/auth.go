package inbound

type AuthPort interface {
	Register(username, password, name, email string) error
	Login(username, password string) (accessToken string, refreshToken string, err error)
	RefreshToken(refreshToken string) (string, string, error)
	ValidateToken(token string) (string, error)
}
