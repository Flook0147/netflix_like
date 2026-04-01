package outbound

import "github.com/google/uuid"

type VideoProcessorPort interface {
	ProcessVideo(movieID uuid.UUID, inputPath string) error
	GenerateSignedURL(objectPath string) (string, error)
}
