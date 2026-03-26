package db

import (
	"github.com/Flook0147/netflix_like/internal/movie/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MovieRepository struct {
	db *gorm.DB
}

func NewMovieRepository(db *gorm.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) CreateMovie(movie *domain.Movie) error {
	return r.db.Create(movie).Error
}

func (r *MovieRepository) GetAllMovies() ([]*domain.Movie, error) {
	var movies []*domain.Movie
	err := r.db.Find(&movies).Error
	return movies, err
}

func (r *MovieRepository) GetMovieByID(id uuid.UUID) (*domain.Movie, error) {
	var movie domain.Movie
	err := r.db.Where("id = ?", id).First(&movie).Error
	return &movie, err
}
