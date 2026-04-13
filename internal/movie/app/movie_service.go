package app

import (
	"fmt"
	"os"

	"github.com/Flook0147/netflix_like/internal/movie/domain"
	"github.com/Flook0147/netflix_like/internal/movie/port/outbound"
	"github.com/google/uuid"
)

const (
	StatusPending  = "PENDING"
	StatusActive   = "ACTIVE"
	StatusInactive = "INACTIVE"
	StatusFailed   = "FAILED"
)

type MovieService struct {
	repo      outbound.MovieRepoPort
	sub       outbound.SubscriptionPort
	processor outbound.VideoProcessorPort
}

func NewMovieService(repo outbound.MovieRepoPort, sub outbound.SubscriptionPort, processor outbound.VideoProcessorPort) *MovieService {
	return &MovieService{
		repo:      repo,
		sub:       sub,
		processor: processor,
	}
}

func (s *MovieService) CreateMovie(movie *domain.Movie, inputPath string) error {
	err := s.repo.CreateMovie(movie)
	if err != nil {
		return err
	}

	// 🔥 async processing
	go func() {
		fmt.Println("Start processing:", movie.ID)

		// ✅ ลบไฟล์หลังใช้งานเสร็จ (อยู่ใน goroutine)
		defer func() {
			err := os.Remove(inputPath)
			if err != nil {
				fmt.Println("failed to remove temp file:", err)
			} else {
				fmt.Println("temp file removed:", inputPath)
			}
		}()

		err := s.processor.ProcessVideo(movie.ID, inputPath) // ✅ ใช้ mp4 path
		if err != nil {
			fmt.Println("process error:", err)
			return
		}

		fmt.Println("Done processing:", movie.ID)
	}()

	return nil
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

	// status, err := s.sub.GetSubscriptionStatus(viewerID)
	// if err != nil {
	// 	return "", err
	// }

	// if status != StatusActive {
	// 	return "", fmt.Errorf("subscription is not active")
	// }

	url, err := s.processor.GenerateSignedURL(movie.HLSPath)
	if err != nil {
		return "", err
	}

	return url, nil
}
