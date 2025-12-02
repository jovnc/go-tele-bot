package middleware

import (
	"context"
	"go-tele-bot/internal/config"
	"slices"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	messageAuthError = "🔒 <b>Access denied</b>\n\n" +
		"You are not authorized to use this bot."
)

// WithAuth wraps a handler with authorization check
func WithAuth(next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		username, chatID := update.Message.From.Username, update.Message.Chat.ID

		if !slices.Contains(config.GlobalConfig.ValidUsernames, username) {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: chatID,
				Text:   messageAuthError,
				ParseMode: models.ParseModeHTML,
			})
			return
		}

		next(ctx, b, update)
	}
}