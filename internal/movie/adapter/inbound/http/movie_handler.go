package http

import (
	"fmt"

	"github.com/Flook0147/netflix_like/internal/movie/domain"
	"github.com/Flook0147/netflix_like/internal/movie/port/inbound"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MovieHandler struct {
	service inbound.MoviePort
}

func NewMovieHandler(service inbound.MoviePort) *MovieHandler {
	return &MovieHandler{service: service}
}

// GetMovies godoc
// @Summary Get all movies
// @Description Retrieve all available movies
// @Tags movies
// @Success 200 {array} domain.Movie
// @Failure 500 {object} map[string]string
// @Router /movies [get]
func (h *MovieHandler) GetMovies(c fiber.Ctx) error {
	movies, err := h.service.GetMovies()
	if err != nil {
		return err
	}
	return c.JSON(movies)
}

// UploadMovie godoc
// @Summary Upload movie
// @Description Upload a video file and convert to HLS
// @Tags movies
// @Accept multipart/form-data
// @Param file formData file true "Video file"
// @Success 200 {object} domain.Movie
// @Failure 400 {object} map[string]string
// @Router /movies [post]
func (h *MovieHandler) UploadMovie(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	movieID := uuid.New()
	inputPath := fmt.Sprintf("./%s.mp4", movieID.String())

	// save file
	if err := c.SaveFile(file, inputPath); err != nil {
		return err
	}

	movie := domain.Movie{
		ID:      movieID,
		HLSPath: fmt.Sprintf("/videos/%s/index.m3u8", movieID.String()),
	}

	// 🔥 ส่ง inputPath ไปให้ service
	if err := h.service.CreateMovie(&movie, inputPath); err != nil {
		return err
	}

	return c.JSON(movie)
}

func (h *MovieHandler) GetMovieStreamURL(c fiber.Ctx) (string, error) {
	viewerIDStr := c.Query("viewer_id")
	movieIDStr := c.Query("movie_id")
	return h.service.GetMovieStreamURL(uuid.MustParse(viewerIDStr), uuid.MustParse(movieIDStr))
}
