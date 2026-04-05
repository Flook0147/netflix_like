package middleware

import "github.com/gofiber/fiber/v3"

func RolesAllowed(allowedRoles ...string) fiber.Handler {
	return func(c fiber.Ctx) error {

		user := c.Locals("user")
		if user == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		claims, ok := user.(map[string]interface{})
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token claims",
			})
		}

		// 🔥 handle type safely
		roleValue, exists := claims["role"]
		if !exists {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "role not found",
			})
		}

		role, ok := roleValue.(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid role type",
			})
		}

		for _, allowed := range allowedRoles {
			if role == allowed {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "forbidden",
		})
	}
}
