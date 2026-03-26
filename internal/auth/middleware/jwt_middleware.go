package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("super-secret-key")

func JWTMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return fiber.ErrUnauthorized
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			return fiber.ErrUnauthorized
		}

		claims := token.Claims.(jwt.MapClaims)

		// 🔥 set userID
		c.Locals("userID", claims["user_id"])

		return c.Next()
	}
}
