package app

import (
	"context"
	"fmt"
	"net/http"

	"go-tele-bot/internal/app/handler"
	"go-tele-bot/internal/config"

	"github.com/go-telegram/bot"
)

type TelegramBot struct {
	bot *bot.Bot
}

func NewTelegramBot(botToken string) (*TelegramBot, error) {
	opts := []bot.Option{
		bot.WithDefaultHandler(handler.DefaultHandler),
	}

	b, err := bot.New(botToken, opts...)
	if err != nil {
		return nil, err
	}

	handler.RegisterHandlers(b)

	return &TelegramBot{
		bot: b,
	}, nil
}

// Start bot using long polling
func (b *TelegramBot) Start(ctx context.Context) {
	b.bot.Start(ctx)
}

// Start bot using webhook (preferred)
func (b *TelegramBot) StartWebhook(ctx context.Context) error {
	_, err := b.bot.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: config.GlobalConfig.WebhookURL,
	})
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	go func() {
		addr := ":" + config.GlobalConfig.Port
		fmt.Printf("Starting webhook server on %s\n", addr)
		if err := http.ListenAndServe(addr, b.bot.WebhookHandler()); err != nil {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	b.bot.StartWebhook(ctx)
	return nil
}
