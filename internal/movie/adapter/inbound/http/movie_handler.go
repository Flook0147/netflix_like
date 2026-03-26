package http

import (
	"fmt"
	"os"

	"github.com/Flook0147/netflix_like/internal/movie/domain"
	"github.com/Flook0147/netflix_like/internal/movie/port/inbound"
	"github.com/Flook0147/netflix_like/internal/movie/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type MovieHandler struct {
	service inbound.MoviePort
}

func NewMovieHandler(service inbound.MoviePort) *MovieHandler {
	return &MovieHandler{service: service}
}

func (h *MovieHandler) GetMovies(c fiber.Ctx) error {
	movies, err := h.service.GetMovies()
	if err != nil {
		return err
	}
	return c.JSON(movies)
}

func (h *MovieHandler) UploadMovie(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}

	movieID := uuid.New()
	inputPath := fmt.Sprintf("./tmp/%s.mp4", movieID.String())
	outputDir := fmt.Sprintf("./videos/%s", movieID.String())

	// save file
	if err := c.SaveFile(file, inputPath); err != nil {
		return err
	}

	// create folder
	os.MkdirAll(outputDir, os.ModePerm)

	// convert to HLS
	if err := utils.ConvertToHLS(inputPath, outputDir); err != nil {
		return err
	}

	// save to DB
	movie := domain.Movie{
		ID:      movieID,
		HLSPath: fmt.Sprintf("/videos/%s/index.m3u8", movieID.String()),
	}

	err = h.service.CreateMovie(&movie)
	if err != nil {
		return err
	}

	return c.JSON(movie)
}
