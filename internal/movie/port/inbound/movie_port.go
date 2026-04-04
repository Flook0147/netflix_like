package inbound

import (
	"github.com/Flook0147/netflix_like/internal/movie/domain"
	"github.com/google/uuid"
)

type MoviePort interface {
	CreateMovie(movie *domain.Movie, path string) error
	GetMovies() ([]*domain.Movie, error)
	GetMovieByID(movieID uuid.UUID) (*domain.Movie, error)
	GetMovieStreamURL(viewerID, movieID uuid.UUID) (string, error)
}
