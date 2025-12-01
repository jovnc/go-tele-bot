package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken   string
	WebhookURL string
	Port       string
}

var GlobalConfig *Config

func LoadConfig() error {
	// Load .env file if it exists
	_ = godotenv.Load()

	botToken := getEnv("BOT_TOKEN", "")
	if botToken == "" {
		return fmt.Errorf("BOT_TOKEN is required")
	}

	webhookURL := getEnv("WEBHOOK_URL", "")
	if webhookURL == "" {
		return fmt.Errorf("WEBHOOK_URL is required")
	}

	port := getEnv("PORT", "8080")

	GlobalConfig = &Config{
		BotToken:   botToken,
		WebhookURL: webhookURL,
		Port:       port,
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
