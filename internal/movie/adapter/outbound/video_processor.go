package outbound

import (
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Flook0147/netflix_like/internal/movie/utils"
	"github.com/google/uuid"
)

type VideoProcessor struct {
	bucketName string
}

func NewVideoProcessor(bucket string) *VideoProcessor {
	return &VideoProcessor{
		bucketName: bucket,
	}
}

func (v *VideoProcessor) ProcessVideo(movieID uuid.UUID, inputPath string) error {
	outputDir := fmt.Sprintf("/hls/%s", movieID)

	err := utils.ConvertToHLS(inputPath, outputDir)
	if err != nil {
		return err
	}

	// update DB (you may inject repo here)
	return nil
}

func (v *VideoProcessor) GenerateSignedURL(objectPath string) (string, error) {
	serviceAccountEmail := os.Getenv("GCS_SERVICE_ACCOUNT_EMAIL")
	privateKey := []byte(os.Getenv("GCS_PRIVATE_KEY"))

	url, err := storage.SignedURL(v.bucketName, objectPath, &storage.SignedURLOptions{
		Method:         "GET",
		Expires:        time.Now().Add(5 * time.Minute), // 🔥 อายุ 5 นาที
		Scheme:         storage.SigningSchemeV4,
		GoogleAccessID: serviceAccountEmail,
		PrivateKey:     privateKey,
	})
	if err != nil {
		return "", err
	}

	return url, nil
}
