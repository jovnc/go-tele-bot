package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	messageStart = "👋 <b>Welcome!</b>\n\n" +
		"I help you share Google Drive folders with others.\n\n" +
		"<b>Available commands:</b>\n" +
		"📤 /share — Share folders with someone"
)

// StartCommandHandler handles the /start command
func StartCommandHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      messageStart,
		ParseMode: models.ParseModeHTML,
	})
}
