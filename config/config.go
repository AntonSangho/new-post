package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	TelegramChatID   int64
	TargetURL        string
	CheckInterval    string
	DataFile         string
}

func Load() (*Config, error) {
	godotenv.Load()

	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		chatID = 0
	}

	config := &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   chatID,
		TargetURL:        os.Getenv("TARGET_URL"),
		CheckInterval:    getEnvWithDefault("CHECK_INTERVAL", "*/5 * * * *"),
		DataFile:         getEnvWithDefault("DATA_FILE", "data.json"),
	}

	return config, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}