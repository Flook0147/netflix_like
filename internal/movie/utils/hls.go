package utils

import (
	"fmt"
	"os/exec"
)

func ConvertToHLS(inputPath, outputDir string) error {
	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-codec:", "copy",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		fmt.Sprintf("%s/index.m3u8", outputDir),
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg error: %s", string(output))
	}

	return nil
}
