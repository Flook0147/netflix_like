package app

import (
	"fmt"

	"github.com/Flook0147/netflix_like/internal/movie/domain"
	"github.com/Flook0147/netflix_like/internal/movie/port/outbound"
	"github.com/Flook0147/netflix_like/internal/movie/utils"
	"github.com/google/uuid"
)

const (
	StatusPending  = "PENDING"
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
	StatusFailed   = "FAILED"
)

type MovieService struct {
	repo outbound.MovieRepoPort
	sub  outbound.SubscriptionPort
}

func NewMovieService(repo outbound.MovieRepoPort, sub outbound.SubscriptionPort) *MovieService {
	return &MovieService{
		repo: repo,
		sub:  sub,
	}
}

func (s *MovieService) CreateMovie(movie *domain.Movie) error {
	return s.repo.CreateMovie(movie)
}

func (s *MovieService) GetMovies() ([]*domain.Movie, error) {
	return s.repo.GetAllMovies()
}

func (s *MovieService) GetMovieByID(movieID uuid.UUID) (*domain.Movie, error) {
	return s.repo.GetMovieByID(movieID)
}

func (s *MovieService) GetMovieStreamURL(viewerID, movieID uuid.UUID) (string, error) {

	movie, err := s.repo.GetMovieByID(movieID)
	if err != nil {
		return "", err
	}

	// TODO: check payment / subscription
	status, err := s.sub.GetSubscriptionStatus(viewerID)

	if err != nil {
		return "", err
	}

	if status != StatusActive {
		return "", fmt.Errorf("Status is not active")
	}

	token, _ := utils.GenerateVideoToken(viewerID, movieID)

	url := fmt.Sprintf(
		"http://localhost:3000%s?token=%s",
		movie.HLSPath,
		token,
	)

	return url, nil
}
