package handler

import (
	"context"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// DefaultHandler processes messages that don't match specific handlers
func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	if ShareEmailHandler(ctx, b, update) {
		return
	}

	ErrorHandler(ctx, b, update)
}

// RegisterHandlers registers all bot handlers
func RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "share", bot.MatchTypeCommand, ShareCommandHandler)
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "share_", bot.MatchTypePrefix, ShareCallbackHandler)
}
