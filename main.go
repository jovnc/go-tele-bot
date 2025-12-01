package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"go-tele-bot/internal/app"
	"go-tele-bot/internal/config"
	"go-tele-bot/internal/data"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Initialize config
	err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load data
	if err := data.Load(); err != nil {
		log.Fatalf("Failed to load data: %v", err)
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
