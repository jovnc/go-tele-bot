package app

import (
	"context"
	"fmt"
	"log"
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

/**
 * Star bot using long polling
 * Useful for development and testing purposes
 */
func (b *TelegramBot) StartPolling(ctx context.Context) error {
	// Delete webhook before starting polling
	_, err := b.bot.DeleteWebhook(ctx, &bot.DeleteWebhookParams{
		DropPendingUpdates: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete webhook: %w", err)
	}
	log.Println("Webhook deleted, starting polling mode...")
	
	b.bot.Start(ctx)
	return nil
}

/**
 * Start bot using webhook (for production)
 * Useful for serverless environments to reduce resource usage
 */
func (b *TelegramBot) StartWebhook(ctx context.Context) error {
	_, err := b.bot.SetWebhook(ctx, &bot.SetWebhookParams{
		URL: config.GlobalConfig.WebhookURL,
	})
	if err != nil {
		return fmt.Errorf("failed to set webhook: %w", err)
	}

	go func() {
		addr := ":" + config.GlobalConfig.Port
		log.Printf("starting webhook server on %s", addr)
		if err := http.ListenAndServe(addr, b.bot.WebhookHandler()); err != nil {
			log.Fatalf("failed to start webhook server: %v", err)
		}
	}()

	b.bot.StartWebhook(ctx)
	return nil
}
