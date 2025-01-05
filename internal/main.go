package main

import (
	"TgVoiceRecognitionBot/internal/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	config := loadConfig()

	bot, err := initTelegramBot(config)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка инициализации бота: %v", err))
		os.Exit(1)
	}

	voskRecognizer, cleanup := initVoskRecognizer("models/vosk-model-small-ru-0.22")
	defer cleanup()

	logger.Info("Бот запущен")

	go handleShutdown()
	processUpdates(bot, voskRecognizer)
}

func handleShutdown() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-sigchan

	logger.Printf("Завершение работы")

	os.Exit(0)
}
