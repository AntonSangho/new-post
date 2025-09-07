package telegram

import (
	"fmt"
	"log"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-newpost/scraper"
)

type Bot struct {
	api    *tgbotapi.BotAPI
	chatID int64
}

func New(token string, chatID int64) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &Bot{
		api:    bot,
		chatID: chatID,
	}, nil
}

func (b *Bot) SendNewPostNotification(post scraper.Post) error {
	message := fmt.Sprintf("🆕 새로운 주택공고가 등록되었습니다!\n\n"+
		"📝 공고명: %s\n"+
		"📅 공고게시일: %s\n"+
		"🔗 링크: %s\n"+
		"⏰ 발견 시간: %s",
		post.Title,
		post.PublishDate,
		post.URL,
		post.Date,
	)

	msg := tgbotapi.NewMessage(b.chatID, message)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.DisableWebPagePreview = false

	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func (b *Bot) SendStartupMessage() error {
	message := "🏠 서울시 주택공고 모니터링 봇이 시작되었습니다.\n새로운 공고가 발견되면 알림을 보내드립니다."
	
	msg := tgbotapi.NewMessage(b.chatID, message)
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send startup message: %w", err)
	}

	return nil
}

func (b *Bot) SendPostSummary(posts []scraper.Post) error {
	if len(posts) == 0 {
		return b.SendNoPostsAlert()
	}

	message := fmt.Sprintf("🏠 서울시 주택공고 현황 (최근 7일)\n\n📝 총 %d건의 공고가 있습니다\n\n", len(posts))

	for i, post := range posts {
		if i >= 10 { // 최대 10개만 표시
			message += fmt.Sprintf("... 외 %d건 더 있습니다", len(posts)-10)
			break
		}

		statusText := ""
		if post.Status != "" {
			statusText = fmt.Sprintf(" [%s]", post.Status)
		}

		message += fmt.Sprintf("%d️⃣ [%s] %s%s\n📅 %s\n🔗 %s\n\n",
			i+1,
			post.Source,
			post.Title,
			statusText,
			post.PublishDate,
			post.URL,
		)
	}

	msg := tgbotapi.NewMessage(b.chatID, message)
	msg.DisableWebPagePreview = true // 링크 미리보기 비활성화
	
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send post summary: %w", err)
	}

	return nil
}

func (b *Bot) SendNoPostsAlert() error {
	message := fmt.Sprintf("🏠 서울시 주택공고 확인 완료\n\n📋 최근 7일간 새로운 공고가 없습니다.\n⏰ 확인 시간: %s",
		time.Now().Format("2006-01-02 15:04"))
	
	msg := tgbotapi.NewMessage(b.chatID, message)
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send no posts alert: %w", err)
	}

	return nil
}

func (b *Bot) SendErrorMessage(errMsg string) error {
	message := fmt.Sprintf("❌ 오류가 발생했습니다:\n%s", errMsg)
	
	msg := tgbotapi.NewMessage(b.chatID, message)
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send error message: %w", err)
	}

	return nil
}