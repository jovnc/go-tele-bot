package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"go-tele-bot/internal/app"
	"go-tele-bot/internal/app/handler"
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

	// Intialise telegram bot handlers
	if err := handler.InitialiseHandlers(); err != nil {
		log.Fatalf("Failed to initialize handlers: %v", err)
	}

	// Start the bot
	bot, err := app.NewTelegramBot(config.GlobalConfig.BotToken)
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	// Run using webhook in production / polling for local
	isLocal := os.Getenv("DEV") == "true"
	if !isLocal {
		log.Println("Starting bot in webhook mode...")
		err = bot.StartWebhook(ctx)
		if err != nil {
			log.Fatalf("Failed to start webhook: %v", err)
		}
	} else {
		log.Println("Starting bot in polling mode (local dev)...")
		err = bot.StartPolling(ctx)
		if err != nil {
			log.Fatalf("Failed to start polling: %v", err)
		}
	}
}
