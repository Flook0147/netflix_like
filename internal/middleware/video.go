package middleware

import (
	"strings"

	"github.com/Flook0147/netflix_like/internal/movie/utils"
	"github.com/gofiber/fiber/v3"
)

func VideoAuthMiddleware() fiber.Handler {
	return func(c fiber.Ctx) error {

		// 🔐 1. get token
		token := c.Query("token")
		if token == "" {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// 🔍 2. verify token
		claims, err := utils.VerifyVideoToken(token)
		if err != nil {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// 🎬 3. get movieID from token
		tokenMovieID, ok := claims["movie_id"].(string)
		if !ok {
			return c.SendStatus(fiber.StatusUnauthorized)
		}

		// 📂 4. get movieID from path
		// path: /videos/{movieID}/index0.ts
		path := c.Path()
		parts := strings.Split(path, "/")

		if len(parts) < 3 {
			return c.SendStatus(fiber.StatusForbidden)
		}

		pathMovieID := parts[2]

		// 🔥 5. compare
		if tokenMovieID != pathMovieID {
			return c.SendStatus(fiber.StatusForbidden)
		}

		return c.Next()
	}
}
