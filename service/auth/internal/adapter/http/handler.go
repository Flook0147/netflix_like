package http

import (
	"encoding/json"
	"net/http"

	"github.com/Flook0147/netflix_like/service/auth/internal/core/auth"
)

type Handler struct {
	authService *auth.AuthService
}

func NewHTTPHandler(authService *auth.AuthService) *Handler {
	return &Handler{authService: authService}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(
		r.Context(),
		req.Username,
		req.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := h.authService.Register(
		r.Context(),
		req.Email,
		req.Username,
		req.Password,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
