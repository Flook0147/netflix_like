package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET"))
}

func GenerateAccessToken(username string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTSecret())
}

func GenerateRefreshToken(username string) (string, error) {

	claims := jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(getJWTSecret())
}

func ValidateToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token username")
	}
	return username, nil
}

func ValidateRefreshToken(tokenString string) (string, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return getJWTSecret(), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}
	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("invalid token username")
	}
	return username, nil

}

func RefreshToken(refreshToken string) (string, string, error) {
	username, err := ValidateRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	newToken, err := GenerateAccessToken(username)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(username)
	if err != nil {
		return "", "", err
	}

	return newToken, newRefreshToken, nil
}
