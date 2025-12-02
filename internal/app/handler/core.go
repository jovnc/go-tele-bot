package handler

import (
	"context"
	"go-tele-bot/internal/config"
	"slices"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// DefaultHandler processes messages that don't match specific handlers
func DefaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	username := update.Message.From.Username
	chatID := update.Message.Chat.ID

	// Check if username is in valid usernames
	if !slices.Contains(config.GlobalConfig.ValidUsernames, username) {
		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: chatID,
			Text:   "You are not authorized to use this bot.",
		})
		return
	}

	// Handle share command
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
