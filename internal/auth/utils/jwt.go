package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	// fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))
	return []byte(os.Getenv("JWT_SECRET"))
}

type TokenClaims struct {
	Username string
	Role     string
}

func GenerateAccessToken(username, role string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTSecret())
}

func GenerateRefreshToken(username, role string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"role":     role,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTSecret())
}

func ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role")
	}

	return &TokenClaims{
		Username: username,
		Role:     role,
	}, nil
}

func ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid role")
	}

	return &TokenClaims{
		Username: username,
		Role:     role,
	}, nil
}

func RefreshToken(refreshToken string) (string, string, error) {
	claims, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	newToken, err := GenerateAccessToken(claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(claims.Username, claims.Role)
	if err != nil {
		return "", "", err
	}

	return newToken, newRefreshToken, nil
}
