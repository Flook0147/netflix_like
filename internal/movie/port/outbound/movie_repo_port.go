package outbound

import (
	"github.com/Flook0147/netflix_like/internal/movie/domain"
	"github.com/google/uuid"
)

type MovieRepoPort interface {
	CreateMovie(movie *domain.Movie) error
	GetAllMovies() ([]*domain.Movie, error)
	GetMovieByID(movieID uuid.UUID) (*domain.Movie, error)
}
