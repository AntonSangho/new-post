package main

import (
	"log"
	"time"

	"go-newpost/config"
	"go-newpost/scraper"
	"go-newpost/storage"
	"go-newpost/telegram"
)

type Monitor struct {
	scraper  *scraper.Scraper
	storage  *storage.Storage
	telegram *telegram.Bot
	config   *config.Config
}

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if cfg.TelegramBotToken == "" {
		log.Fatal("TELEGRAM_BOT_TOKEN is required")
	}

	if cfg.TelegramChatID == 0 {
		log.Fatal("TELEGRAM_CHAT_ID is required")
	}

	if cfg.TargetURL == "" {
		log.Fatal("TARGET_URL is required")
	}

	monitor, err := NewMonitor(cfg)
	if err != nil {
		log.Fatalf("Failed to create monitor: %v", err)
	}

	log.Printf("Monitoring URL: %s", cfg.TargetURL)
	
	monitor.CheckForNewPosts()
	
	log.Println("Monitor completed successfully")
}

func NewMonitor(cfg *config.Config) (*Monitor, error) {
	telegramBot, err := telegram.New(cfg.TelegramBotToken, cfg.TelegramChatID)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		scraper:  scraper.New(cfg.TargetURL),
		storage:  storage.New(cfg.DataFile),
		telegram: telegramBot,
		config:   cfg,
	}, nil
}

func (m *Monitor) CheckForNewPosts() {
	log.Println("Checking for new posts...")

	posts, err := m.scraper.GetLatestPosts()
	if err != nil {
		log.Printf("Failed to scrape posts: %v", err)
		m.telegram.SendErrorMessage("스크래핑 중 오류가 발생했습니다: " + err.Error())
		return
	}

	log.Printf("Found %d posts", len(posts))

	newPosts, err := m.storage.GetNewPosts(posts)
	if err != nil {
		log.Printf("Failed to get new posts: %v", err)
		m.telegram.SendErrorMessage("데이터 처리 중 오류가 발생했습니다: " + err.Error())
		return
	}

	log.Printf("Found %d new posts", len(newPosts))

	for _, post := range newPosts {
		log.Printf("New post: %s", post.Title)
		if err := m.telegram.SendNewPostNotification(post); err != nil {
			log.Printf("Failed to send notification for post %s: %v", post.Title, err)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}