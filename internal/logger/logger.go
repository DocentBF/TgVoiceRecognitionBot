package logger

import (
	"log"
	"os"
)

type Logger struct {
	internalLogger *log.Logger
}

var instance *Logger

func init() {
	instance = &Logger{
		internalLogger: log.New(os.Stdout, "[TgVoiceRecognitionBot] ", log.LstdFlags|log.Lshortfile),
	}
}

func Info(msg string) {
	instance.internalLogger.Println("[INFO] " + msg)
}

func Error(msg string) {
	instance.internalLogger.Println("[ERROR] " + msg)
}

func Debug(msg string) {
	instance.internalLogger.Println("[DEBUG] " + msg)
}

func Printf(format string, v ...interface{}) {
	instance.internalLogger.Printf(format, v...)
}
