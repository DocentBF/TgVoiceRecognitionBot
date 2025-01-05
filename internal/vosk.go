package main

import (
	"TgVoiceRecognitionBot/internal/logger"
	"fmt"
	vosk "github.com/alphacep/vosk-api/go"
	"os"
)

type recognizerWrapper struct {
	model      *vosk.VoskModel
	recognizer *vosk.VoskRecognizer
}

func initVoskRecognizer(modelPath string) (*recognizerWrapper, func()) {
	model, err := vosk.NewModel(modelPath)
	if err != nil {
		logger.Error(fmt.Sprintf("Не удалось загрузить модель: %v", err))
		os.Exit(1)
	}

	recognizer, err := vosk.NewRecognizer(model, 16000)
	if err != nil {
		model.Free()
		logger.Error(fmt.Sprintf("Не удалось создать распознаватель: %v", err))
		os.Exit(1)
	}

	recognizer.SetWords(1)
	recognizer.SetPartialWords(1)

	cleanup := func() {
		recognizer.Free()
		model.Free()
	}

	return &recognizerWrapper{model: model, recognizer: recognizer}, cleanup
}
