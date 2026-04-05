package http

import (
	"github.com/Flook0147/netflix_like/internal/middleware"
	"github.com/gofiber/fiber/v3"
)

func RegisterMovieRoutes(public fiber.Router, protected fiber.Router, handler *MovieHandler) {
	moviePublic := public.Group("/movies")
	moviePublic.Get("/", handler.GetMovies)

	movieProtected := protected.Group("/movies")
	movieProtected.Post("/upload", middleware.RolesAllowed("admin"), handler.UploadMovie)
	movieProtected.Get("/stream/:id", middleware.RolesAllowed("admin", "viewer"), handler.GetMovieStreamURL)
}
