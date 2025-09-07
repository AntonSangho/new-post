package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramBotToken string
	TelegramChatID   int64
	TargetURLs       []string
}

func Load() (*Config, error) {
	godotenv.Load()

	chatID, err := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
	if err != nil {
		chatID = 0
	}

	// TARGET_URL을 콤마로 분리하여 배열로 변환
	urlsStr := os.Getenv("TARGET_URL")
	var targetURLs []string
	if urlsStr != "" {
		urls := strings.Split(urlsStr, ",")
		for _, url := range urls {
			url = strings.TrimSpace(url)
			if url != "" {
				targetURLs = append(targetURLs, url)
			}
		}
	}

	config := &Config{
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		TelegramChatID:   chatID,
		TargetURLs:       targetURLs,
	}

	return config, nil
}
