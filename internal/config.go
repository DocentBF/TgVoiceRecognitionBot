package main

import (
	"TgVoiceRecognitionBot/internal/logger"
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	TelegramToken string  `json:"telegram_token"`
	IsBotDebug    bool    `json:"is_bot_debug"`
	AllowedUsers  []int64 `json:"allowed_users"`
}

func loadConfig() Config {
	file, err := os.Open("config.json")
	if err != nil {
		logger.Error(fmt.Sprintf("Не удалось открыть файл конфигурации: %v", err))
		os.Exit(1)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		logger.Error(fmt.Sprintf("Не удалось декодировать конфигурацию: %v", err))
		os.Exit(1)
	}
	return config
}
