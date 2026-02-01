package handler

import (
	"context"

	"go-tele-bot/internal/app/middleware"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// DefaultHandler processes messages that don't match specific handlers
func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	middleware.WithAuth(func(ctx context.Context, b *bot.Bot, update *models.Update) {
		// Handle LC mode messages first
		if LcMessageHandler(ctx, b, update) {
			return
		}

		// Handle stateful flows (e.g., share email input)
		if ShareEmailInputHandler(ctx, b, update) {
			return
		}

		// No match - show error
		ErrorHandler(ctx, b, update)
	})(ctx, b, update)
}

// RegisterHandlers registers all bot handlers
func RegisterHandlers(b *bot.Bot) {
	b.RegisterHandler(bot.HandlerTypeMessageText, "share", bot.MatchTypeCommand, middleware.WithAuth(ShareCommandHandler))
	b.RegisterHandler(bot.HandlerTypeCallbackQueryData, "share_", bot.MatchTypePrefix, ShareCallbackHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "start", bot.MatchTypeCommand, middleware.WithAuth(StartCommandHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "lc", bot.MatchTypeCommand, middleware.WithAuth(LcHandler))
	b.RegisterHandler(bot.HandlerTypeMessageText, "exit", bot.MatchTypeCommand, middleware.WithAuth(ExitLcHandler))
}
