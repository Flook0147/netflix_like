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

func (v *VideoProcessor) ProcessVideo(movieID uuid.UUID, inputPath string) error {
	outputDir := fmt.Sprintf("./tmp/hls/%s", movieID.String())

	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	// ✅ convert mp4 → HLS
	if err := utils.ConvertToHLS(inputPath, outputDir); err != nil {
		return err
	}

	// ✅ rewrite m3u8 → signed ts
	if err := v.rewriteM3U8(outputDir, movieID.String()); err != nil {
		return err
	}

	// ✅ upload
	if err := v.uploadDirectory(outputDir, movieID.String()); err != nil {
		return err
	}
	// cleanup
	if err := os.RemoveAll(outputDir); err != nil {
		fmt.Println("failed to clean output:", err)
	}

	return nil
}

func (v *VideoProcessor) GenerateSignedURL(objectPath string) (string, error) {
	return v.client.Bucket(v.bucketName).SignedURL(objectPath, &storage.SignedURLOptions{
		Method:  "GET",
		Expires: time.Now().Add(5 * time.Minute),
	})
}

func (v *VideoProcessor) rewriteM3U8(dir, movieID string) error {
	m3u8Path := filepath.Join(dir, "index.m3u8")

	data, err := os.ReadFile(m3u8Path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// 🎯 rewrite only ts lines
		if strings.HasSuffix(line, ".ts") {
			objectPath := fmt.Sprintf("videos/%s/%s", movieID, line)

			signedURL, err := v.GenerateSignedURL(objectPath)
			if err != nil {
				return err
			}

			fmt.Println("Rewrite:", line, "→", signedURL)

			lines[i] = signedURL
		}
	}

	newContent := strings.Join(lines, "\n")
	return os.WriteFile(m3u8Path, []byte(newContent), 0644)
}

func (v *VideoProcessor) uploadFile(objectName, filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	wc := v.client.Bucket(v.bucketName).Object(objectName).NewWriter(ctx)

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

func (v *VideoProcessor) uploadWithRetry(objectName, filePath string) error {
	var err error

	for i := 0; i < 3; i++ {
		err = v.uploadFile(objectName, filePath)
		if err == nil {
			return nil
		}

		fmt.Println("Retry", i+1, "error:", err)
		time.Sleep(2 * time.Second)
	}

	return err
}

func (v *VideoProcessor) uploadDirectory(localDir, movieID string) error {
	return filepath.Walk(localDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(localDir, path)
		if err != nil {
			return err
		}

		relPath = filepath.ToSlash(relPath)

		// 🔥 FIX path
		relPath = strings.TrimPrefix(relPath, "./")
		relPath = strings.TrimPrefix(relPath, "/")

		if strings.Contains(relPath, "..") {
			return fmt.Errorf("invalid path: %s", relPath)
		}

		objectName := fmt.Sprintf("videos/%s/%s", movieID, relPath)

		fmt.Println("Uploading:", objectName)

		return v.uploadWithRetry(objectName, path)
	})
}
