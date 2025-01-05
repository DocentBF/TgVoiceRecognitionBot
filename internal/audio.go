package main

import (
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func downloadFile(url string) (string, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	tmpFile, err := os.CreateTemp("", "*.ogg")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, response.Body); err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}

func convertToWav(filePath string) (string, error) {
	wavFilePath := filePath[:len(filePath)-len(filepath.Ext(filePath))] + ".wav"
	cmd := exec.Command("ffmpeg", "-i", filePath, "-ar", "16000", "-ac", "1", wavFilePath)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return wavFilePath, nil
}
