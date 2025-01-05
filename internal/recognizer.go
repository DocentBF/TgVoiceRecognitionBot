package main

import (
	"TgVoiceRecognitionBot/internal/logger"
	"encoding/json"
	"io"
	"os"
)

type RecognitionResult struct {
	Text string `json:"text"`
}

func transcribeAudio(filePath string, recognizer *recognizerWrapper) string {
	audioFile, err := os.Open(filePath)
	if err != nil {
		logger.Printf("Не удалось открыть аудиофайл: %v", err)
		return ""
	}
	defer audioFile.Close()

	buf := make([]byte, 4096)
	resultStr := ""
	recognizer.recognizer.Reset()
	for {
		n, err := audioFile.Read(buf)
		if err != nil && err != io.EOF {
			logger.Printf("Ошибка чтения файла: %v", err)
			return ""
		}
		if n == 0 {
			break
		}

		if recognizer.recognizer.AcceptWaveform(buf[:n]) != 0 {
			bufferResultStr := recognizer.recognizer.Result()
			var result RecognitionResult
			if err := json.Unmarshal([]byte(bufferResultStr), &result); err != nil {
				logger.Printf("Не удалось распознать результат: %v", err)
				continue
			}
			if result.Text != "" {
				resultStr += " " + result.Text
			}
		}
	}
	return resultStr
}
