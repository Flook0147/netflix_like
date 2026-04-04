package outbound

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/Flook0147/netflix_like/internal/movie/utils"
	"github.com/google/uuid"
)

type VideoProcessor struct {
	bucketName string
	client     *storage.Client
}

// 🔥 init client ครั้งเดียว (สำคัญมาก)
func NewVideoProcessor(bucket string) *VideoProcessor {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}

	return &VideoProcessor{
		bucketName: bucket,
		client:     client,
	}
}

// 🔥 main flow
func (v *VideoProcessor) ProcessVideo(movieID uuid.UUID, inputPath string) error {
	outputDir := fmt.Sprintf("./tmp/hls/%s", movieID.String())

	// create output folder
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	// convert mp4 → HLS
	if err := utils.ConvertToHLS(inputPath, outputDir); err != nil {
		return err
	}

	// upload all files
	if err := v.uploadDirectory(outputDir, movieID.String()); err != nil {
		return err
	}

	// cleanup
	if err := os.RemoveAll(outputDir); err != nil {
		fmt.Println("failed to clean output:", err)
	}

	return nil
}

// 🔥 generate signed url
func (v *VideoProcessor) GenerateSignedURL(objectPath string) (string, error) {
	serviceAccountEmail := os.Getenv("GCS_SERVICE_ACCOUNT_EMAIL")
	privateKey := []byte(strings.ReplaceAll(os.Getenv("GCS_PRIVATE_KEY"), `\n`, "\n"))

	url, err := storage.SignedURL(v.bucketName, objectPath, &storage.SignedURLOptions{
		Method:         "GET",
		Expires:        time.Now().Add(5 * time.Minute),
		Scheme:         storage.SigningSchemeV4,
		GoogleAccessID: serviceAccountEmail,
		PrivateKey:     privateKey,
	})
	if err != nil {
		return "", err
	}

	return url, nil
}

// 🔥 upload file (มี timeout + content type)
func (v *VideoProcessor) uploadFile(objectName, filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := v.client.Bucket(v.bucketName).Object(objectName).NewWriter(ctx)

	// set content type
	if strings.HasSuffix(objectName, ".m3u8") {
		wc.ContentType = "application/vnd.apple.mpegurl"
	} else if strings.HasSuffix(objectName, ".ts") {
		wc.ContentType = "video/MP2T"
	}

	if _, err := io.Copy(wc, f); err != nil {
		return err
	}

	return wc.Close()
}

// 🔥 retry upload
func (v *VideoProcessor) uploadWithRetry(objectName, filePath string) error {
	var err error

	for i := 0; i < 3; i++ {
		err = v.uploadFile(objectName, filePath)
		if err == nil {
			return nil
		}

		fmt.Println("retry", i+1, "error:", err)
		time.Sleep(2 * time.Second)
	}

	return err
}

// 🔥 upload ทั้ง folder
func (v *VideoProcessor) uploadDirectory(localDir, movieID string) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// get relative path
		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		// fix Windows path
		relPath = filepath.ToSlash(relPath)

		objectName := fmt.Sprintf("videos/%s/%s", movieID, relPath)

		fmt.Println("Uploading:", objectName)

		// upload with retry
		if err := v.uploadWithRetry(objectName, path); err != nil {
			return err
		}

		return nil
	})
}
