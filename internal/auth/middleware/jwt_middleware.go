package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		panic("JWT_SECRET not set")
	}
	return []byte(secret)
}

func JWTMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		// fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))

		// fmt.Println("Token String:", tokenString)

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil {
			return err
		}

		if !token.Valid {
			return fmt.Errorf("token is invalid")
		}

		claims := token.Claims.(jwt.MapClaims)

		// 🔥 set userID
		c.Locals("userID", claims["user_id"])

		return c.Next()
	}
}
