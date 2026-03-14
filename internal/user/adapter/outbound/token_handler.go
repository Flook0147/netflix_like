package outbound

import (
	"github.com/Flook0147/netflix_like/internal/auth/app"
	inbound "github.com/Flook0147/netflix_like/internal/auth/port/inbound"
)

type TokenHandler struct {
	authPort inbound.AuthPort
}

func (h TokenHandler) NewTokenHandler(authService *app.AuthService) *TokenHandler {
	return &TokenHandler{
		authPort: authService,
	}
}

func (h *TokenHandler) ValidateToken(token string) (string, error) {
	return h.authPort.ValidateToken(token)
}
