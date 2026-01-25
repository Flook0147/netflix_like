// internal/auth/adapter/jwt/provider.go
package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Provider struct {
	secret string
	ttl    time.Duration
}

func New(secret string, ttl time.Duration) *Provider {
	return &Provider{secret: secret, ttl: ttl}
}

func (p *Provider) IssueToken(userID string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(p.ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secret))
}

func (p *Provider) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(p.secret), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	return claims["sub"].(string), nil
}
