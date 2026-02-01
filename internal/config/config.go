package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BotToken       string
	WebhookURL     string
	Port           string
	ValidUsernames []string
	OpenAIAPIKey   string
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
	validUsernames := strings.Split(getEnv("VALID_USERNAMES", ""), ";")
	
	openAIAPIKey := getEnv("OPEN_AI_API_KEY", "")
	if openAIAPIKey == "" {
		return fmt.Errorf("OPEN_AI_API_KEY is required")
	}

	GlobalConfig = &Config{
		BotToken:       botToken,
		WebhookURL:     webhookURL,
		Port:           port,
		ValidUsernames: validUsernames,
		OpenAIAPIKey:   openAIAPIKey,
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
