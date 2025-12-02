package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	messageError = "🤔 Unknown command"
)

// ErrorHandler sends a message to the user if the message is not understood
func ErrorHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      messageError,
		ParseMode: models.ParseModeHTML,
	})
}
