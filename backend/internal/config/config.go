package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL     string
	TelegramToken   string
	TelegramWebApp  string
	APIPort         int
	BotWebhookURL   string
	UseWebhook      bool
}

func Load() *Config {
	port, _ := strconv.Atoi(getEnv("API_PORT", "8080"))
	
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://localhost:5432/football_bot"),
		TelegramToken:  getEnv("TELEGRAM_BOT_TOKEN", ""),
		TelegramWebApp: getEnv("TELEGRAM_WEBAPP_URL", ""),
		APIPort:        port,
		BotWebhookURL:  getEnv("BOT_WEBHOOK_URL", ""),
		UseWebhook:     getEnv("USE_WEBHOOK", "false") == "true",
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
