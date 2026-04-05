package middleware

import (
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
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		// ✅ safe split
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token format",
			})
		}

		tokenString := parts[1]

		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			return getJWTSecret(), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "token not valid",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid claims",
			})
		}

		// 🔥 สำคัญ: set user ให้ RolesMiddleware ใช้
		c.Locals("user", map[string]interface{}(claims))

		// optional
		c.Locals("userID", claims["user_id"])

		return c.Next()
	}
}
