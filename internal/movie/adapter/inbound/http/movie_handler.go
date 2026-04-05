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
// @Summary Upload a movie file
// @Description Upload a video file, process it into HLS format, and return movie metadata
// @Tags movies
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Video file (e.g., MP4)"
// @Success 200 {object} domain.Movie
// @Failure 400 {object} map[string]string "Invalid request or file upload error"
// @Failure 500 {object} map[string]string "Internal server error during processing"
// @Router /movies [post]
func (h *MovieHandler) UploadMovie(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	movieID := uuid.New()
	inputPath := fmt.Sprintf("./tmp/%s.mp4", movieID.String())

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

// GetMovieStreamURL godoc
// @Summary Get signed streaming URL
// @Description Generate a signed HLS URL for streaming if user subscription is active
// @Tags movies
// @Produce json
// @Param id path string true "Movie ID"
// @Param viewer_id query string true "Viewer ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string "Invalid input"
// @Failure 403 {object} map[string]string "Subscription inactive"
// @Failure 500 {object} map[string]string "Internal error"
// @Router /movies/stream/{id} [get]
func (h *MovieHandler) GetMovieStreamURL(c fiber.Ctx) error {
	viewerIDStr := c.Query("viewer_id")
	movieIDStr := c.Params("id") // 🔥 ใช้ param จาก route

	// validate input
	viewerID, err := uuid.Parse(viewerIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid viewer_id",
		})
	}

	movieID, err := uuid.Parse(movieIDStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid movie_id",
		})
	}

	url, err := h.service.GetMovieStreamURL(viewerID, movieID)
	if err != nil {
		return c.Status(403).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"stream_url": url,
	})
}
