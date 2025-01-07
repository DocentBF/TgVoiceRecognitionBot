package main

import (
	"TgVoiceRecognitionBot/internal/logger"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
	"slices"
)

var botConfig Config

func initTelegramBot(config Config) (*tgbotapi.BotAPI, error) {
	botConfig = config
	bot, err := tgbotapi.NewBotAPI(botConfig.TelegramToken)
	if err != nil {
		return nil, err
	}
	bot.Debug = botConfig.IsBotDebug
	return bot, nil
}

func processUpdates(bot *tgbotapi.BotAPI, recognizer *recognizerWrapper) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка получения обновлений: %v", err))
		os.Exit(1)
	}

	for update := range updates {
		if len(botConfig.AllowedUsers) > 0 && update.Message != nil && !isAllowedUser(int64(update.Message.From.ID)) {
			logger.Printf("Сообщение от неразрешенного пользователя: %d", update.Message.From.ID)
			continue
		}
		if update.Message.Voice != nil || update.Message.Audio != nil {
			sentMessage := replyMessage(bot, update.Message.Chat.ID, "Распознаю сообщение...")
			text, err := handleVoiceMessage(bot, update.Message, recognizer)
			if err != nil {
				text = "Не удалось распознать сообщение"
				logger.Printf("Не удалось распознать сообщение: %v", err)
			}
			editSentMessage(bot, update.Message.Chat.ID, sentMessage, text)
		}
		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Нажмите кнопку 'Пинг', чтобы проверить работу бота.")
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("/ping"),
				),
			)
			bot.Send(msg)

		case "/ping":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Понг! Бот работает.")
			bot.Send(msg)
		}
	}
}

func handleVoiceMessage(bot *tgbotapi.BotAPI, message *tgbotapi.Message, recognizer *recognizerWrapper) (
	string,
	error,
) {
	var fileID string
	text := ""
	if message.Voice != nil {
		fileID = message.Voice.FileID
	} else {
		fileID = message.Audio.FileID
	}

	fileURL, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		return "", fmt.Errorf("не удалось получить URL файла: %v", err)
	}

	filePath, err := downloadFile(fileURL)
	if err != nil {
		return "", fmt.Errorf("не удалось скачать файл: %v", err)
	}
	defer os.Remove(filePath)

	wavFilePath, err := convertToWav(filePath)
	if err != nil {
		return "", fmt.Errorf("не удалось конвертировать файл в WAV: %v", err)
	}
	defer os.Remove(wavFilePath)

	text = transcribeAudio(wavFilePath, recognizer)
	if text == "" {
		text = "Не удалось распознать голосовое сообщение."
	}

	return text, nil
}

func isAllowedUser(userID int64) bool {
	return slices.Contains(botConfig.AllowedUsers, userID)
}

func replyMessage(bot *tgbotapi.BotAPI, chatID int64, text string) *tgbotapi.Message {
	msg := tgbotapi.NewMessage(chatID, text)
	sentMessage, err := bot.Send(msg)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка при отравке сообщения: %v", err))
	}

	return &sentMessage
}

func editSentMessage(bot *tgbotapi.BotAPI, chatID int64, sentMessage *tgbotapi.Message, newText string) {
	editMsg := tgbotapi.NewEditMessageText(chatID, sentMessage.MessageID, newText)
	_, err := bot.Send(editMsg)
	if err != nil {
		logger.Error(fmt.Sprintf("Ошибка при редактировании сообщения: %v", err))
	}
}
