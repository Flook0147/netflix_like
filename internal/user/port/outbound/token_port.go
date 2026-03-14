package outbound

type TokenPort interface {
	ValidateToken(token string) (string, error)
}
