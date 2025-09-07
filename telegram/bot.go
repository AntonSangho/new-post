package telegram

import (
	"fmt"
	"log"

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
	message := fmt.Sprintf("🆕 새로운 게시글이 등록되었습니다!\n\n"+
		"📝 제목: %s\n"+
		"🔗 링크: %s\n"+
		"📅 발견 시간: %s",
		post.Title,
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
	message := "🤖 게시판 모니터링 봇이 시작되었습니다.\n새로운 글이 발견되면 알림을 보내드립니다."
	
	msg := tgbotapi.NewMessage(b.chatID, message)
	_, err := b.api.Send(msg)
	if err != nil {
		return fmt.Errorf("failed to send startup message: %w", err)
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