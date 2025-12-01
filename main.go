package main

import (
	"context"
	"go-tele-bot/internal/app"
	"go-tele-bot/internal/config"
	"log"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize config
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Start the bot
	bot, err := app.NewTelegramBot(config.GlobalConfig.BotToken)
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	err = bot.StartWebhook(ctx)
	if err != nil {
		log.Fatalf("Failed to start webhook: %v", err)
	}
}
