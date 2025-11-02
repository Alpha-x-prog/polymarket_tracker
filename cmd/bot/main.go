package main

import (
	"log"
	"os"
	"polymarket_tg_bot/internal/bot"

	"github.com/joho/godotenv"
)

func mustEnv(key string) string {

	_ = godotenv.Load(".env")
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("environment variable %s is not set", key)
	}
	return v
}

func main() {
	token := mustEnv("TG_BOT_TOKEN")
	b, err := bot.New(token)
	if err != nil {
		log.Fatalf("init bot: %v", err)
	}

	if err := b.Run(); err != nil {
		log.Fatalf("run bot: %v", err)
	}
}
