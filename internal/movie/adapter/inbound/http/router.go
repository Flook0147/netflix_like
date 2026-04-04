package http

import "github.com/gofiber/fiber/v3"

func RegisterMovieRoutes(public fiber.Router, protected fiber.Router, handler *MovieHandler) {
	moviePublic := public.Group("/movies")
	moviePublic.Get("/", handler.GetMovies)

	movieProtected := protected.Group("/movies")
	movieProtected.Post("/upload", handler.UploadMovie)
	movieProtected.Get("/stream/:id", handler.GetMovies)
}
