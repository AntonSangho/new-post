package main

import (
	"log"
	"sort"

	"go-newpost/config"
	"go-newpost/scraper"
	"go-newpost/telegram"
)

type Monitor struct {
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

	if len(cfg.TargetURLs) == 0 {
		log.Fatal("TARGET_URL is required")
	}

	monitor, err := NewMonitor(cfg)
	if err != nil {
		log.Fatalf("Failed to create monitor: %v", err)
	}

	log.Printf("Monitoring %d URLs: %v", len(cfg.TargetURLs), cfg.TargetURLs)
	
	monitor.CheckForNewPosts()
	
	log.Println("Monitor completed successfully")
}

func NewMonitor(cfg *config.Config) (*Monitor, error) {
	telegramBot, err := telegram.New(cfg.TelegramBotToken, cfg.TelegramChatID)
	if err != nil {
		return nil, err
	}

	return &Monitor{
		telegram: telegramBot,
		config:   cfg,
	}, nil
}

func (m *Monitor) CheckForNewPosts() {
	log.Println("Checking for recent posts...")

	var allPosts []scraper.Post

	// 각 URL에 대해 스크래핑 실행
	for i, url := range m.config.TargetURLs {
		log.Printf("Scraping URL %d/%d: %s", i+1, len(m.config.TargetURLs), url)
		
		s := scraper.New(url)
		posts, err := s.GetRecentPosts(7)
		if err != nil {
			log.Printf("Failed to scrape URL %s: %v", url, err)
			continue // 하나 실패해도 다른 URL은 계속 처리
		}
		
		log.Printf("Found %d posts from %s", len(posts), url)
		allPosts = append(allPosts, posts...)
	}

	// 날짜순으로 정렬 (최신순)
	sort.Slice(allPosts, func(i, j int) bool {
		return allPosts[i].PublishDate > allPosts[j].PublishDate
	})

	// 중복 제거 (같은 제목+날짜)
	allPosts = removeDuplicates(allPosts)

	log.Printf("Found %d total recent posts (last 7 days)", len(allPosts))

	if err := m.telegram.SendPostSummary(allPosts); err != nil {
		log.Printf("Failed to send post summary: %v", err)
		m.telegram.SendErrorMessage("알림 전송 중 오류가 발생했습니다: " + err.Error())
		return
	}

	log.Println("Post summary sent successfully")
}

func removeDuplicates(posts []scraper.Post) []scraper.Post {
	seen := make(map[string]bool)
	var result []scraper.Post
	
	for _, post := range posts {
		// 제목+날짜로 중복 체크 (소스는 제외)
		key := post.Title + post.PublishDate
		if !seen[key] {
			seen[key] = true
			result = append(result, post)
		}
	}
	
	return result
}